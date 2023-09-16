package ipn

import (
	"fmt"
	"merchants.sidooh/pkg/clients"
	"merchants.sidooh/pkg/services/merchant"
	"merchants.sidooh/pkg/services/payment"
	"merchants.sidooh/pkg/services/transaction"
	"merchants.sidooh/utils"
	"strconv"
)

type Service interface {
	HandlePaymentIpn(data *utils.Payment) error
}

type service struct {
	notifyApi             *clients.ApiClient
	accountApi            *clients.ApiClient
	paymentRepository     payment.Repository
	transactionRepository transaction.Repository
	merchantRepository    merchant.Repository
}

func (s *service) HandlePaymentIpn(data *utils.Payment) error {
	payment, err := s.paymentRepository.ReadPaymentByColumn("payment_id", data.Id)
	if err != nil {
		return err
	}

	payment.Status = data.Status
	updatedPayment, err := s.paymentRepository.UpdatePayment(payment)

	tx, err := s.transactionRepository.ReadTransaction(updatedPayment.TransactionId)
	mt, err := s.merchantRepository.ReadMerchant(tx.MerchantId)

	account, err := s.accountApi.GetAccountById(strconv.Itoa(int(mt.AccountId)))
	if err != nil {
		return err
	}

	go func() {
		message := fmt.Sprintf("KES%v Float for %s purchased successfully", payment.Amount, tx.Destination)
		if payment.Status != "COMPLETED" {
			message = fmt.Sprintf("Sorry, KES%v Float for %s could not be purchased", payment.Amount, tx.Destination)
		}
		s.notifyApi.SendSMS("DEFAULT", account.Phone, message)
	}()

	return err

}

func NewService(r payment.Repository, transactionRep transaction.Repository, merchantRep merchant.Repository) Service {
	return &service{paymentRepository: r, transactionRepository: transactionRep, merchantRepository: merchantRep, notifyApi: clients.GetNotifyClient(), accountApi: clients.GetAccountClient()}
}
