package earning_account

import (
	"merchants.sidooh/api/presenter"
	"merchants.sidooh/pkg/datastore"
	"merchants.sidooh/pkg/entities"
)

// Repository interface allows us to access the CRUD Operations here.
type Repository interface {
	CreateAccount(store *entities.EarningAccount) (*entities.EarningAccount, error)
	ReadAccountsByMerchant(merchantId uint) (*[]presenter.EarningAccount, error)
}
type repository struct {
}

func (r *repository) CreateAccount(store *entities.EarningAccount) (*entities.EarningAccount, error) {
	result := datastore.DB.Create(&store)
	if result.Error != nil {
		return nil, result.Error
	}

	return store, nil
}

func (r *repository) ReadAccountsByMerchant(merchantId uint) (stores *[]presenter.EarningAccount, err error) {
	err = datastore.DB.Where("merchant_id", merchantId).Find(&stores).Error
	return
}

// NewRepo is the single instance repo that is being created.
func NewRepo() Repository {
	return &repository{}
}
