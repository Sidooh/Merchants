package merchant

import (
	"github.com/gofiber/fiber/v2"
	"merchants.sidooh/api/presenter"
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
	apiClient  *fiber.Client
	repository Repository
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

func (s *service) CreateMerchant(merchant *entities.Merchant) (*entities.Merchant, error) {
	return s.repository.CreateMerchant(merchant)
}

func (s *service) UpdateMerchant(merchant *entities.Merchant) (*presenter.Merchant, error) {
	return s.repository.UpdateMerchant(merchant)
}

func NewService(r Repository) Service {
	return &service{repository: r}
}
