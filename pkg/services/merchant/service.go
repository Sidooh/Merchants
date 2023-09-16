package merchant

import (
	"merchants.sidooh/api/presenter"
	"merchants.sidooh/pkg/clients"
	"merchants.sidooh/pkg/entities"
)

type Service interface {
	FetchMerchants() (*[]presenter.Merchant, error)
	GetMerchant(id uint) (*presenter.Merchant, error)
	GetMerchantByAccount(accountId uint) (*presenter.Merchant, error)
	CreateMerchant(merchant *entities.Merchant) (*entities.Merchant, error)
	UpdateMerchant(merchant *entities.Merchant) (*presenter.Merchant, error)
}

type service struct {
	paymentsApi *clients.ApiClient
	repository  Repository
}

func (s *service) FetchMerchants() (*[]presenter.Merchant, error) {
	return s.repository.ReadMerchants()
}

func (s *service) GetMerchant(id uint) (*presenter.Merchant, error) {
	return s.repository.ReadMerchant(id)
}

func (s *service) GetMerchantByAccount(accountId uint) (*presenter.Merchant, error) {
	return s.repository.ReadMerchantByAccount(accountId)
}

func (s *service) CreateMerchant(data *entities.Merchant) (merchant *entities.Merchant, err error) {
	merchant, err = s.repository.CreateMerchant(data)

	// TODO: Generate code and assign float account
	//floatAccount, err := s.paymentsApi.CreateFloatAccount(int(merchant.Id), int(merchant.AccountId))
	//if err != nil {
	//	return nil, pkg.ErrServerError
	//}
	//
	//id := uint(floatAccount.Id)
	//merchant.FloatAccountId = &id
	//merchant, err = s.repository.UpdateMerchant(merchant)
	//if err != nil {
	//	return nil, err
	//}

	return
}

func (s *service) UpdateMerchant(merchant *entities.Merchant) (*presenter.Merchant, error) {
	return s.repository.UpdateMerchant(merchant)
}

func NewService(r Repository) Service {
	return &service{repository: r, paymentsApi: clients.GetPaymentClient()}
}
