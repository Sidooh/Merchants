package earning

import (
	"fmt"
	"merchants.sidooh/pkg/clients"
	"merchants.sidooh/pkg/entities"
)

type Service interface {
	SaveEarnings() error
	CreateEarning(data *entities.Earning) (*entities.Earning, error)
}

type service struct {
	repository Repository
	savingsApi *clients.ApiClient
	notifyApi  *clients.ApiClient
}

func (s *service) SaveEarnings() error {
	earnings, err := s.repository.ReadPendingEarnings()

	savings := map[uint]clients.Investment{}
	for _, earning := range *earnings {
		inv := savings[earning.AccountId]
		inv.AccountId = earning.AccountId

		if earning.Type == "SELF" {
			inv.CashbackAmount += .2 * earning.Amount
		}

		if earning.Type == "INVITE" {
			inv.CommissionAmount += .2 * earning.Amount
		}

		savings[earning.AccountId] = inv
	}

	var investments []clients.Investment

	for _, investment := range savings {
		investments = append(investments, investment)
	}

	//message := "STATUS:SAVINGS\n\n"

	if len(*earnings) > 0 {
		savedEarnings, err := s.savingsApi.SaveEarnings(investments)
		if err != nil {
			return err
		}

		if _, ok := savedEarnings["completed"]; ok {
			fmt.Println(ok)
			//message += fmt.Sprintf("Processed earnings for %s'.count($completed)."  accounts\n");
		}
	} else {

	}

	//for _, earning := range *earnings {
	//	earning.Status = "COMPLETED"
	//	s.repository.UpdateEarning(&earning)
	//}

	return err
}

func (s *service) CreateEarning(data *entities.Earning) (*entities.Earning, error) {
	return s.repository.CreateEarning(data)
}

func NewService(r Repository) Service {
	return &service{repository: r, savingsApi: clients.GetSavingsClient(), notifyApi: clients.GetNotifyClient()}
}
