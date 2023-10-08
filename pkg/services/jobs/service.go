package jobs

import (
	"fmt"
	"merchants.sidooh/pkg/clients"
	"merchants.sidooh/pkg/logger"
	"merchants.sidooh/pkg/services/earning"
)

type Service interface {
	EarningsInvestments() error
}

type service struct {
	earningService earning.Service
	notifyApi      *clients.ApiClient
}

func (s *service) EarningsInvestments() error {
	go func() {
		_, err := s.earningService.SaveEarnings()
		if err != nil {
			message := fmt.Sprintf("Failed to save process merchant earnings")
			logger.ClientLog.Error(message, err, err)

			s.notifyApi.SendSMS("DEFAULT", "0780611696", message)
		}
	}()

	return nil
}

func NewService(earningSrv earning.Service) Service {
	return &service{
		earningService: earningSrv,
		notifyApi:      clients.GetNotifyClient(),
	}
}
