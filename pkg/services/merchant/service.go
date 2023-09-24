package merchant

import (
	"merchants.sidooh/api/presenter"
	"merchants.sidooh/pkg"
	"merchants.sidooh/pkg/clients"
	"merchants.sidooh/pkg/entities"
	"merchants.sidooh/utils"
	"strconv"
)

type Service interface {
	FetchMerchants() (*[]presenter.Merchant, error)
	GetMerchant(id uint) (*presenter.Merchant, error)
	GetMerchantByAccount(accountId uint) (*presenter.Merchant, error)
	GetMerchantByIdNumber(idNumber string) (*presenter.Merchant, error)
	CreateMerchant(merchant *entities.Merchant) (*entities.Merchant, error)
	UpdateMerchantKYB(merchant *entities.Merchant) (*presenter.Merchant, error)
}

type service struct {
	paymentsApi *clients.ApiClient
	notifyApi   *clients.ApiClient
	accountApi  *clients.ApiClient
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

func (s *service) GetMerchantByIdNumber(idNumber string) (*presenter.Merchant, error) {
	return s.repository.ReadMerchantByIdNumber(idNumber)
}

func (s *service) CreateMerchant(data *entities.Merchant) (merchant *entities.Merchant, err error) {
	merchant, err = s.repository.CreateMerchant(data)

	account, err := s.accountApi.GetAccountById(strconv.Itoa(int(merchant.AccountId)))
	if err != nil {
		return nil, err
	}

	go s.notifyApi.SendSMS("DEFAULT", account.Phone, "KYC details created")

	return
}

func (s *service) UpdateMerchantKYB(data *entities.Merchant) (merchant *presenter.Merchant, err error) {
	merchant, err = s.repository.UpdateMerchant(data)
	if err != nil {
		return nil, pkg.ErrServerError
	}

	// TODO: Generate code and assign float account
	// TODO: Fix this to ensure uniqueness - get all codes and generate while comparing... or generate and check loop
	code := uint(utils.RandomIntBetween(10000, 99999))
	data.Code = &code

	floatAccount, err := s.paymentsApi.CreateFloatAccount(int(merchant.Id), int(merchant.AccountId))
	if err != nil {
		return nil, pkg.ErrServerError
	}
	id := uint(floatAccount.Id)
	data.FloatAccountId = &id

	merchant, err = s.repository.UpdateMerchant(data)
	if err != nil {
		return nil, err
	}

	account, err := s.accountApi.GetAccountById(strconv.Itoa(int(merchant.AccountId)))
	if err != nil {
		return nil, err
	}

	go s.notifyApi.SendSMS("DEFAULT", account.Phone, "KYB details updated")

	return
}

func NewService(r Repository) Service {
	return &service{repository: r, paymentsApi: clients.GetPaymentClient(), notifyApi: clients.GetNotifyClient(), accountApi: clients.GetAccountClient()}
}
