package service

import (
"github.com/messwith/coding_challenge/models"
"github.com/messwith/coding_challenge/repository"
"github.com/messwith/coding_challenge/errors"
)

// AccountService is an interface which allows to abstract from details of storing accounts in some storage
type AccountService interface {
	GetAccounts() ([]models.Account, error)
}

// SqlAccountService uses sql db as storage, currently it is postgres
// but you can change it with simple import replacement, e.g. _ ".../lib/pq" -> _ ".../lib/mysql"
type SqlAccountService struct {
	accountRep *repository.AccountRepository
}

// NewSqlAccountService creates ready to use SqlAccountService instance
func NewSqlAccountService(
	accountRep *repository.AccountRepository) *SqlAccountService {

	return &SqlAccountService{
		accountRep: accountRep,
	}
}

// GetAccounts returns all accounts from db without filtering and pagination
func (ps *SqlAccountService) GetAccounts() ([]models.Account, error) {

	accounts, err := ps.accountRep.GetAccounts()
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	return accounts, nil
}
