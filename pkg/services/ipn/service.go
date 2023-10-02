package ipn

import (
	"fmt"
	"merchants.sidooh/pkg/clients"
	"merchants.sidooh/pkg/entities"
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

	// Compute cashback and commissions
	// Compute cashback
	cashback := float32(data.Charge) * .2

	s.earningRepository.CreateEarning(&entities.Earning{
		Amount:        cashback,
		Type:          "SELF",
		TransactionId: tx.Id,
		AccountId:     mt.AccountId,
	})

	andType, err := s.earningAccountRepository.ReadAccountByAccountIdAndType(mt.AccountId, "CASHBACK")
	if err == nil {
		andType.Amount += cashback
		s.earningAccountRepository.UpdateAccount(andType)
	} else {
		_, err = s.earningAccountRepository.CreateAccount(&entities.EarningAccount{
			Type:      "CASHBACK",
			Amount:    cashback,
			AccountId: mt.AccountId,
		})
		if err != nil {
			return err
		}
	}

	// Compute commissions
	commsission := cashback / 2

	inviters, err := s.accountApi.GetInviters(strconv.Itoa(int(mt.AccountId)))
	if err != nil {
		return err
	}

	if len(inviters) > 1 {
		for _, inviter := range inviters[1:] {
			s.earningRepository.CreateEarning(&entities.Earning{
				Amount:        commsission,
				Type:          "INVITE",
				TransactionId: tx.Id,
				AccountId:     uint(inviter.Id),
			})

			andType, err := s.earningAccountRepository.ReadAccountByAccountIdAndType(uint(inviter.Id), "COMMISSION")
			if err == nil {
				andType.Amount += commsission
				s.earningAccountRepository.UpdateAccount(andType)
			} else {
				_, err = s.earningAccountRepository.CreateAccount(&entities.EarningAccount{
					Type:      "COMMISSION",
					Amount:    commsission,
					AccountId: uint(inviter.Id),
				})
				if err != nil {
					return err
				}
			}
		}
	}

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

func NewService(r payment.Repository, transactionRep transaction.Repository, merchantRep merchant.Repository, mpesaStoreRep mpesa_store.Repository, earningAccRep earning_account.Repository, earningRep earning.Repository) Service {
	return &service{paymentRepository: r,
		transactionRepository:    transactionRep,
		merchantRepository:       merchantRep,
		mpesaStoreRepository:     mpesaStoreRep,
		earningAccountRepository: earningAccRep,
		earningRepository:        earningRep,
		notifyApi:                clients.GetNotifyClient(),
		accountApi:               clients.GetAccountClient(),
	}
}
