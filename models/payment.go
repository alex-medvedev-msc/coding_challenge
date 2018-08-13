package models

import (
	"github.com/shopspring/decimal"
	"errors"
)

// Payment holds information about moving funds from one account to other
// There is always two matching payments in opposite direction for each transaction
type Payment struct {
	Account string `json:"account"`
	ToAccount string `json:"to_account"`
	FromAccount string `json:"from_account"`
	Direction string `json:"direction"`
	Amount decimal.Decimal `json:"amount"`
}

const (
	// DirectionIn means that payment is incoming for payment owner account
	DirectionIn = "ingoing"
	// DirectionOut means that payment is outgoing for payment owner account
	DirectionOut = "outgoing"
)

// Validate validates payment
// checking consistence of accounts fields and direction
// Validate DOES NOT perform checking of account ability to make this payment
func (p *Payment) Validate() error {

	if !(p.Direction == DirectionIn || p.Direction == DirectionOut) {
		return errors.New("unknown direction")
	}
	if p.Amount.Cmp(decimal.Zero) != 1 {
		return errors.New("amount must be positive")
	}
	if p.Account == "" {
		return errors.New("payment owner must be present")
	}
	if p.Direction == DirectionOut && p.ToAccount == "" {
		return errors.New("outgoing payment must have a destination")
	}
	if p.Direction == DirectionIn && p.FromAccount == "" {
		return errors.New("ingoing payment must have a source")
	}

	return nil
}

// NewIncomingPayment creates new payment with direction = DirectionIn and owner = to
func NewIncomingPayment(from, to string, amount decimal.Decimal) *Payment {
	return &Payment{
		Account: to,
		FromAccount: from,
		Direction: DirectionIn,
		Amount: amount,
	}
}

// NewOutgoingPayment creates new payment with direction = DirectionOut and owner = from
func NewOutgoingPayment(from, to string, amount decimal.Decimal) *Payment {
	return &Payment{
		Account: from,
		ToAccount: to,
		Direction: DirectionOut,
		Amount: amount,
	}
}