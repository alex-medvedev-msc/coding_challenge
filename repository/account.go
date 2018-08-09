package repository

import (
	"database/sql"
	"github.com/messwith/coding_challenge/models"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

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


