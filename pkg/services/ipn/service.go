package ipn

import (
	"merchants.sidooh/pkg/clients"
	"merchants.sidooh/pkg/services/payment"
	"merchants.sidooh/utils"
)

type Service interface {
	HandlePaymentIpn(data *utils.Payment) error
}

type service struct {
	paymentsApi       *clients.ApiClient
	paymentRepository payment.Repository
}

func (s *service) HandlePaymentIpn(data *utils.Payment) error {
	payment, err := s.paymentRepository.ReadPaymentByColumn("payment_id", data.Id)
	if err != nil {
		return err
	}

	payment.Status = data.Status
	_, err = s.paymentRepository.UpdatePayment(payment)
	return err

}

func NewService(r payment.Repository) Service {
	return &service{paymentRepository: r}
}
