package repository

import (
	"database/sql"
	"github.com/messwith/coding_challenge/models"
)

// PaymentRepository is a collection of methods for storing and querying payments in db
type PaymentRepository struct {
	db *sql.DB
}

// NewPaymentRepository creates ready to use payment repository objects
func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// CreatePayment creates new payment in db
func (pr *PaymentRepository) CreatePayment(tx *sql.Tx, payment *models.Payment) error {

	_, err := tx.Exec(`INSERT INTO payments (account, to_account, from_account, direction, amount) 
												VALUES ($1, $2, $3, $4, $5)`,
		payment.Account, payment.ToAccount, payment.FromAccount, payment.Direction, payment.Amount)
	return err
}

// GetPayments returns all payments from db without filtering and pagination
func (pr *PaymentRepository) GetPayments() ([]models.Payment, error) {

	payments := []models.Payment{}

	rows, err := pr.db.Query(`SELECT account, to_account, from_account, direction, amount from payments`)
	if err == sql.ErrNoRows {
		return payments, err
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		payment := models.Payment{}
		err = rows.Scan(&payment.Account, &payment.ToAccount, &payment.FromAccount, &payment.Direction, &payment.Amount)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}
