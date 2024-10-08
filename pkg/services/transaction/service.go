package transaction

import (
	"cmp"
	"context"
	"errors"
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
	"merchants.sidooh/pkg/services/savings"
	"merchants.sidooh/utils"
	"merchants.sidooh/utils/consts"
	"slices"
	"strconv"
	"strings"
)

type Service interface {
	FetchTransactions(filters Filters) (*[]presenter.Transaction, error)
	GetTransaction(id uint) (*presenter.Transaction, error)
	GetTransactionsByMerchant(merchantId uint) (*[]presenter.Transaction, error)
	UpdateTransaction(transaction *entities.Transaction) (*entities.Transaction, error)

	PurchaseMpesaFloat(transaction *entities.Transaction, agent, store, source, sourceAccount string) (*entities.Transaction, error)
	MpesaWithdrawal(transaction *entities.Transaction) (*entities.Transaction, error)
	FloatPurchase(transaction *entities.Transaction) (*entities.Transaction, error)
	FloatTransfer(transaction *entities.Transaction) (*entities.Transaction, error)
	FloatWithdraw(transaction *entities.Transaction, destination, account string) (*entities.Transaction, error)
	WithdrawEarnings(transaction *entities.Transaction, source, destination, account string) (*entities.Transaction, error)

	WithdrawSavings(transaction *entities.Transaction, source, destination, account string) (*entities.Transaction, error)

	CompleteTransaction(payment *entities.Payment, ipn *utils.Payment) error
}

type service struct {
	repository Repository

	merchantRepository   merchant.Repository
	paymentRepository    payment.Repository
	savingsRepository    savings.Repository
	earningAccRepository earning_account.Repository
	earningRepository    earning.Repository
	mpesaStoreRepository mpesa_store.Repository

	earningAccService earning_account.Service
	earningService    earning.Service

	accountsApi *clients.ApiClient
	paymentsApi *clients.ApiClient
	savingsApi  *clients.ApiClient
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

	results = &presenter.Transaction{
		Id:          tx.Id,
		Description: tx.Description,
		Destination: tx.Destination,
		Status:      tx.Status,
		Amount:      tx.Amount,
		MerchantId:  tx.MerchantId,
		Product:     tx.Product,
		CreatedAt:   tx.CreatedAt,
		UpdatedAt:   tx.UpdatedAt,
	}

	if tx.Payment != nil {
		results.Payment = &presenter.Payment{
			Id:            tx.Payment.Id,
			Description:   tx.Payment.Description,
			Destination:   tx.Payment.Destination,
			Amount:        tx.Payment.Amount,
			Charge:        tx.Payment.Charge,
			TransactionId: tx.Payment.TransactionId,
			Status:        tx.Payment.Status,
		}
	}

	return
}

func (s *service) GetTransactionsByMerchant(merchantId uint) (*[]presenter.Transaction, error) {
	return s.repository.ReadTransactionsByMerchant(merchantId)
}

func (s *service) UpdateTransaction(transaction *entities.Transaction) (*entities.Transaction, error) {
	return s.repository.UpdateTransaction(transaction)
}

func (s *service) PurchaseMpesaFloat(data *entities.Transaction, agent, store, source, sourceAccount string) (tx *entities.Transaction, err error) {
	merchant, err := s.merchantRepository.ReadMerchant(data.MerchantId)
	if err != nil {
		return nil, err
	}

	tx, err = s.repository.CreateTransaction(data)
	if err != nil {
		return nil, err
	}

	if source == "" || sourceAccount == "" {
		source = "FLOAT"
		sourceAccount = strconv.Itoa(int(merchant.FloatAccountId))
	}
	payment, err := s.paymentsApi.BuyMpesaFloat(merchant.AccountId, int(tx.Amount), agent, store, source, sourceAccount)
	if err != nil {
		logger.ClientLog.Error("Error buying mpesa float", "tx", tx, "error", err)

		if errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
		// TODO: check on connect timeout and exclude from failed tx...
		tx.Status = "FAILED"
		tx, err := s.repository.UpdateTransaction(tx)
		if err != nil {
			return nil, err
		}

		go func() {
			message := fmt.Sprintf("Sorry, KES%v Float could not be purchased", tx.Amount)

			s.notifyApi.SendSMS("DEFAULT", merchant.Phone, message)
		}()

		return nil, err
	}

	s.paymentRepository.CreatePayment(&entities.Payment{
		Amount: payment.Amount,
		Charge: float32(payment.Charge),
		Status: payment.Status,
		//Description:     payment.,
		Destination:   payment.Destination,
		TransactionId: tx.Id,
		PaymentId:     payment.Id,
	})

	return
}

