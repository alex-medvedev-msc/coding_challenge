package api_test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/messwith/coding_challenge/models"
)

func TestGetAccounts(t *testing.T) {
	defer clearDB()
	require.Nil(t, createRandomAccounts(100))
	var accounts []models.Account
	require.Nil(t, request("GET", "/accounts", nil, &accounts))
	require.Len(t, accounts, 100)
	require.Equal(t, accounts[0].ID, "1")
}

func TestGetAccountsEmpty(t *testing.T) {
	defer clearDB()
	var accounts []models.Account
	require.Nil(t, request("GET", "/accounts", nil, &accounts))
	require.NotNil(t, accounts)
	require.Len(t, accounts, 0)
}
