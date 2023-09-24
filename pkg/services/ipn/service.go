package ipn

import (
	"fmt"
	"merchants.sidooh/pkg/clients"
	"merchants.sidooh/pkg/entities"
	"merchants.sidooh/pkg/services/merchant"
	"merchants.sidooh/pkg/services/mpesa_store"
	"merchants.sidooh/pkg/services/payment"
	"merchants.sidooh/pkg/services/transaction"
	"merchants.sidooh/utils"
	"strconv"
	"strings"
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
	mpesaStoreRepository  mpesa_store.Repository
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

	// Update TX; after cashback?
	tx, err = s.transactionRepository.UpdateTransaction(&entities.Transaction{
		ModelID: entities.ModelID{tx.Id},
		Status:  payment.Status,
	})

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

	_, _ = s.mpesaStoreRepository.CreateStore(&entities.MpesaAgentStoreAccount{
		Agent:      strings.Split(tx.Destination, "-")[0],
		Store:      strings.Split(tx.Destination, "-")[1],
		Name:       strings.Join(strings.Split(data.Store, " ")[0:4], " "),
		MerchantId: mt.Id,
	})

	return err

}

func NewService(r payment.Repository, transactionRep transaction.Repository, merchantRep merchant.Repository, mpesaStoreRep mpesa_store.Repository) Service {
	return &service{paymentRepository: r, transactionRepository: transactionRep, merchantRepository: merchantRep, mpesaStoreRepository: mpesaStoreRep, notifyApi: clients.GetNotifyClient(), accountApi: clients.GetAccountClient()}
}
