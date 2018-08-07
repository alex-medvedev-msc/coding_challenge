package models

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/shopspring/decimal"
)

func TestPayment_Validate(t *testing.T) {
	p := Payment{}
	assert.NotNil(t, p.Validate())

	// zero amount
	p = Payment{FromAccount: "1", Direction: DirectionIn, Account: "2"}
	assert.NotNil(t, p.Validate())

	// negative amount
	p = Payment{FromAccount: "1", Direction: DirectionIn, Account: "2", Amount: decimal.NewFromFloat(-1.2)}
	assert.NotNil(t, p.Validate())

	// empty account
	p = Payment{FromAccount: "1", Direction: DirectionIn, Amount: decimal.NewFromFloat(1.2)}
	assert.NotNil(t, p.Validate())

	// wrong direction
	p = Payment{FromAccount: "1", Direction: DirectionOut, Account: "2", Amount: decimal.NewFromFloat(1.2)}
	assert.NotNil(t, p.Validate())

	// valid incoming
	p = Payment{FromAccount: "1", Direction: DirectionIn, Account: "2", Amount: decimal.NewFromFloat(0.00001)}
	assert.Nil(t, p.Validate())

	// valid outgoing
	p = Payment{Account: "1", Direction: DirectionOut, ToAccount: "2", Amount: decimal.NewFromFloat(13423423423.2)}
	assert.Nil(t, p.Validate())
}

func TestNewIncomingPayment(t *testing.T) {
	incoming := NewIncomingPayment("from", "to", decimal.NewFromFloat(1.432))

	assert.Equal(t, incoming.Direction, DirectionIn)
	assert.Equal(t, incoming.Account, "to")
}

func TestNewOutgoingPayment(t *testing.T) {
	outgoing := NewOutgoingPayment("from", "to", decimal.NewFromFloat(1.432))

	assert.Equal(t, outgoing.Direction, DirectionOut)
	assert.Equal(t, outgoing.Account, "from")
}
