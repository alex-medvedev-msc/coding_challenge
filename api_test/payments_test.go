package api_test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/messwith/coding_challenge/api"
	"github.com/shopspring/decimal"
	"github.com/messwith/coding_challenge/models"
	"github.com/stretchr/testify/assert"
)

func TestGetPayments(t *testing.T) {
	defer clearDB()
	require.Nil(t, createRandomAccounts(2))
	pr := api.PaymentRequest{
		Amount: decimal.NewFromFloat(0.1),
		FromAccount: "1",
		ToAccount: "2",
	}
	require.Nil(t, request("POST", "/payments", pr, nil))
	pr2 := api.PaymentRequest{
		Amount: decimal.NewFromFloat(0.1),
		FromAccount: "2",
		ToAccount: "1",
	}
	require.Nil(t, request("POST", "/payments", pr2, nil))

	var payments []models.Payment
	require.Nil(t, request("GET", "/payments", nil, &payments))

	require.Len(t, payments, 4)
	assert.Equal(t, payments[0].Direction, models.DirectionOut)
	assert.Equal(t, payments[1].Direction, models.DirectionIn)
	assert.Equal(t, payments[2].Direction, models.DirectionOut)
	assert.Equal(t, payments[3].Direction, models.DirectionIn)

	assert.Equal(t, payments[0].Account, pr.FromAccount)
	assert.Equal(t, payments[1].Account, pr.ToAccount)
	assert.Equal(t, payments[2].Account, pr2.FromAccount)
	assert.Equal(t, payments[3].Account, pr2.ToAccount)
}

func TestCreatePayment(t *testing.T) {

}



