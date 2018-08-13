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
	defer clearDB()
	require.Nil(t, createAccount("1", decimal.NewFromFloat(0.5)))
	require.Nil(t, createAccount("2", decimal.NewFromFloat(1.5)))
	pr := api.PaymentRequest{
		Amount: decimal.NewFromFloat(0.1),
		FromAccount: "1",
		ToAccount: "2",
	}
	require.Nil(t, request("POST", "/payments", pr, nil))

	var accounts []models.Account
	require.Nil(t, request("GET", "/accounts", nil, &accounts))
	assert.True(t, accounts[0].Balance.Equal(decimal.NewFromFloat(0.4)))
	assert.True(t, accounts[1].Balance.Equal(decimal.NewFromFloat(1.6)))

	var payments []models.Payment
	require.Nil(t, request("GET", "/payments", nil, &payments))

	require.Len(t, payments, 2)
	assert.Equal(t, payments[0].Direction, models.DirectionOut)
	assert.Equal(t, payments[1].Direction, models.DirectionIn)

	assert.Equal(t, payments[0].Account, pr.FromAccount)
	assert.Equal(t, payments[1].Account, pr.ToAccount)

	pr.Amount = decimal.NewFromFloat(0.4)
	require.Nil(t, request("POST", "/payments", pr, nil))

	require.Nil(t, request("GET", "/accounts", nil, &accounts))
	assert.True(t, accounts[0].Balance.Equal(decimal.Zero))
	assert.True(t, accounts[1].Balance.Equal(decimal.NewFromFloat(2.0)))
}

func TestCreatePaymentWrong(t *testing.T) {
	defer clearDB()
	require.Nil(t, createAccount("1", decimal.NewFromFloat(0.5)))
	require.Nil(t, createAccount("2", decimal.NewFromFloat(1.5)))
	pr := api.PaymentRequest{
		Amount: decimal.NewFromFloat(0),
		FromAccount: "1",
		ToAccount: "2",
	}
	require.NotNil(t, request("POST", "/payments", pr, nil))

	pr.Amount = decimal.NewFromFloat(1)
	require.NotNil(t, request("POST", "/payments", pr, nil))

	pr.Amount = decimal.NewFromFloat(-1)
	require.NotNil(t, request("POST", "/payments", pr, nil))

	pr.Amount, _ = decimal.NewFromString("0.50000000000001")
	require.NotNil(t, request("POST", "/payments", pr, nil))


	var accounts []models.Account
	require.Nil(t, request("GET", "/accounts", nil, &accounts))
	assert.True(t, accounts[0].Balance.Equal(decimal.NewFromFloat(0.5)))
	assert.True(t, accounts[1].Balance.Equal(decimal.NewFromFloat(1.5)))

}



