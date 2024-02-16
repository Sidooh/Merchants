package transaction

import (
	"merchants.sidooh/api/presenter"
	"merchants.sidooh/pkg/datastore"
	"merchants.sidooh/pkg/entities"
	"time"
)

// Repository interface allows us to access the CRUD Operations here.
type Repository interface {
	CreateTransaction(transaction *entities.Transaction) (*entities.Transaction, error)
	ReadTransactions(filters Filters) (*[]presenter.Transaction, error)
	ReadTransaction(id uint) (*entities.Transaction, error)
	ReadTransactionsByMerchant(merchantId uint) (*[]presenter.Transaction, error)
	UpdateTransaction(transaction *entities.Transaction) (*entities.Transaction, error)
}
type repository struct {
}

type Filters struct {
	Columns []string

	Accounts  []string
	Merchants []string
	Days      int
}

func (r *repository) CreateTransaction(transaction *entities.Transaction) (*entities.Transaction, error) {
	result := datastore.DB.Create(&transaction)
	if result.Error != nil {
		return nil, result.Error
	}

	return transaction, nil
}

func (r *repository) ReadTransactions(filters Filters) (transactions *[]presenter.Transaction, err error) {
	query := datastore.DB.Order("id desc")
	if len(filters.Columns) > 0 {
		query = query.Select(filters.Columns)
	}
	if len(filters.Merchants) > 0 {
		query = query.Where("merchant_id in ?", filters.Merchants)
	}
	if filters.Days > 0 {
		duration := time.Duration(filters.Days) * 24 * time.Hour
		query = query.Where("transactions.created_at > ?", time.Now().Add(-duration))
	}

	// TODO: Fix this into the filters struct
	err = query.Joins("Payment").Find(&transactions).Error

	return
}

func (r *repository) ReadTransaction(id uint) (transaction *entities.Transaction, err error) {
	err = datastore.DB.First(&transaction, id).Error
	return
}

func (r *repository) ReadTransactionsByMerchant(merchantId uint) (transaction *[]presenter.Transaction, err error) {
	err = datastore.DB.Where("merchant_id", merchantId).Find(&transaction).Error
	return
}

func (r *repository) UpdateTransaction(transaction *entities.Transaction) (*entities.Transaction, error) {
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
