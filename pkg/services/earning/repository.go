package earning

import (
	"merchants.sidooh/pkg/datastore"
	"merchants.sidooh/pkg/entities"
)

// Repository interface allows us to access the CRUD Operations here.
type Repository interface {
	CreateEarning(data *entities.Earning) (*entities.Earning, error)
	ReadEarnings() (*[]entities.Earning, error)
	ReadPendingEarnings() (*[]entities.Earning, error)
}
type repository struct {
}

func (r *repository) CreateEarning(data *entities.Earning) (*entities.Earning, error) {
	result := datastore.DB.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}

	return data, nil
}

func (r *repository) ReadEarnings() (results *[]entities.Earning, err error) {
	err = datastore.DB.Find(&results).Error
	return
}

func (r *repository) ReadPendingEarnings() (results *[]entities.Earning, err error) {
	err = datastore.DB.Select("account_id, SUM(amount) amount").Where("status = ?", "PENDING").Group("account_id").Find(&results).Error
	return
}

// NewRepo is the single instance repo that is being created.
func NewRepo() Repository {
	return &repository{}
}
