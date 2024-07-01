package ipn

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"merchants.sidooh/pkg/clients"
	"merchants.sidooh/pkg/entities"
	"merchants.sidooh/pkg/logger"
	"merchants.sidooh/pkg/services/earning"
	"merchants.sidooh/pkg/services/earning_account"
	"merchants.sidooh/pkg/services/merchant"
	"merchants.sidooh/pkg/services/mpesa_store"
	"merchants.sidooh/pkg/services/payment"
	"merchants.sidooh/pkg/services/savings"
	"merchants.sidooh/pkg/services/transaction"
	"merchants.sidooh/utils"
	"strings"
)

type Service interface {
	HandlePaymentIpn(data *utils.Payment) error
	HandleSavingsIpn(data map[int]string) error
}

type service struct {
	notifyApi                *clients.ApiClient
	accountApi               *clients.ApiClient
	paymentRepository        payment.Repository
	savingsRepository        savings.Repository
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

func (s *service) HandleSavingsIpn(data map[int]string) error {

	for id, res := range data {
		go func() {
			tx, err := s.savingsRepository.ReadTransactionByColumn("savings_id", id)
			if err != nil {
				log.Error(err)
			}

			if tx.Status == "PENDING" && (res == "COMPLETED" || res == "FAILED") {
				tx, err = s.savingsRepository.UpdateTransaction(&entities.SavingsTransaction{
					ModelID: entities.ModelID{Id: tx.Id},
					Status:  res,
				})

				t, _ := s.transactionRepository.UpdateTransaction(&entities.Transaction{
					ModelID: entities.ModelID{Id: tx.TransactionId},
					Status:  res,
				})

				merchant, _ := s.merchantRepository.ReadMerchant(t.MerchantId)

				// Send Notification
				accType := strings.Split(t.Description, " - ")[1]
				destination := *t.Destination
				if strings.Split(*t.Destination, "-")[0] == "FLOAT" {
					destination = "VOUCHER"
				}
				date := tx.CreatedAt.Format("02/01/2006, 3:04 PM")

				message := fmt.Sprintf("Withdrawal of KES%v from savings %s to %s on %s was successful.",
					tx.Amount, accType, destination, date)
				if t.Status == "FAILED" {
					message = fmt.Sprintf("Sorry, KES%v Withdrawal to %s could not be processed", t.Amount, destination)
				}

				s.notifyApi.SendSMS("DEFAULT", merchant.Phone, message)

			}

		}()
	}

	return nil
}

func NewService(r payment.Repository, savingsRep savings.Repository, transactionRep transaction.Repository, merchantRep merchant.Repository, mpesaStoreRep mpesa_store.Repository, earningAccRep earning_account.Repository, earningRep earning.Repository, transactionSrv transaction.Service, earningAccSrv earning_account.Service, earningSrv earning.Service) Service {
	return &service{
		paymentRepository:        r,
		savingsRepository:        savingsRep,
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
