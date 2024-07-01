package savings

import (
	"merchants.sidooh/pkg/services/merchant"
)

type Service interface {
	//FetchPayments() (*[]presenter.Payment, error)
	//GetPendingPayments() (*[]entities.Payment, error)
	//GetPayment(id uint) (*presenter.Payment, error)
	//CreatePayment(transaction *entities.SavingsTransaction) (*entities.SavingsTransaction, error)
	//UpdatePayment(payment *entities.Payment) (*presenter.Payment, error)
}

type service struct {
	//paymentsApi        *clients.ApiClient
	repository         Repository
	merchantRepository merchant.Repository
	paymentRepository  merchant.Repository
}

//
//func (s *service) FetchPayments() (*[]presenter.Payment, error) {
//	return s.repository.ReadPayments()
//}
//
//func (s *service) GetPayment(id uint) (*presenter.Payment, error) {
//	return s.repository.ReadPayment(id)
//}
//
//func (s *service) GetPendingPayments() (*[]entities.Payment, error) {
//	return s.repository.ReadPaymentsWhere("status", "PENDING")
//}

//func (s *service) CreatePayment(transaction *entities.SavingsTransaction) (*entities.SavingsTransaction, error) {
//	return s.repository.CreateSavingsTransaction(transaction)
//}

//
//func (s *service) UpdatePayment(payment *entities.Payment) (*presenter.Payment, error) {
//	return s.repository.UpdatePayment(payment)
//}

func NewService(r Repository, merchantRepo merchant.Repository) Service {
	return &service{repository: r, merchantRepository: merchantRepo}
}