func (s *service) MpesaWithdrawal(data *entities.Transaction) (tx *entities.Transaction, err error) {
	merchant, err := s.merchantRepository.ReadMerchant(data.MerchantId)
	if err != nil {
		return nil, err
	}

	tx, err = s.repository.CreateTransaction(data)
	if err != nil {
		return nil, err
	}

	payment, err := s.paymentsApi.MpesaWithdraw(merchant.AccountId, merchant.FloatAccountId, int(tx.Amount), *tx.Destination)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
		// TODO: check on connect timeout and exclude from failed tx...
		tx.Status = "FAILED"
		tx, err := s.repository.UpdateTransaction(tx)
		if err != nil {
			return nil, err
		}

		logger.ClientLog.Error("Error withdrawing mpesa", "tx", tx, "error", err)

		go func() {
			message := fmt.Sprintf("Sorry, KES%v Withdrawal could not be processed", tx.Amount)

			s.notifyApi.SendSMS("DEFAULT", merchant.Phone, message)
		}()

		return nil, err
	}

	s.paymentRepository.CreatePayment(&entities.Payment{
		Amount: payment.Amount,
		Charge: float32(payment.Charge),
		Status: payment.Status,
		//Description:     payment.,
		Destination:   payment.Destination,
		TransactionId: tx.Id,
		PaymentId:     payment.Id,
	})

	return
}

func (s *service) FloatPurchase(data *entities.Transaction) (tx *entities.Transaction, err error) {
	merchant, err := s.merchantRepository.ReadMerchant(data.MerchantId)
	if err != nil {
		return nil, err
	}

	tx, err = s.repository.CreateTransaction(data)
	if err != nil {
		return nil, err
	}

	payment, err := s.paymentsApi.FloatPurchase(merchant.AccountId, merchant.FloatAccountId, int(tx.Amount), *tx.Destination)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
		// TODO: check on connect timeout and exclude from failed tx...
		tx.Status = "FAILED"
		tx, err := s.repository.UpdateTransaction(tx)
		if err != nil {
			return nil, err
		}

		logger.ClientLog.Error("Error purchasing float", "tx", tx, "error", err)

		go func() {
			message := fmt.Sprintf("Sorry, KES%v Voucher purchase could not be processed", tx.Amount)

			s.notifyApi.SendSMS("DEFAULT", merchant.Phone, message)
		}()

		return nil, err
	}

	s.paymentRepository.CreatePayment(&entities.Payment{
		Amount: payment.Amount,
		Charge: float32(payment.Charge),
		Status: payment.Status,
		//Description:     payment.,ba
		Destination:   payment.Destination,
		TransactionId: tx.Id,
		PaymentId:     payment.Id,
	})

	return
}

