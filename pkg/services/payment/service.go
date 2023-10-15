package payment

import (
	"merchants.sidooh/api/presenter"
	"merchants.sidooh/pkg/clients"
	"merchants.sidooh/pkg/entities"
	"merchants.sidooh/pkg/services/merchant"
)

type Service interface {
	FetchPayments() (*[]presenter.Payment, error)
	GetPayment(id uint) (*presenter.Payment, error)
	CreatePayment(payment *entities.Payment) (*entities.Payment, error)
	UpdatePayment(payment *entities.Payment) (*presenter.Payment, error)
}

type service struct {
	paymentsApi        *clients.ApiClient
	repository         Repository
	merchantRepository merchant.Repository
	paymentRepository  merchant.Repository
}

func (s *service) FetchPayments() (*[]presenter.Payment, error) {
	return s.repository.ReadPayments()
}

func (s *service) GetPayment(id uint) (*presenter.Payment, error) {
	return s.repository.ReadPayment(id)
}

func (s *service) CreatePayment(payment *entities.Payment) (*entities.Payment, error) {
	return s.repository.CreatePayment(payment)
}

func (s *service) UpdatePayment(payment *entities.Payment) (*presenter.Payment, error) {
	return s.repository.UpdatePayment(payment)
}

func NewService(r Repository, merchantRepo merchant.Repository) Service {
	return &service{repository: r, paymentsApi: clients.GetPaymentClient(), merchantRepository: merchantRepo}
}
