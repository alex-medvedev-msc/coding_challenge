package models

import "github.com/shopspring/decimal"

type Account struct {
	Id string `json:"id"`
	Owner string `json:"owner"`
	Balance decimal.Decimal `json:"balance"`
	Currency string `json:"currency"`
}
