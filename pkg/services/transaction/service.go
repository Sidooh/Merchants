package transaction

import (
	"cmp"
	"fmt"
	"merchants.sidooh/api/presenter"
	"merchants.sidooh/pkg"
	"merchants.sidooh/pkg/clients"
	"merchants.sidooh/pkg/entities"
	"merchants.sidooh/pkg/logger"
	"merchants.sidooh/pkg/services/earning"
	"merchants.sidooh/pkg/services/earning_account"
	"merchants.sidooh/pkg/services/merchant"
	"merchants.sidooh/pkg/services/mpesa_store"
	"merchants.sidooh/pkg/services/payment"
	"merchants.sidooh/utils"
	"slices"
	"strconv"
	"strings"
)

type Service interface {
	FetchTransactions(filters Filters) (*[]presenter.Transaction, error)
	GetTransaction(id uint) (*presenter.Transaction, error)
	GetTransactionsByMerchant(merchantId uint) (*[]presenter.Transaction, error)
	UpdateTransaction(transaction *entities.Transaction) (*entities.Transaction, error)

	PurchaseFloat(transaction *entities.Transaction, agent, store string) (*entities.Transaction, error)
	WithdrawEarnings(transaction *entities.Transaction, source, destination, account string) (*entities.Transaction, error)

	CompleteTransaction(payment *entities.Payment, ipn *utils.Payment) error
}

type service struct {
	repository Repository

	merchantRepository   merchant.Repository
	paymentRepository    payment.Repository
	earningAccRepository earning_account.Repository
	earningRepository    earning.Repository
	mpesaStoreRepository mpesa_store.Repository

	earningAccService earning_account.Service
	earningService    earning.Service

	accountsApi *clients.ApiClient
	paymentsApi *clients.ApiClient
	notifyApi   *clients.ApiClient
}

func (s *service) FetchTransactions(filters Filters) (*[]presenter.Transaction, error) {
	if len(filters.Accounts) > 0 {
		merchants, err := s.merchantRepository.ReadMerchants(merchant.Filters{
			Columns:  []string{"account_id", "id"},
			Accounts: filters.Accounts,
		})
		if err != nil || len(*merchants) == 0 {
			return &[]presenter.Transaction{}, nil
		}

		for _, m := range *merchants {
			filters.Merchants = append(filters.Merchants, strconv.Itoa(int(m.Id)))
		}
	}

	return s.repository.ReadTransactions(filters)
}

func (s *service) GetTransaction(id uint) (results *presenter.Transaction, err error) {
	tx, err := s.repository.ReadTransaction(id)
	if err != nil {
		return nil, err
	}

	utils.ConvertStruct(tx, &results)

	return
}

func (s *service) GetTransactionsByMerchant(merchantId uint) (*[]presenter.Transaction, error) {
	return s.repository.ReadTransactionsByMerchant(merchantId)
}

func (s *service) UpdateTransaction(transaction *entities.Transaction) (*entities.Transaction, error) {
	return s.repository.UpdateTransaction(transaction)
}

func (s *service) PurchaseFloat(data *entities.Transaction, agent, store string) (tx *entities.Transaction, err error) {
	merchant, err := s.merchantRepository.ReadMerchant(data.MerchantId)
	if err != nil {
		return nil, err
	}

	tx, err = s.repository.CreateTransaction(data)
	if err != nil {
		return nil, err
	}

	payment, err := s.paymentsApi.BuyMpesaFloat(merchant.AccountId, merchant.FloatAccountId, int(tx.Amount), agent, store)
	fmt.Println("Payment", payment, err)
	if err != nil {
		tx.Status = "FAILED"
		tx, err := s.repository.UpdateTransaction(tx)
		if err != nil {
			return nil, err
		}

		logger.ClientLog.Error("Error buying float", "tx", tx, "error", err)

		go func() {
			account, _ := s.accountsApi.GetAccountById(strconv.Itoa(int(merchant.AccountId)))

			message := fmt.Sprintf("Sorry, KES%v Float could not be purchased", tx.Amount)

			s.notifyApi.SendSMS("DEFAULT", account.Phone, message)
		}()

		return nil, err
	}

	s.paymentRepository.CreatePayment(&entities.Payment{
		Amount: payment.Amount,
		Status: payment.Status,
		//Description:     payment.,
		Destination:   payment.Destination,
		TransactionId: tx.Id,
		PaymentId:     payment.Id,
	})

	return
}

