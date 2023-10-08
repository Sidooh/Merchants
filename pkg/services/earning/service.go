package earning

import (
	"fmt"
	"merchants.sidooh/api/presenter"
	"merchants.sidooh/pkg/clients"
	"merchants.sidooh/pkg/entities"
	"merchants.sidooh/utils"
)

type Service interface {
	FetchEarnings() (*[]presenter.Earning, error)
	//FetchPendingEarnings() (*[]presenter.Earning, error)
	SaveEarnings() (interface{}, error)
	CreateEarning(data *entities.Earning) (*entities.Earning, error)
}

type service struct {
	repository Repository
	savingsApi *clients.ApiClient
}

func (s *service) FetchEarnings() (results *[]presenter.Earning, err error) {
	earnings, err := s.repository.ReadEarnings()
	if err != nil {
		return nil, err
	}
	utils.ConvertStruct(earnings, results)

	return
}

//func (s *service) FetchPendingEarnings() (results *[]presenter.Earning, err error) {
//	earnings, err := s.repository.ReadPendingEarnings()
//	if err != nil {
//		return nil, err
//	}
//	utils.ConvertStruct(earnings, results)
//
//	return
//}

func (s *service) SaveEarnings() (interface{}, error) {
	earnings, err := s.repository.ReadPendingEarnings()

	savings := map[int]int{}
	for _, earning := range *earnings {
		savings[int(earning.AccountId)] = int(earning.Amount)
	}

	saveEarnings, err := s.savingsApi.SaveEarnings()
	fmt.Println(saveEarnings, err)
	if err != nil {
		return nil, err
	}

	return savings, err
}

func (s *service) CreateEarning(data *entities.Earning) (*entities.Earning, error) {
	return s.repository.CreateEarning(data)
}

func NewService(r Repository) Service {
	return &service{repository: r, savingsApi: clients.GetSavingsClient()}
}
