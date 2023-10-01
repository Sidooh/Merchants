package earning_account

import (
	"merchants.sidooh/api/presenter"
	"merchants.sidooh/pkg/entities"
)

type Service interface {
	FetchAccountsByMerchant(merchantId uint) (*[]presenter.EarningAccount, error)
	CreateAccount(store *entities.EarningAccount) (*entities.EarningAccount, error)
}

type service struct {
	repository Repository
}

func (s *service) FetchAccountsByMerchant(merchantId uint) (*[]presenter.EarningAccount, error) {
	return s.repository.ReadAccountsByMerchant(merchantId)
}

func (s *service) CreateAccount(store *entities.EarningAccount) (*entities.EarningAccount, error) {
	return s.repository.CreateAccount(store)
}

func NewService(r Repository) Service {
	return &service{repository: r}
}
