package earning_account_transaction

import (
	"merchants.sidooh/pkg/datastore"
	"merchants.sidooh/pkg/entities"
)

// Repository interface allows us to access the CRUD Operations here.
type Repository interface {
	CreateTransaction(data *entities.EarningAccountTransaction) (*entities.EarningAccountTransaction, error)
}
type repository struct {
}

func (r *repository) CreateTransaction(data *entities.EarningAccountTransaction) (*entities.EarningAccountTransaction, error) {
	result := datastore.DB.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}

	return data, nil
}

// NewRepo is the single instance repo that is being created.
func NewRepo() Repository {
	return &repository{}
}
