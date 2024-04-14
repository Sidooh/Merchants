package mpesa_store

import (
	"merchants.sidooh/api/presenter"
	"merchants.sidooh/pkg/datastore"
	"merchants.sidooh/pkg/entities"
)

// Repository interface allows us to access the CRUD Operations here.
type Repository interface {
	CreateStore(store *entities.MpesaAgentStoreAccount) (*entities.MpesaAgentStoreAccount, error)
	ReadStoresByMerchant(merchantId uint) (*[]presenter.MpesaAgentStoreAccount, error)
	ReadAllStores() ([]*presenter.MpesaStore, error)
}
type repository struct {
}

func (r *repository) CreateStore(store *entities.MpesaAgentStoreAccount) (*entities.MpesaAgentStoreAccount, error) {
	result := datastore.DB.Create(&store)
	if result.Error != nil {
		return nil, result.Error
	}

	return store, nil
}

func (r *repository) ReadStoresByMerchant(merchantId uint) (stores *[]presenter.MpesaAgentStoreAccount, err error) {
	err = datastore.DB.Where("merchant_id", merchantId).Find(&stores).Error
	return
}

func (r *repository) ReadAllStores() (stores []*presenter.MpesaStore, err error) {
	err = datastore.DB.Model(entities.MpesaAgentStoreAccount{}).
		Select("agent", "store", "name").
		Distinct("agent", "store", "name").
		Find(&stores).
		Error
	return
}

// NewRepo is the single instance repo that is being created.
func NewRepo() Repository {
	return &repository{}
}
