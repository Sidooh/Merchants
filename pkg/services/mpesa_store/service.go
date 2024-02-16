package mpesa_store

import (
	"merchants.sidooh/api/presenter"
	"merchants.sidooh/pkg/entities"
)

type Service interface {
	FetchStoresByMerchant(merchantId uint) (*[]presenter.MpesaAgentStoreAccount, error)
	CreateStore(store *entities.MpesaAgentStoreAccount) (*entities.MpesaAgentStoreAccount, error)
}

type service struct {
	repository Repository
}

func (s *service) FetchStoresByMerchant(merchantId uint) (*[]presenter.MpesaAgentStoreAccount, error) {
	return s.repository.ReadStoresByMerchant(merchantId)
}

func (s *service) CreateStore(store *entities.MpesaAgentStoreAccount) (*entities.MpesaAgentStoreAccount, error) {
	return s.repository.CreateStore(store)
}

func NewService(r Repository) Service {
	return &service{repository: r}
}
