package payment

import (
	"merchants.sidooh/api/presenter"
	"merchants.sidooh/pkg/datastore"
	"merchants.sidooh/pkg/entities"
)

// Repository interface allows us to access the CRUD Operations here.
type Repository interface {
	CreatePayment(payment *entities.Payment) (*entities.Payment, error)
	ReadPayments() (*[]presenter.Payment, error)
	ReadPayment(id uint) (*presenter.Payment, error)
	ReadPaymentByColumn(column string, value interface{}) (*entities.Payment, error)
	UpdatePayment(payment *entities.Payment) (*presenter.Payment, error)
}
type repository struct {
}

func (r *repository) CreatePayment(payment *entities.Payment) (*entities.Payment, error) {
	result := datastore.DB.Create(&payment)
	if result.Error != nil {
		return nil, result.Error
	}

	return payment, nil
}

func (r *repository) ReadPayments() (payments *[]presenter.Payment, err error) {
	err = datastore.DB.Find(&payments).Error
	return
}

func (r *repository) ReadPayment(id uint) (*presenter.Payment, error) {
	var payment presenter.Payment
	result := datastore.DB.First(&payment, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &payment, nil
}

func (r *repository) ReadPaymentByColumn(column string, value interface{}) (payment *entities.Payment, err error) {
	err = datastore.DB.Where(column, value).First(&payment).Error
	return
}

func (r *repository) UpdatePayment(payment *entities.Payment) (*presenter.Payment, error) {
	result := datastore.DB.Updates(payment)
	if result.Error != nil {
		return nil, result.Error
	}

	return r.ReadPayment(payment.Id)
}

// NewRepo is the single instance repo that is being created.
func NewRepo() Repository {
	return &repository{}
}
