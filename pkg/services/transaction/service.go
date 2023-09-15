package transaction

import (
	"merchants.sidooh/api/presenter"
	"merchants.sidooh/pkg/clients"
	"merchants.sidooh/pkg/entities"
	"merchants.sidooh/pkg/services/merchant"
	"merchants.sidooh/pkg/services/payment"
)

type Service interface {
	FetchTransactions() (*[]presenter.Transaction, error)
	GetTransaction(id uint) (*presenter.Transaction, error)
	GetTransactionsByMerchant(merchantId uint) (*[]presenter.Transaction, error)
	CreateTransaction(transaction *entities.Transaction) (*entities.Transaction, error)
	PurchaseFloat(transaction *entities.Transaction, agent, store string) (*entities.Transaction, error)
	UpdateTransaction(transaction *entities.Transaction) (*presenter.Transaction, error)
}

type service struct {
	paymentsApi        *clients.ApiClient
	repository         Repository
	merchantRepository merchant.Repository
	paymentRepository  payment.Repository
}

func (s *service) FetchTransactions() (*[]presenter.Transaction, error) {
	return s.repository.ReadTransactions()
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

func (s *service) CreateTransaction(transaction *entities.Transaction) (*entities.Transaction, error) {
	return s.repository.CreateTransaction(transaction)
}

func (s *service) UpdateTransaction(transaction *entities.Transaction) (*presenter.Transaction, error) {
	return s.repository.UpdateTransaction(transaction)
}

func NewService(r Repository, merchantRepo merchant.Repository, paymentRepo payment.Repository) Service {
	return &service{repository: r, paymentsApi: clients.GetPaymentClient(), merchantRepository: merchantRepo, paymentRepository: paymentRepo}
}
