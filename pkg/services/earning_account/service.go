package earning_account

import (
	"merchants.sidooh/api/presenter"
	"merchants.sidooh/pkg"
	"merchants.sidooh/pkg/entities"
	"merchants.sidooh/pkg/services/earning_account_transaction"
	"merchants.sidooh/utils"
)

type Service interface {
	FetchAccountsByAccountId(accountId uint) (*[]presenter.EarningAccount, error)
	FetchAccountsByMerchant(merchantId uint) (*[]presenter.EarningAccount, error)
	CreateAccount(data *entities.EarningAccount) (*entities.EarningAccount, error)

	CreditAccount(accountId uint, amount float32) (*entities.EarningAccount, error)
	DebitAccount(accountId uint, amount float32) (*entities.EarningAccount, error)
}

type service struct {
	repository             Repository
	earningAccTxRepository earning_account_transaction.Repository
}

func (s *service) FetchAccountsByAccountId(accountId uint) (*[]presenter.EarningAccount, error) {
	return s.repository.ReadAccountsByAccountId(accountId)
}

func (s *service) FetchAccountsByMerchant(merchantId uint) (results *[]presenter.EarningAccount, err error) {
	accounts, err := s.repository.ReadAccountsByMerchant(merchantId)
	if err != nil {
		return nil, err
	}
	utils.ConvertStruct(accounts, results)

	return
}

func (s *service) CreditAccount(accountId uint, amount float32) (*entities.EarningAccount, error) {
	// TODO: use db tx
	account, err := s.repository.ReadAccount(accountId)
	if err != nil {
		return nil, err
	}

	// create tx
	s.earningAccTxRepository.CreateTransaction(&entities.EarningAccountTransaction{
		Type:             "CREDIT",
		Amount:           amount,
		EarningAccountId: accountId,
	})

	// update acc
	account.Amount += amount
	account, err = s.repository.UpdateAccount(account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *service) DebitAccount(accountId uint, amount float32) (*entities.EarningAccount, error) {
	// TODO: use db tx
	account, err := s.repository.ReadAccount(accountId)
	if err != nil {
		return nil, err
	}

	if account.Amount < amount {
		return nil, pkg.ErrInsufficientBalance
	}

	// create tx
	s.earningAccTxRepository.CreateTransaction(&entities.EarningAccountTransaction{
		Type:             "DEBIT",
		Amount:           amount,
		EarningAccountId: accountId,
	})

	// update acc
	account.Amount -= amount
	account, err = s.repository.UpdateColumn(account, "amount", account.Amount)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *service) CreateAccount(data *entities.EarningAccount) (*entities.EarningAccount, error) {
	return s.repository.CreateAccount(data)
}

func NewService(r Repository, earningAccTxRepo earning_account_transaction.Repository) Service {
	return &service{repository: r, earningAccTxRepository: earningAccTxRepo}
}
