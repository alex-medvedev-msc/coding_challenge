package repository

import (
	"database/sql"
	"github.com/messwith/coding_challenge/models"
)

// AccountRepository is essentially set of methods for working with accounts in db
type AccountRepository struct {
	db *sql.DB
}

// NewAccountRepository creates ready to use account repository
func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

// GetAccounts loads all accounts from db without pagination
func (ar *AccountRepository) GetAccounts() ([]models.Account, error) {

	accounts := []models.Account{}

	rows, err := ar.db.Query(`SELECT id, owner, balance, currency FROM accounts`)
	if err == sql.ErrNoRows {
		return accounts, nil
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		account := models.Account{}
		err = rows.Scan(&account.Id, &account.Owner, &account.Balance, &account.Currency)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}


