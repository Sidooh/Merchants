package earning_account

import (
	"merchants.sidooh/api/presenter"
	"merchants.sidooh/pkg/datastore"
	"merchants.sidooh/pkg/entities"
)

// Repository interface allows us to access the CRUD Operations here.
type Repository interface {
	CreateAccount(data *entities.EarningAccount) (*entities.EarningAccount, error)
	ReadAccount(id uint) (*entities.EarningAccount, error)
	ReadAccountsByAccountId(accountId uint) (*[]presenter.EarningAccount, error)
	ReadAccountsByMerchant(merchantId uint) (*[]presenter.EarningAccount, error)
	ReadAccountByAccountIdAndType(accountId uint, accType string) (*entities.EarningAccount, error)
	UpdateAccount(data *entities.EarningAccount) (*entities.EarningAccount, error)
}
type repository struct {
}

func (r *repository) CreateAccount(data *entities.EarningAccount) (*entities.EarningAccount, error) {
	result := datastore.DB.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}

	return data, nil
}

func (r *repository) ReadAccount(id uint) (result *entities.EarningAccount, err error) {
	err = datastore.DB.First(&result, id).Error
	return
}

func (r *repository) ReadAccountsByAccountId(accountId uint) (result *[]presenter.EarningAccount, err error) {
	err = datastore.DB.Where("account_id", accountId).Find(&result).Error
	return
}

func (r *repository) ReadAccountsByMerchant(merchantId uint) (result *[]presenter.EarningAccount, err error) {
	merchant := entities.Merchant{}
	err = datastore.DB.First(&merchant, merchantId).Error
	if err != nil {
		return nil, err
	}

	err = datastore.DB.Where("account_id", merchant.AccountId).Find(&result).Error
	return
}

func (r *repository) ReadAccountByAccountIdAndType(accountId uint, accType string) (result *entities.EarningAccount, err error) {
	err = datastore.DB.Where("account_id", accountId).Where("type", accType).First(&result).Error
	return
}

func (r *repository) UpdateAccount(data *entities.EarningAccount) (*entities.EarningAccount, error) {
	result := datastore.DB.Updates(data)
	if result.Error != nil {
		return nil, result.Error
	}

	return r.ReadAccount(data.Id)
}

// NewRepo is the single instance repo that is being created.
func NewRepo() Repository {
	return &repository{}
}