func (s *service) WithdrawEarnings(data *entities.Transaction, source, destination, account string) (tx *entities.Transaction, err error) {
	merchant, err := s.merchantRepository.ReadMerchant(data.MerchantId)
	if err != nil {
		return nil, err
	}

	if destination == "FLOAT" && strconv.Itoa(int(merchant.FloatAccountId)) != account {
		return nil, pkg.ErrUnauthorized
	}

	charge := float32(s.getWithdrawalCharge(int(data.Amount)))
	if destination == "FLOAT" {
		charge = 0.0
	}

	var earningTXs []entities.EarningAccountTransaction

	if source != "" {
		earningAccount, err := s.earningAccRepository.ReadAccountByAccountIdAndType(merchant.AccountId, source)
		if err != nil {
			return nil, err
		}

		if earningAccount.Amount < data.Amount+charge {
			return nil, pkg.ErrInsufficientBalance
		}

		_, tx, err := s.earningAccService.DebitAccount(earningAccount.Id, data.Amount+charge)
		if err != nil {
			return nil, err
		}

		earningTXs = append(earningTXs, *tx)

	} else {
		earningAccounts, err := s.earningAccRepository.ReadAccountsByMerchant(data.MerchantId)
		if err != nil {
			return nil, err
		}

		//sort with highest balance first
		slices.SortFunc(earningAccounts, func(a, b entities.EarningAccount) int {
			return 0 - cmp.Compare(a.Amount, b.Amount) // reversed
		})

		var totalBalance float32

		for _, earningAccount := range earningAccounts {
			totalBalance += earningAccount.Amount
		}

		if totalBalance < data.Amount+charge {
			return nil, pkg.ErrInsufficientBalance
		}

		totalWithdrawal := data.Amount + charge

		for _, earningAccount := range earningAccounts {
			toDebit := totalWithdrawal

			if earningAccount.Amount > totalWithdrawal {
				totalWithdrawal -= totalWithdrawal
			} else {
				totalWithdrawal -= earningAccount.Amount
				toDebit = earningAccount.Amount
			}

			_, tx, err := s.earningAccService.DebitAccount(earningAccount.Id, toDebit)
			if err != nil {
				return nil, err
			}
			earningTXs = append(earningTXs, *tx)

			if totalWithdrawal == 0 {
				break
			}
		}
	}

	tx, err = s.repository.CreateTransaction(data)
	if err != nil {
		return nil, err
	}

	paymentData, err := s.paymentsApi.Withdraw(merchant.AccountId, 1, int(tx.Amount), destination, account)
	if err != nil {

		// TODO: reverse earningTXs
		for _, earningTx := range earningTXs {
			_, err := s.earningAccService.CreditAccount(earningTx.EarningAccountId, earningTx.Amount)
			if err != nil {
				return nil, err
			}
		}

		tx.Status = "FAILED"
		_, err := s.repository.UpdateTransaction(tx)

		return nil, err
	}

	payment, err := s.paymentRepository.CreatePayment(&entities.Payment{
		Amount: paymentData.Amount,
		Status: "PENDING",
		//Description:     payment.,
		Destination:   paymentData.Destination,
		TransactionId: tx.Id,
		PaymentId:     paymentData.Id,
	})
	if err != nil {
		return nil, err
	}

	if paymentData.Status == "COMPLETED" {
		err = s.CompleteTransaction(payment, paymentData)

		return nil, err
	}

	return
}

func (s *service) CompleteTransaction(payment *entities.Payment, ipn *utils.Payment) error {
	payment.Status = ipn.Status
	updatedPayment, err := s.paymentRepository.UpdatePayment(payment)
	if err != nil {
		return err
	}

	// TODO: convert this to return entity which means time conversion below can be removed
	transaction, err := s.repository.ReadTransaction(updatedPayment.TransactionId)
	mt, err := s.merchantRepository.ReadMerchant(transaction.MerchantId)

	switch transaction.Product {
	case "FLOAT":
		if ipn.Status == "FAILED" {

			date := transaction.CreatedAt.Format("02/01/2006, 3:04 PM")

			float, _ := s.paymentsApi.FetchFloatAccount(strconv.Itoa(int(mt.FloatAccountId)))

			message := fmt.Sprintf("Hi, we have added KES%v to your voucher account "+
				"because we could not complete your"+
				" KES%v float purchase for %s on %s. New voucher balance is KES%v.",
				transaction.Amount, transaction.Amount, *transaction.Destination, date, float.Balance)

			account, err := s.accountsApi.GetAccountById(strconv.Itoa(int(mt.AccountId)))
			if err != nil {
				return err
			}

			s.notifyApi.SendSMS("ERROR", account.Phone, message)

			transaction, err = s.UpdateTransaction(&entities.Transaction{
				ModelID: entities.ModelID{Id: transaction.Id},
				Status:  payment.Status,
			})

			return nil
		}

		err := s.computeCashback(mt, transaction, payment, ipn)
		if err != nil {
			return err
		}

	case "WITHDRAWAL":
		account, err := s.accountsApi.GetAccountById(strconv.Itoa(int(mt.AccountId)))
		if err != nil {
			return err
		}

		go func() {
			earningType := strings.Split(transaction.Description, " - ")[1]
			earningAcc, _ := s.earningAccRepository.ReadAccountByAccountIdAndType(mt.AccountId, earningType)
			destination := *transaction.Destination
			if strings.Split(*transaction.Destination, "-")[0] == "FLOAT" {
				destination = "VOUCHER"
			}
			date := transaction.CreatedAt.Format("02/01/2006, 3:04 PM")
			message := fmt.Sprintf("Withdrawal of KES%v from %s to %s on %s was successful. Cost KES%v. New %s Balance is KES%v",
				transaction.Amount, earningAcc.Type, destination, date, ipn.Charge, earningAcc.Type, earningAcc.Amount)
			if payment.Status == "FAILED" {
				message = fmt.Sprintf("Sorry, KES%v Withdrawal to %s could not be processed", transaction.Amount, destination)
			}
			s.notifyApi.SendSMS("DEFAULT", account.Phone, message)
		}()
	}

	// Update TX; after cashback?
	transaction, err = s.UpdateTransaction(&entities.Transaction{
		ModelID: entities.ModelID{Id: transaction.Id},
		Status:  payment.Status,
	})

	return err
}

