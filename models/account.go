package models

import "github.com/shopspring/decimal"

// Account is an entity which can transfer funds to other account
type Account struct {
	ID       string          `json:"id"`
	Owner    string          `json:"owner"`
	Balance  decimal.Decimal `json:"balance"`
	Currency string          `json:"currency"`
}
