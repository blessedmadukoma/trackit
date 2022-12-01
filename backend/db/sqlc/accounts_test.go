package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// createRandomAccount creates a random account in the database
func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	arg := CreateAccountParams{
		UserID:  user.ID,
		Balance: 0,
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.UserID, account.UserID)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

// TestCreateAccount tests the CreateUserAccount method which creates a new user's account
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

// TestGetAccountByID gets the account information through the ID
func TestGetAccountByID(t *testing.T) {
	account1 := createRandomAccount(t)

	account2, err := testQueries.GetAccountByID(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.UserID, account2.UserID)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

// TestListAccounts lists all the accounts
func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

// TestUpdateAccount tests the UpdateAccounts method which updates the details of the user based on the ID
func TestUpdateAccounts(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: 10,
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.UserID, account2.UserID)

	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}
