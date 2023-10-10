package ipn

import (
	"fmt"
	"merchants.sidooh/api/presenter"
	"merchants.sidooh/pkg/clients"
	"merchants.sidooh/pkg/entities"
	"merchants.sidooh/pkg/logger"
	"merchants.sidooh/pkg/services/earning"
	"merchants.sidooh/pkg/services/earning_account"
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
	notifyApi                *clients.ApiClient
	accountApi               *clients.ApiClient
	paymentRepository        payment.Repository
	transactionRepository    transaction.Repository
	merchantRepository       merchant.Repository
	mpesaStoreRepository     mpesa_store.Repository
	earningAccountRepository earning_account.Repository
	earningRepository        earning.Repository
	earningAccountService    earning_account.Service
	earningService           earning.Service
}

func (s *service) HandlePaymentIpn(data *utils.Payment) error {
	payment, err := s.paymentRepository.ReadPaymentByColumn("payment_id", data.Id)
	if err != nil {
		return err
	}

	if payment.Status != "PENDING" {
		go func() {
			message := fmt.Sprintf("Merchant Payment is not pending, check %v", payment.Id)
			logger.ClientLog.Error(message, "payment", payment)

			s.notifyApi.SendSMS("DEFAULT", "0780611696", message)
		}()

		return nil
	}

	payment.Status = data.Status
	updatedPayment, err := s.paymentRepository.UpdatePayment(payment)

	tx, err := s.transactionRepository.ReadTransaction(updatedPayment.TransactionId)
	mt, err := s.merchantRepository.ReadMerchant(tx.MerchantId)

	// Update TX; after cashback?
	tx, err = s.transactionRepository.UpdateTransaction(&entities.Transaction{
		ModelID: entities.ModelID{Id: tx.Id},
		Status:  payment.Status,
	})

	switch tx.Product {
	case "FLOAT":
		return s.computeCashback(mt, tx, payment, data)

	case "WITHDRAWAL":
		account, err := s.accountApi.GetAccountById(strconv.Itoa(int(mt.AccountId)))
		if err != nil {
			return err
		}

		go func() {
			message := fmt.Sprintf("KES%v Withdrawal to %s was successful", payment.Amount, tx.Destination)
			if payment.Status != "COMPLETED" {
				message = fmt.Sprintf("Sorry, KES%v Withdrawal to %s could not be processed", payment.Amount, tx.Destination)
			}
			s.notifyApi.SendSMS("DEFAULT", account.Phone, message)
		}()
	}

	return err

}

func (s *service) computeCashback(mt *presenter.Merchant, tx *presenter.Transaction, payment *entities.Payment, data *utils.Payment) error {

	// Compute cashback and commissions
	// Compute cashback
	cashback := float32(data.Charge) * .2

	s.earningRepository.CreateEarning(&entities.Earning{
		Amount:        cashback,
		Type:          "SELF",
		TransactionId: tx.Id,
		AccountId:     mt.AccountId,
	})

	earningAcc, err := s.earningAccountRepository.ReadAccountByAccountIdAndType(mt.AccountId, "CASHBACK")
	if err != nil {
		earningAcc, err = s.earningAccountRepository.CreateAccount(&entities.EarningAccount{
			Type:      "CASHBACK",
			AccountId: mt.AccountId,
		})
		if err != nil {
			return err
		}
	}
	s.earningAccountService.CreditAccount(earningAcc.AccountId, cashback)
	//s.earningAccountService.DebitAccount(earningAcc.AccountId, cashback) // Debit acc for savings

	// Compute commissions
	commission := float32(data.Charge) * .1

	inviters, err := s.accountApi.GetInviters(strconv.Itoa(int(mt.AccountId)))
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

			earningAcc, err := s.earningAccountRepository.ReadAccountByAccountIdAndType(uint(inviter.Id), "COMMISSION")
			if err != nil {
				earningAcc, err = s.earningAccountRepository.CreateAccount(&entities.EarningAccount{
					Type:      "COMMISSION",
					AccountId: uint(inviter.Id),
				})
				if err != nil {
					return err
				}
			}
			s.earningAccountService.CreditAccount(earningAcc.AccountId, commission)

		}
	}

	account, err := s.accountApi.GetAccountById(strconv.Itoa(int(mt.AccountId)))
	if err != nil {
		return err
	}

	go func() {
		message := fmt.Sprintf("KES%v Float for %s purchased successfully", payment.Amount, strings.Join(strings.Split(data.Store, " ")[0:4], " "))
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

	// TODO: add go func with code to debit savings and send to save platform

	return err
}

func NewService(r payment.Repository, transactionRep transaction.Repository, merchantRep merchant.Repository, mpesaStoreRep mpesa_store.Repository, earningAccRep earning_account.Repository, earningRep earning.Repository, earningAccSrv earning_account.Service, earningSrv earning.Service) Service {
	return &service{paymentRepository: r,
		transactionRepository:    transactionRep,
		merchantRepository:       merchantRep,
		mpesaStoreRepository:     mpesaStoreRep,
		earningAccountRepository: earningAccRep,
		earningRepository:        earningRep,
		earningAccountService:    earningAccSrv,
		notifyApi:                clients.GetNotifyClient(),
		accountApi:               clients.GetAccountClient(),
		earningService:           earningSrv,
	}
}