func (s *service) FloatTransfer(data *entities.Transaction) (transaction *entities.Transaction, err error) {
	merchant, err := s.merchantRepository.ReadMerchant(data.MerchantId)
	if err != nil {
		return nil, err
	}

	recipientAccount, _ := strconv.Atoi(*data.Destination)
	recipient, err := s.merchantRepository.ReadMerchant(uint(recipientAccount))
	if err != nil {
		return nil, err
	}

	if recipient.FloatAccountId == merchant.FloatAccountId {
		return nil, pkg.ErrInvalidMerchant
	}

	transaction, err = s.repository.CreateTransaction(data)
	if err != nil {
		return nil, err
	}

	paymentData, err := s.paymentsApi.FloatTransfer(merchant.AccountId, merchant.FloatAccountId, int(transaction.Amount), strconv.Itoa(int(recipient.FloatAccountId)))
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
		// TODO: check on connect timeout and exclude from failed tx...
		transaction.Status = "FAILED"
		transaction, err = s.repository.UpdateTransaction(transaction)
		if err != nil {
			return nil, err
		}

		logger.ClientLog.Error("Error transferring float", "tx", transaction, "error", err)

		go func() {
			message := fmt.Sprintf("Sorry, KES%v Voucher transfer could not be processed", transaction.Amount)

			s.notifyApi.SendSMS("DEFAULT", merchant.Phone, message)
		}()

		return nil, err
	}

	s.paymentRepository.CreatePayment(&entities.Payment{
		Amount:        paymentData.Amount,
		Charge:        float32(paymentData.Charge),
		Status:        paymentData.Status,
		Description:   paymentData.Description,
		Destination:   paymentData.Destination,
		TransactionId: transaction.Id,
		PaymentId:     paymentData.Id,
	})

	//TODO: SMS
	if paymentData.Status == "COMPLETED" {
		// Update TX; after cashback?
		transaction, err := s.UpdateTransaction(&entities.Transaction{
			ModelID: entities.ModelID{Id: transaction.Id},
			Status:  paymentData.Status,
		})

		go func() {
			recipientAcc, _ := s.accountsApi.GetAccountById(strconv.Itoa(int(recipient.AccountId)))

			float, _ := s.paymentsApi.FetchFloatAccount(strconv.Itoa(int(merchant.FloatAccountId)))

			date := transaction.CreatedAt.Format("02/01/2006, 3:04 PM")

			// sender
			message := fmt.Sprintf("Voucher transfer of KES%v to %s on %s was successful. Cost KES%v. New Voucher Balance is KES%v",
				transaction.Amount, recipientAcc.Phone+" - "+recipient.BusinessName, date, paymentData.Charge, float.Balance)

			s.notifyApi.SendSMS("DEFAULT", merchant.Phone, message)

			float, _ = s.paymentsApi.FetchFloatAccount(strconv.Itoa(int(recipient.FloatAccountId)))

			// recipient
			message = fmt.Sprintf("You have received KES%v Voucher from %s on %s. New Voucher Balance is KES%v",
				transaction.Amount, merchant.Phone+" - "+merchant.BusinessName, date, float.Balance)

			s.notifyApi.SendSMS("DEFAULT", recipientAcc.Phone, message)

		}()

		return nil, err
	}

	return
}

