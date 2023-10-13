package ipn

import (
	"fmt"
	"merchants.sidooh/pkg/clients"
	"merchants.sidooh/pkg/logger"
	"merchants.sidooh/pkg/services/earning"
	"merchants.sidooh/pkg/services/earning_account"
	"merchants.sidooh/pkg/services/merchant"
	"merchants.sidooh/pkg/services/mpesa_store"
	"merchants.sidooh/pkg/services/payment"
	"merchants.sidooh/pkg/services/transaction"
	"merchants.sidooh/utils"
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
	transactionService       transaction.Service
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

	err = s.transactionService.CompleteTransaction(payment, data)
	if err != nil {
		return err
	}

	return err

}

func NewService(r payment.Repository, transactionRep transaction.Repository, merchantRep merchant.Repository, mpesaStoreRep mpesa_store.Repository, earningAccRep earning_account.Repository, earningRep earning.Repository, transactionSrv transaction.Service, earningAccSrv earning_account.Service, earningSrv earning.Service) Service {
	return &service{paymentRepository: r,
		transactionRepository:    transactionRep,
		merchantRepository:       merchantRep,
		mpesaStoreRepository:     mpesaStoreRep,
		earningAccountRepository: earningAccRep,
		earningRepository:        earningRep,
		transactionService:       transactionSrv,
		earningAccountService:    earningAccSrv,
		earningService:           earningSrv,
		notifyApi:                clients.GetNotifyClient(),
		accountApi:               clients.GetAccountClient(),
	}
}
