package earning_account

import (
	"merchants.sidooh/api/presenter"
	"merchants.sidooh/pkg/entities"
)

type Service interface {
	FetchAccountsByAccountId(accountId uint) (*[]presenter.EarningAccount, error)
	FetchAccountsByMerchant(merchantId uint) (*[]presenter.EarningAccount, error)
	CreateAccount(data *entities.EarningAccount) (*entities.EarningAccount, error)
	UpdateAccount(data *entities.EarningAccount) (*entities.EarningAccount, error)
}

type service struct {
	repository Repository
}

func (s *service) FetchAccountsByAccountId(accountId uint) (*[]presenter.EarningAccount, error) {
	return s.repository.ReadAccountsByAccountId(accountId)
}

func (s *service) FetchAccountsByMerchant(merchantId uint) (*[]presenter.EarningAccount, error) {
	return s.repository.ReadAccountsByMerchant(merchantId)
}

func (s *service) CreateAccount(data *entities.EarningAccount) (*entities.EarningAccount, error) {
	return s.repository.CreateAccount(data)
}

func (s *service) UpdateAccount(data *entities.EarningAccount) (*entities.EarningAccount, error) {
	return s.repository.UpdateAccount(data)
}

func NewService(r Repository) Service {
	return &service{repository: r}
}