func (s *service) FloatWithdraw(data *entities.Transaction, destination, account string) (transaction *entities.Transaction, err error) {
	merchant, err := s.merchantRepository.ReadMerchant(data.MerchantId)
	if err != nil {
		return nil, err
	}

	transaction, err = s.repository.CreateTransaction(data)
	if err != nil {
		return nil, err
	}

	paymentData, err := s.paymentsApi.FloatWithdraw(merchant.AccountId, merchant.FloatAccountId, int(transaction.Amount), destination, account)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
		// TODO: check on connect timeout and exclude from failed tx...
		transaction.Status = "FAILED"
		transaction, err = s.repository.UpdateTransaction(transaction)
		if err != nil {
			return nil, err
		}

		logger.ClientLog.Error("Error withdrawing float", "tx", transaction, "error", err)

		go func() {
			message := fmt.Sprintf("Sorry, KES%v Voucher withdrawal could not be processed", transaction.Amount)

			s.notifyApi.SendSMS("DEFAULT", merchant.Phone, message)
		}()

		return nil, err
	}

	s.paymentRepository.CreatePayment(&entities.Payment{
		Amount:        paymentData.Amount,
		Charge:        float32(paymentData.Charge),
		Status:        paymentData.Status,
		Description:   paymentData.Description,
		Destination:   paymentData.Destination,
		TransactionId: transaction.Id,
		PaymentId:     paymentData.Id,
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
		Charge: float32(paymentData.Charge),
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

func (s *service) WithdrawSavings(data *entities.Transaction, source, destination, account string) (tx *entities.Transaction, err error) {
	merchant, err := s.merchantRepository.ReadMerchant(data.MerchantId)
	if err != nil {
		return nil, err
	}

	if destination == "FLOAT" {
		return nil, pkg.ErrUnauthorized
	}

	personalAccounts, err := s.savingsApi.GetPersonalAccounts(strconv.Itoa(int(merchant.AccountId)))
	if err != nil {
		return nil, err
	}
	if len(personalAccounts) == 0 {
		return nil, pkg.ErrInvalidAccount
	}

	var personalAcc *clients.PersonalAccount
	for _, personalAccount := range personalAccounts {
		if personalAccount.Type == "MERCHANT_"+source {
			personalAcc = &personalAccount
			break
		}
	}
	if personalAcc == nil {
		return nil, pkg.ErrInvalidAccount
	}
	if personalAcc.Balance <= float64(data.Amount) {
		return nil, pkg.ErrInsufficientBalance
	}

	tx, err = s.repository.CreateTransaction(data)
	if err != nil {
		return nil, err
	}

	withdrawalData, err := s.savingsApi.WithdrawSavings(personalAcc.Id, destination, account, strconv.Itoa(int(tx.Id)), int(data.Amount))
	if err != nil {

		tx.Status = "FAILED"
		_, err := s.repository.UpdateTransaction(tx)

		return nil, err
	}

	_, err = s.savingsRepository.CreateSavingsTransaction(&entities.SavingsTransaction{
		Type:              withdrawalData.Type,
		Amount:            float32(withdrawalData.Amount),
		Status:            withdrawalData.Status,
		Description:       withdrawalData.Description,
		Extra:             withdrawalData.Extra,
		TransactionId:     tx.Id,
		SavingsId:         withdrawalData.Id,
		PersonalAccountId: withdrawalData.PersonalAccountId,
	})
	if err != nil {
		return nil, err
	}

	//if withdrawalData.Status == "COMPLETED" {
	//	err = s.CompleteTransaction(payment, paymentData)
	//
	//	return nil, err
	//}

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
	merchant, err := s.merchantRepository.ReadMerchant(transaction.MerchantId)

	switch transaction.Product {
	case consts.FLOAT_PURCHASE:
		if ipn.Status == "FAILED" {

			date := transaction.CreatedAt.Format("02/01/2006, 3:04 PM")

			message := ""

			if ipn.ErrorCode != 0 {
				message = fmt.Sprintf("Hi, we could not complete the"+
					" KES%v voucher purchase for %s on %s. %s. Please try again later.",
					transaction.Amount, *transaction.Destination, date, ipn.ErrorMessage)
			} else {
				message = fmt.Sprintf("Hi, we could not complete the"+
					" KES%v voucher purchase by %s on %s. Please try again later.",
					transaction.Amount, *transaction.Destination, date)
			}

			s.notifyApi.SendSMS("ERROR", merchant.Phone, message)

			transaction, err = s.UpdateTransaction(&entities.Transaction{
				ModelID: entities.ModelID{Id: transaction.Id},
				Status:  payment.Status,
			})

			return nil
		}

		go func() {
			float, _ := s.paymentsApi.FetchFloatAccount(strconv.Itoa(int(merchant.FloatAccountId)))

			date := transaction.CreatedAt.Format("02/01/2006, 3:04 PM")
			message := fmt.Sprintf("Ksh%v has been added to your merchant voucher account on %s via Mpesa. New balance is Ksh%v",
				transaction.Amount, date, float.Balance)

			s.notifyApi.SendSMS("DEFAULT", merchant.Phone, message)
		}()

	case consts.FLOAT_WITHDRAW:
		if ipn.Status == "FAILED" {

			date := transaction.CreatedAt.Format("02/01/2006, 3:04 PM")

			message := fmt.Sprintf("Hi, we could not complete the"+
				" KES%v voucher withdrawal by %s on %s. Please try again later.",
				transaction.Amount, *transaction.Destination, date)

			s.notifyApi.SendSMS("ERROR", merchant.Phone, message)

			transaction, err = s.UpdateTransaction(&entities.Transaction{
				ModelID: entities.ModelID{Id: transaction.Id},
				Status:  payment.Status,
			})

			return nil
		}

		go func() {
			float, _ := s.paymentsApi.FetchFloatAccount(strconv.Itoa(int(merchant.FloatAccountId)))

			date := transaction.CreatedAt.Format("02/01/2006, 3:04 PM")
			message := fmt.Sprintf("Voucher withdrawal of KES%v for %s on %s was successful. Cost KES%v. New Voucher Balance is KES%v",
				transaction.Amount, merchant.Phone, date, ipn.Charge, float.Balance)

			s.notifyApi.SendSMS("DEFAULT", merchant.Phone, message)
		}()

	case consts.CASH_WITHDRAW:
		if ipn.Status == "FAILED" {

			date := transaction.CreatedAt.Format("02/01/2006, 3:04 PM")

			message := fmt.Sprintf("Hi, we could not complete the"+
				" KES%v cash withdrawal by %s on %s. Please try again later.",
				transaction.Amount, *transaction.Destination, date)

			s.notifyApi.SendSMS("ERROR", merchant.Phone, message)

			transaction, err = s.UpdateTransaction(&entities.Transaction{
				ModelID: entities.ModelID{Id: transaction.Id},
				Status:  payment.Status,
			})

			return nil
		}

		err := s.computeMpesaWithdrawalCashback(merchant, transaction, payment, ipn)
		if err != nil {
			return err
		}

	case consts.MPESA_FLOAT:
		if ipn.Status == "FAILED" {

			date := transaction.CreatedAt.Format("02/01/2006, 3:04 PM")

			float, _ := s.paymentsApi.FetchFloatAccount(strconv.Itoa(int(merchant.FloatAccountId)))

			message := ""

			if ipn.ErrorCode != 0 {
				message = fmt.Sprintf("Hi, we could not complete your"+
					" KES%v float purchase for %s on %s. %s. Please try again later.",
					transaction.Amount, *transaction.Destination, date, ipn.ErrorMessage)
			} else {
				message = fmt.Sprintf("Hi, we have added KES%v to your voucher account "+
					"because we could not complete your"+
					" KES%v float purchase for %s on %s. New voucher balance is KES%v.",
					transaction.Amount, transaction.Amount, *transaction.Destination, date, float.Balance)
			}

			s.notifyApi.SendSMS("ERROR", merchant.Phone, message)

			transaction, err = s.UpdateTransaction(&entities.Transaction{
				ModelID: entities.ModelID{Id: transaction.Id},
				Status:  payment.Status,
			})

			return nil
		}

		err := s.computeCashback(merchant, transaction, payment, ipn)
		if err != nil {
			return err
		}

	case consts.EARNINGS_WITHDRAW:
		go func() {
			earningType := strings.Split(transaction.Description, " - ")[1]
			earningAcc, _ := s.earningAccRepository.ReadAccountByAccountIdAndType(merchant.AccountId, earningType)
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

			s.notifyApi.SendSMS("DEFAULT", merchant.Phone, message)
		}()
	}

	// Update TX; after cashback?
	transaction, err = s.UpdateTransaction(&entities.Transaction{
		ModelID: entities.ModelID{Id: transaction.Id},
		Status:  payment.Status,
	})

	return err
}

func (s *service) computeCashback(merchant *presenter.Merchant, tx *entities.Transaction, payment *entities.Payment, data *utils.Payment) (err error) {
	// Compute cashback and commissions
	// Compute cashback
	//TODO Fix this for float purchase using mpesa
	cashback := float32(30) * .2
	if payment.Charge == 0 {
		cashback = 0
	}

	if cashback > 0 {

		s.earningRepository.CreateEarning(&entities.Earning{
			Amount:        cashback,
			Type:          "SELF",
			TransactionId: tx.Id,
			AccountId:     merchant.AccountId,
		})

		earningAcc, err := s.earningAccRepository.ReadAccountByAccountIdAndType(merchant.AccountId, "CASHBACK")
		if err != nil {
			earningAcc, err = s.earningAccRepository.CreateAccount(&entities.EarningAccount{
				Type:      "CASHBACK",
				AccountId: merchant.AccountId,
			})
			if err != nil {
				return err
			}
		}
		s.earningAccService.CreditAccount(earningAcc.Id, cashback)
		s.earningAccService.DebitAccount(earningAcc.Id, cashback*.8) // Debit acc for savings

	}

	// Compute commissions
	//TODO Fix this for float purchase using mpesa
	commission := float32(30) * .1
	if payment.Charge == 0 {
		commission = 0
	}

	if commission > 0 {
		inviters, err := s.accountsApi.GetInviters(strconv.Itoa(int(merchant.AccountId)))
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
				s.earningAccService.DebitAccount(earningAcc.Id, commission*.8) // Debit acc for savings
			}
		}
	}

	go func() {
		date := tx.CreatedAt.Format("02/01/2006, 3:04 PM")

		destination := *tx.Destination
		if data.Store != "" && len(strings.Split(data.Store, " ")) != 1 {
			destination = strings.Join(strings.Split(data.Store, " ")[0:4], " ")
		}

		message := fmt.Sprintf("You have purchased KES%v float for %s on %s using Voucher. Cost KES%v. "+
			"You have received KES%v cashback.",
			payment.Amount,
			destination, date, data.Charge, cashback)

		if payment.Status != "COMPLETED" {
			message = fmt.Sprintf("Sorry, KES%v Float for %s could not be purchased", payment.Amount, *tx.Destination)
		}
		s.notifyApi.SendSMS("DEFAULT", merchant.Phone, message)
	}()

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		_, _ = s.mpesaStoreRepository.CreateStore(&entities.MpesaAgentStoreAccount{
			Agent:      strings.Split(*tx.Destination, "-")[0],
			Store:      strings.Split(*tx.Destination, "-")[1],
			Name:       strings.Join(strings.Split(data.Store, " ")[0:4], " "),
			MerchantId: merchant.Id,
		})
	}()

	// TODO: add go func with code to debit savings and send to save platform
	go s.earningService.SaveEarnings()

	return
}

