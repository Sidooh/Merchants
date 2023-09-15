package merchant

import (
	"merchants.sidooh/api/presenter"
	"merchants.sidooh/pkg/datastore"
	"merchants.sidooh/pkg/entities"
)

// Repository interface allows us to access the CRUD Operations here.
type Repository interface {
	CreateMerchant(merchant *entities.Merchant) (*entities.Merchant, error)
	ReadMerchants() (*[]presenter.Merchant, error)
	ReadMerchant(id uint) (*presenter.Merchant, error)
	ReadMerchantByAccount(accountId uint) (*presenter.Merchant, error)
	UpdateMerchant(merchant *entities.Merchant) (*presenter.Merchant, error)
}
type repository struct {
}

func (r *repository) CreateMerchant(merchant *entities.Merchant) (*entities.Merchant, error) {
	result := datastore.DB.Create(&merchant)
	if result.Error != nil {
		return nil, result.Error
	}

	return merchant, nil
}

func (r *repository) ReadMerchants() (*[]presenter.Merchant, error) {
	var merchants []presenter.Merchant
	result := datastore.DB.Find(&merchants)
	if result.Error != nil {
		return nil, result.Error
	}

	return &merchants, nil
}

func (r *repository) ReadMerchant(id uint) (*presenter.Merchant, error) {
	var merchant presenter.Merchant
	result := datastore.DB.First(&merchant, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &merchant, nil
}

func (r *repository) ReadMerchantByAccount(accountId uint) (merchant *presenter.Merchant, err error) {
	err = datastore.DB.Where("account_id", accountId).First(&merchant).Error
	return
}

func (r *repository) UpdateMerchant(merchant *entities.Merchant) (*presenter.Merchant, error) {
	result := datastore.DB.Updates(merchant)
	if result.Error != nil {
		return nil, result.Error
	}

	return r.ReadMerchant(merchant.Id)
}

// NewRepo is the single instance repo that is being created.
func NewRepo() Repository {
	return &repository{}
}
