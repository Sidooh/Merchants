package jobs

import (
	"fmt"
	"merchants.sidooh/pkg/clients"
	"merchants.sidooh/pkg/logger"
	"merchants.sidooh/pkg/services/earning"
	"merchants.sidooh/pkg/services/payment"
	"merchants.sidooh/pkg/services/transaction"
	"strconv"
)

type Service interface {
	EarningsInvestments() error
	QueryPaymentsStatus() error
}

type service struct {
	earningService     earning.Service
	paymentService     payment.Service
	transactionService transaction.Service

	paymentsApi *clients.ApiClient
	notifyApi   *clients.ApiClient
}

func (s *service) EarningsInvestments() error {
	go func() {
		err := s.earningService.SaveEarnings()
		if err != nil {
			message := fmt.Sprintf("Failed to save process merchant earnings")
			logger.ClientLog.Error(message, err, err)

			_ = s.notifyApi.SendSMS("DEFAULT", "0780611696", message)
		}

	}()

	return nil
}

func (s *service) QueryPaymentsStatus() error {
	payments, err := s.paymentService.GetPendingPayments()
	if err != nil {
		return err
	}

	go func() {
		for _, payment := range *payments {
			paymentData, err := s.paymentsApi.Find(strconv.Itoa(int(payment.PaymentId)))
			if err != nil {
				logger.ClientLog.Error("failed to fetch payment", err)
			}

			if paymentData != nil && paymentData.Status != "PENDING" {

				err := s.transactionService.CompleteTransaction(&payment, paymentData)
				if err != nil {
					logger.ClientLog.Error("failed to complete transaction", err)
				}

			}
		}
	}()

	return nil
}

func NewService(earningSrv earning.Service, paymentSrv payment.Service, transactionSrv transaction.Service) Service {
	return &service{
		earningService:     earningSrv,
		paymentService:     paymentSrv,
		transactionService: transactionSrv,

		paymentsApi: clients.GetPaymentClient(),
		notifyApi:   clients.GetNotifyClient(),
	}
}
