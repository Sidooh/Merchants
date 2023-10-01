package earning

import (
	"merchants.sidooh/api/presenter"
	"merchants.sidooh/pkg/datastore"
	"merchants.sidooh/pkg/entities"
)

// Repository interface allows us to access the CRUD Operations here.
type Repository interface {
	CreateEarning(data *entities.Earning) (*entities.Earning, error)
	ReadEarnings() (*[]presenter.Earning, error)
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

func (r *repository) ReadEarnings() (results *[]presenter.Earning, err error) {
	err = datastore.DB.Find(&results).Error
	return
}

// NewRepo is the single instance repo that is being created.
func NewRepo() Repository {
	return &repository{}
}