func (s *service) computeMpesaWithdrawalCashback(merchant *presenter.Merchant, tx *entities.Transaction, payment *entities.Payment, data *utils.Payment) error {
	// Compute cashback and commissions
	// Compute cashback
	cashback := float32(s.getMpesaWithdrawalCashback(int(tx.Amount)))

	s.earningRepository.CreateEarning(&entities.Earning{
		Amount:        cashback,
		Type:          "SELF",
		TransactionId: tx.Id,
		AccountId:     merchant.AccountId,
	})

	earningAcc, err := s.earningAccRepository.ReadAccountByAccountIdAndType(merchant.AccountId, "COMMISSION")
	if err != nil {
		earningAcc, err = s.earningAccRepository.CreateAccount(&entities.EarningAccount{
			Type:      "COMMISSION",
			AccountId: merchant.AccountId,
		})
		if err != nil {
			return err
		}
	}
	s.earningAccService.CreditAccount(earningAcc.Id, cashback)
	s.earningAccService.DebitAccount(earningAcc.Id, cashback*.2) // Debit acc for savings

	// Compute commissions
	commission := float32(s.getMpesaWithdrawalCommission(int(tx.Amount)))

	inviters, err := s.accountsApi.GetInviters(strconv.Itoa(int(merchant.AccountId)))
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
			s.earningAccService.DebitAccount(earningAcc.Id, commission*.2) // Debit acc for savings
		}
	}

	go func() {
		float, _ := s.paymentsApi.FetchFloatAccount(strconv.Itoa(int(merchant.FloatAccountId)))

		date := tx.CreatedAt.Format("02/01/2006, 3:04 PM")

		message := fmt.Sprintf("KES%v cash withdrawal by %s on %s was successful. "+
			"New voucher balance KES%v. Commission earned KES%v. Commission saved KES%v",
			payment.Amount, *tx.Destination, date, float.Balance, cashback, cashback*.2)

		s.notifyApi.SendSMS("DEFAULT", merchant.Phone, message)
	}()

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