func (s *service) computeCashback(mt *presenter.Merchant, tx *entities.Transaction, payment *entities.Payment, data *utils.Payment) error {
	// Compute cashback and commissions
	// Compute cashback
	cashback := float32(data.Charge) * .2

	s.earningRepository.CreateEarning(&entities.Earning{
		Amount:        cashback,
		Type:          "SELF",
		TransactionId: tx.Id,
		AccountId:     mt.AccountId,
	})

	earningAcc, err := s.earningAccRepository.ReadAccountByAccountIdAndType(mt.AccountId, "CASHBACK")
	if err != nil {
		earningAcc, err = s.earningAccRepository.CreateAccount(&entities.EarningAccount{
			Type:      "CASHBACK",
			AccountId: mt.AccountId,
		})
		if err != nil {
			return err
		}
	}
	s.earningAccService.CreditAccount(earningAcc.Id, cashback)
	//s.earningAccountService.DebitAccount(earningAcc.AccountId, cashback) // Debit acc for savings

	// Compute commissions
	commission := float32(data.Charge) * .1

	inviters, err := s.accountsApi.GetInviters(strconv.Itoa(int(mt.AccountId)))
	if err != nil {
		return err
	}

	if len(inviters) > 1 {
		for _, inviter := range inviters[1:] {
			s.earningRepository.CreateEarning(&entities.Earning{
				Amount:        commission,
				Type:          "INVITE",
				TransactionId: tx.Id,
				AccountId:     uint(inviter.Id),
			})

			earningAcc, err := s.earningAccRepository.ReadAccountByAccountIdAndType(uint(inviter.Id), "COMMISSION")
			if err != nil {
				earningAcc, err = s.earningAccRepository.CreateAccount(&entities.EarningAccount{
					Type:      "COMMISSION",
					AccountId: uint(inviter.Id),
				})
				if err != nil {
					return err
				}
			}
			s.earningAccService.CreditAccount(earningAcc.Id, commission)

		}
	}

	account, err := s.accountsApi.GetAccountById(strconv.Itoa(int(mt.AccountId)))
	if err != nil {
		return err
	}

	go func() {
		date := tx.CreatedAt.Format("02/01/2006, 3:04 PM")

		message := fmt.Sprintf("You have purchased KES%v float for %s on %s using Voucher. Cost KES%v. "+
			"You have received KES%v cashback.",
			payment.Amount,
			strings.Join(strings.Split(data.Store, " ")[0:4], " "), date, data.Charge, cashback)
		//message := fmt.Sprintf("KES%v Float for %s purchased successfully", payment.Amount, strings.Join(strings.Split(data.Store, " ")[0:4], " "))
		if payment.Status != "COMPLETED" {
			message = fmt.Sprintf("Sorry, KES%v Float for %s could not be purchased", payment.Amount, tx.Destination)
		}
		s.notifyApi.SendSMS("DEFAULT", account.Phone, message)
	}()

	_, _ = s.mpesaStoreRepository.CreateStore(&entities.MpesaAgentStoreAccount{
		Agent:      strings.Split(*tx.Destination, "-")[0],
		Store:      strings.Split(*tx.Destination, "-")[1],
		Name:       strings.Join(strings.Split(data.Store, " ")[0:4], " "),
		MerchantId: mt.Id,
	})

	// TODO: add go func with code to debit savings and send to save platform
	go s.earningService.SaveEarnings()

	return err
}

func (s *service) getWithdrawalCharge(amount int) int {
	charges, err := s.paymentsApi.GetWithdrawalCharges()
	if err != nil {
		return 0
	}

	for _, charge := range charges {
		if charge.Min <= amount && amount <= charge.Max {
			return charge.Charge
		}
	}

	return 0
}

func NewService(r Repository, merchantRepo merchant.Repository, paymentRepo payment.Repository, earningAccRepo earning_account.Repository, earningRepo earning.Repository, mpesaStoreRepo mpesa_store.Repository, earningAccSrv earning_account.Service, earningSrv earning.Service) Service {
	return &service{
		repository: r,

		merchantRepository:   merchantRepo,
		paymentRepository:    paymentRepo,
		earningAccRepository: earningAccRepo,
		earningRepository:    earningRepo,
		mpesaStoreRepository: mpesaStoreRepo,

		earningAccService: earningAccSrv,
		earningService:    earningSrv,

		accountsApi: clients.GetAccountClient(),
		paymentsApi: clients.GetPaymentClient(),
		notifyApi:   clients.GetNotifyClient(),
	}
}
