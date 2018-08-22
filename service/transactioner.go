package service

import (
	"github.com/messwith/coding_challenge/models"
	"github.com/messwith/coding_challenge/repository"
	"github.com/shopspring/decimal"
	"database/sql"
	"github.com/messwith/coding_challenge/errors"
)

// Transactioner is an interface which allows to abstract from details of storing payments in some storage
type Transactioner interface {
	Do(incoming *models.Payment, outgoing *models.Payment) error
}

// SqlTransactioner uses sql db as storage, currently it is postgres
// but you can change it with simple import replacement, e.g. _ ".../lib/pq" -> _ ".../lib/mysql"
type SqlTransactioner struct {
	accountRep *repository.AccountRepository
	paymentRep *repository.PaymentRepository
}

// NewSqlTransactioner creates ready to use SqlTransactioner instance
func NewSqlTransactioner(
		accountRep *repository.AccountRepository,
		paymentRep *repository.PaymentRepository) *SqlTransactioner {

		return &SqlTransactioner{
			accountRep: accountRep,
			paymentRep: paymentRep,
		}
}

// Do performs transfer of funds between two accounts, requires already validated matching payments
func (t *SqlTransactioner) Do(incoming *models.Payment, outgoing *models.Payment) error {

	tx, err := t.paymentRep.BeginTx()
	if err != nil {
		return errors.NewInternalError(err)
	}
	defer tx.Rollback()

	// locking sender balance for update
	sender, err := t.accountRep.LockAccount(tx, incoming.FromAccount)
	if err == sql.ErrNoRows{
		return errors.NewAccountNotFoundError(incoming.FromAccount)
	} else if err != nil {
		return errors.NewInternalError(err)
	}

	// checking if sender has enough funds to perform payment
	sender.Balance = sender.Balance.Sub(incoming.Amount)
	if sender.Balance.Cmp(decimal.Zero) == -1 {
		return errors.NewDataError("sender has unsufficient funds for this payment")
	}

	// locking receiver for update
	receiver, err := t.accountRep.LockAccount(tx, outgoing.ToAccount)
	if err == sql.ErrNoRows{
		return errors.NewAccountNotFoundError(incoming.FromAccount)
	} else if err != nil {
		return errors.NewInternalError(err)
	}

	// we do not support other currencies
	if sender.Currency != receiver.Currency {
		return errors.NewDataError("sender and receiver has different currencies, conversion is not supported yet")
	}
	receiver.Balance = receiver.Balance.Add(incoming.Amount)

	// creating payments in db
	if err := t.createPayments(tx, incoming, outgoing); err != nil {
		return err
	}

	// finally, updating account balances
	if err := t.updateAccountBalances(tx, &sender, &receiver); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return errors.NewInternalError(err)
	}
	return nil
}

func (t *SqlTransactioner) createPayments(tx *sql.Tx, incoming, outgoing *models.Payment) error {
	if err := t.paymentRep.CreatePayment(tx, outgoing); err != nil {
		return errors.NewInternalError(err)
	}

	if err := t.paymentRep.CreatePayment(tx, incoming); err != nil {
		return errors.NewInternalError(err)
	}
	return nil
}

func (t *SqlTransactioner) updateAccountBalances(tx *sql.Tx, sender *models.Account, receiver *models.Account) error {
	if err := t.accountRep.UpdateAccountBalance(tx, sender.ID, sender.Balance); err != nil {
		return errors.NewInternalError(err)
	}

	if err := t.accountRep.UpdateAccountBalance(tx, receiver.ID, receiver.Balance); err != nil {
		return errors.NewInternalError(err)
	}
	return nil
}
