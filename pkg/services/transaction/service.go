package transaction

import (
	"cmp"
	"fmt"
	"merchants.sidooh/api/presenter"
	"merchants.sidooh/pkg"
	"merchants.sidooh/pkg/clients"
	"merchants.sidooh/pkg/entities"
	"merchants.sidooh/pkg/services/earning_account"
	"merchants.sidooh/pkg/services/merchant"
	"merchants.sidooh/pkg/services/payment"
	"slices"
	"strconv"
)

type Service interface {
	FetchTransactions(filters Filters) (*[]presenter.Transaction, error)
	GetTransaction(id uint) (*presenter.Transaction, error)
	GetTransactionsByMerchant(merchantId uint) (*[]presenter.Transaction, error)
	CreateTransaction(transaction *entities.Transaction) (*entities.Transaction, error)
	PurchaseFloat(transaction *entities.Transaction, agent, store string) (*entities.Transaction, error)
	WithdrawEarnings(transaction *entities.Transaction, destination, account string) (*entities.Transaction, error)
	UpdateTransaction(transaction *entities.Transaction) (*presenter.Transaction, error)
}

type service struct {
	paymentsApi          *clients.ApiClient
	repository           Repository
	merchantRepository   merchant.Repository
	paymentRepository    payment.Repository
	earningAccRepository earning_account.Repository
	earningAccService    earning_account.Service
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

func (s *service) GetTransaction(id uint) (*presenter.Transaction, error) {
	return s.repository.ReadTransaction(id)
}

func (s *service) GetTransactionsByMerchant(merchantId uint) (*[]presenter.Transaction, error) {
	return s.repository.ReadTransactionsByMerchant(merchantId)
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
	if err != nil {
		tx.Status = "FAILED"
		s.repository.UpdateTransaction(tx)
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

func (s *service) WithdrawEarnings(data *entities.Transaction, destination, account string) (tx *entities.Transaction, err error) {
	merchant, err := s.merchantRepository.ReadMerchant(data.MerchantId)
	if err != nil {
		return nil, err
	}

	if destination == "FLOAT" && strconv.Itoa(int(merchant.FloatAccountId)) != account {
		return nil, pkg.ErrUnauthorized
	}

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

	//TODO: get actual charges here
	if totalBalance < data.Amount+15 {
		return nil, pkg.ErrInsufficientBalance
	}

	totalWithdrawal := data.Amount + 15

	var earningTXs []uint
	for _, earningAccount := range earningAccounts {
		toDebit := totalWithdrawal

		if earningAccount.Amount > totalWithdrawal {
			totalWithdrawal -= totalWithdrawal
		} else {
			totalWithdrawal -= earningAccount.Amount
			toDebit = earningAccount.Amount
		}

		tx, err := s.earningAccService.DebitAccount(earningAccount.Id, toDebit)
		if err != nil {
			return nil, err
		}
		earningTXs = append(earningTXs, tx.Id)

		if totalWithdrawal == 0 {
			break
		}
	}

	tx, err = s.repository.CreateTransaction(data)
	if err != nil {
		return nil, err
	}

	payment, err := s.paymentsApi.Withdraw(merchant.AccountId, 1, int(tx.Amount), destination, account)
	if err != nil {
		tx.Status = "FAILED"
		s.repository.UpdateTransaction(tx)

		// TODO: reverse earningTXs
		fmt.Println(earningTXs)
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

func (s *service) CreateTransaction(transaction *entities.Transaction) (*entities.Transaction, error) {
	return s.repository.CreateTransaction(transaction)
}

func (s *service) UpdateTransaction(transaction *entities.Transaction) (*presenter.Transaction, error) {
	return s.repository.UpdateTransaction(transaction)
}

func NewService(r Repository, merchantRepo merchant.Repository, paymentRepo payment.Repository, earningAccRepo earning_account.Repository, earningAccSrv earning_account.Service) Service {
	return &service{repository: r, paymentsApi: clients.GetPaymentClient(), merchantRepository: merchantRepo, paymentRepository: paymentRepo, earningAccRepository: earningAccRepo, earningAccService: earningAccSrv}
}
