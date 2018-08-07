package models

import (
	"github.com/shopspring/decimal"
	"errors"
)

type Payment struct {
	Account string `json:"account"`
	ToAccount string `json:"to_account"`
	FromAccount string `json:"from_account"`
	Direction string `json:"direction"`
	Amount decimal.Decimal `json:"amount"`
}

const (
	DirectionIn = "ingoing"
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
