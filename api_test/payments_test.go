package api_test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/messwith/coding_challenge/api"
	"github.com/shopspring/decimal"
	"github.com/messwith/coding_challenge/models"
	"github.com/stretchr/testify/assert"
	"strconv"
	"math/rand"
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

	// simple happy case
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

	// check if we can completely drain account balance
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

	// payment with zero amount is not allowed
	pr := api.PaymentRequest{
		Amount: decimal.NewFromFloat(0),
		FromAccount: "1",
		ToAccount: "2",
	}
	require.NotNil(t, request("POST", "/payments", pr, nil))

	// trying to send more than we can
	pr.Amount = decimal.NewFromFloat(1)
	require.NotNil(t, request("POST", "/payments", pr, nil))

	// trying to send negative amount
	pr.Amount = decimal.NewFromFloat(-1)
	require.NotNil(t, request("POST", "/payments", pr, nil))

	// checking system precision
	pr.Amount, _ = decimal.NewFromString("0.50000000000001")
	require.NotNil(t, request("POST", "/payments", pr, nil))

	// balance of accounts must be intact after test
	var accounts []models.Account
	require.Nil(t, request("GET", "/accounts", nil, &accounts))
	assert.True(t, accounts[0].Balance.Equal(decimal.NewFromFloat(0.5)))
	assert.True(t, accounts[1].Balance.Equal(decimal.NewFromFloat(1.5)))
}

func BenchmarkCreatePaymentParallel(b *testing.B) {
	defer clearDB()
	require.Nil(b, createAccount("1", decimal.NewFromFloat(10000)))
	require.Nil(b, createAccount("2", decimal.NewFromFloat(10000)))
	require.Nil(b, createAccount("3", decimal.NewFromFloat(10000)))

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			from := rand.Intn(3)+1
			to := from -1
			if to < 1 {
				to = 3
			}
			pr := api.PaymentRequest{
				Amount: decimal.NewFromFloat(0.01),
				FromAccount: strconv.Itoa(from),
				ToAccount: strconv.Itoa(to),
			}
			assert.Nil(b, request("POST", "/payments", pr, nil))
		}
	})
}
