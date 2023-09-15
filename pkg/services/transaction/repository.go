package transaction

import (
	"merchants.sidooh/api/presenter"
	"merchants.sidooh/pkg/datastore"
	"merchants.sidooh/pkg/entities"
)

// Repository interface allows us to access the CRUD Operations here.
type Repository interface {
	CreateTransaction(transaction *entities.Transaction) (*entities.Transaction, error)
	ReadTransactions() (*[]presenter.Transaction, error)
	ReadTransaction(id uint) (*presenter.Transaction, error)
	ReadTransactionsByMerchant(merchantId uint) (*[]presenter.Transaction, error)
	UpdateTransaction(transaction *entities.Transaction) (*presenter.Transaction, error)
}
type repository struct {
}

func (r *repository) CreateTransaction(transaction *entities.Transaction) (*entities.Transaction, error) {
	result := datastore.DB.Create(&transaction)
	if result.Error != nil {
		return nil, result.Error
	}

	return transaction, nil
}

func (r *repository) ReadTransactions() (*[]presenter.Transaction, error) {
	var transactions []presenter.Transaction
	result := datastore.DB.Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}

	return &transactions, nil
}

func (r *repository) ReadTransaction(id uint) (*presenter.Transaction, error) {
	var transaction presenter.Transaction
	result := datastore.DB.First(&transaction, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &transaction, nil
}

func (r *repository) ReadTransactionsByMerchant(merchantId uint) (transaction *[]presenter.Transaction, err error) {
	err = datastore.DB.Where("merchant_id", merchantId).Find(&transaction).Error
	return
}

func (r *repository) UpdateTransaction(transaction *entities.Transaction) (*presenter.Transaction, error) {
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
