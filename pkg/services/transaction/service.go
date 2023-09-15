package transaction

import (
	"github.com/gofiber/fiber/v2"
	"merchants.sidooh/api/presenter"
	"merchants.sidooh/pkg/entities"
)

type Service interface {
	FetchTransactions() (*[]presenter.Transaction, error)
	GetTransaction(id uint) (*presenter.Transaction, error)
	GetTransactionsByMerchant(merchantId uint) (*[]presenter.Transaction, error)
	CreateTransaction(transaction *entities.Transaction) (*entities.Transaction, error)
	UpdateTransaction(transaction *entities.Transaction) (*presenter.Transaction, error)
}

type service struct {
	apiClient  *fiber.Client
	repository Repository
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

func (s *service) CreateTransaction(transaction *entities.Transaction) (*entities.Transaction, error) {
	return s.repository.CreateTransaction(transaction)
}

func (s *service) UpdateTransaction(transaction *entities.Transaction) (*presenter.Transaction, error) {
	return s.repository.UpdateTransaction(transaction)
}

func NewService(r Repository) Service {
	return &service{repository: r}
}
