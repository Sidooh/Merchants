package earning

import (
	"merchants.sidooh/api/presenter"
	"merchants.sidooh/pkg/entities"
)

type Service interface {
	FetchEarnings() (*[]presenter.Earning, error)
	CreateEarning(data *entities.Earning) (*entities.Earning, error)
}

type service struct {
	repository Repository
}

func (s *service) FetchEarnings() (*[]presenter.Earning, error) {
	return s.repository.ReadEarnings()
}

func (s *service) CreateEarning(data *entities.Earning) (*entities.Earning, error) {
	return s.repository.CreateEarning(data)
}

func NewService(r Repository) Service {
	return &service{repository: r}
}