func (s *service) getMpesaWithdrawalCashback(amount int) int {
	cashbacks, err := s.paymentsApi.GetMpesaWithdrawalCommissions()
	if err != nil {
		return 0
	}

	for _, cashback := range cashbacks {
		if cashback.Min <= amount && amount <= cashback.Max {
			return cashback.Charge
		}
	}

	return 0
}

func (s *service) getMpesaWithdrawalCommission(amount int) int {
	commissions, err := s.paymentsApi.GetMpesaWithdrawalInviterCommissions()
	if err != nil {
		return 0
	}

	for _, commission := range commissions {
		if commission.Min <= amount && amount <= commission.Max {
			return commission.Charge
		}
	}

	return 0
}

func NewService(r Repository, merchantRepo merchant.Repository, paymentRepo payment.Repository, savingsRepo savings.Repository, earningAccRepo earning_account.Repository, earningRepo earning.Repository, mpesaStoreRepo mpesa_store.Repository, earningAccSrv earning_account.Service, earningSrv earning.Service) Service {
	return &service{
		repository: r,

		merchantRepository:   merchantRepo,
		paymentRepository:    paymentRepo,
		savingsRepository:    savingsRepo,
		earningAccRepository: earningAccRepo,
		earningRepository:    earningRepo,
		mpesaStoreRepository: mpesaStoreRepo,

		earningAccService: earningAccSrv,
		earningService:    earningSrv,

		accountsApi: clients.GetAccountClient(),
		paymentsApi: clients.GetPaymentClient(),
		savingsApi:  clients.GetSavingsClient(),
		notifyApi:   clients.GetNotifyClient(),
	}
}
