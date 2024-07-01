package savings

import (
	"merchants.sidooh/pkg/datastore"
	"merchants.sidooh/pkg/entities"
)

// Repository interface allows us to access the CRUD Operations here.
type Repository interface {
	CreateSavingsTransaction(transaction *entities.SavingsTransaction) (*entities.SavingsTransaction, error)
	//ReadPayments() (*[]presenter.Payment, error)
	//ReadPaymentsWhere(column string, value interface{}) (*[]entities.Payment, error)
	ReadTransaction(id uint) (transaction *entities.SavingsTransaction, err error)
	ReadTransactionByColumn(column string, value interface{}) (*entities.SavingsTransaction, error)
	UpdateTransaction(transaction *entities.SavingsTransaction) (*entities.SavingsTransaction, error)
}
type repository struct {
}

func (r *repository) CreateSavingsTransaction(transaction *entities.SavingsTransaction) (*entities.SavingsTransaction, error) {
	result := datastore.DB.Create(&transaction)
	if result.Error != nil {
		return nil, result.Error
	}

	return transaction, nil
}

//	func (r *repository) ReadPayments() (payments *[]presenter.Payment, err error) {
//		err = datastore.DB.Find(&payments).Error
//		return
//	}
//
//	func (r *repository) ReadPaymentsWhere(column string, value interface{}) (payments *[]entities.Payment, err error) {
//		err = datastore.DB.Where(column, value).Find(&payments).Error
//		return
//	}
func (r *repository) ReadTransaction(id uint) (transaction *entities.SavingsTransaction, err error) {
	err = datastore.DB.First(&transaction, id).Error
	return
}

func (r *repository) ReadTransactionByColumn(column string, value interface{}) (transaction *entities.SavingsTransaction, err error) {
	err = datastore.DB.Where(column, value).First(&transaction).Error
	return
}

func (r *repository) UpdateTransaction(transaction *entities.SavingsTransaction) (*entities.SavingsTransaction, error) {
	result := datastore.DB.Updates(transaction)
	if result.Error != nil {
		return nil, result.Error
	}

	return r.ReadTransaction(transaction.Id)
}

// NewRepo is the single instance repo that is being created.
func NewRepo() Repository {
	return &repository{}
}
