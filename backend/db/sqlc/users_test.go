package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/blessedmadukoma/trackit-chima/util"
	"github.com/stretchr/testify/require"
)

// createRandomUser creates a random user account in the database
func createRandomUser(t *testing.T) User {
	arg := CreateUserAccountParams{
		Firstname: util.RandomUser(),
		Lastname:  util.RandomUser(),
		Email:     util.RandomEmail(),
		Mobile:    util.RandomMobileNumber(),
		Password:  util.RandomPassword(),
	}

	userAccount, err := testQueries.CreateUserAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, userAccount)

	require.Equal(t, arg.Firstname, userAccount.Firstname)
	require.Equal(t, arg.Lastname, userAccount.Lastname)
	require.Equal(t, arg.Email, userAccount.Email)
	require.Equal(t, arg.Mobile, userAccount.Mobile)
	require.Equal(t, arg.Password, userAccount.Password)

	require.NotZero(t, userAccount.ID)
	require.NotZero(t, userAccount.CreatedAt)

	return userAccount
}

// TestCreateUserAccount tests the CreateUserAccount method which creates a new user's account
func TestCreateUserAccount(t *testing.T) {
	createRandomUser(t)
}

// TestGetUserAccountByID tests the GetUserAccountByID method which gets a user's account information through the ID
func TestGetUserAccountByID(t *testing.T) {
 userAccount1 := createRandomUser(t)

	userAccount2, err := testQueries.GetUserAccountByID(context.Background(), userAccount1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, userAccount2)

	require.Equal(t, userAccount1.ID, userAccount2.ID)
	require.Equal(t, userAccount1.Firstname, userAccount2.Firstname)
	require.Equal(t, userAccount1.Lastname, userAccount2.Lastname)
	require.Equal(t, userAccount1.Email, userAccount2.Email)
	require.Equal(t, util.PasswordMatch(userAccount1.Password, userAccount2.Password), util.PasswordMatch(userAccount2.Password, userAccount1.Password))
	require.WithinDuration(t, userAccount1.CreatedAt, userAccount2.CreatedAt, time.Second)
}

// TestListUsers tests the ListUsers method which lists all the users' accounts
func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	arg := ListUsersParams{
		Limit: 5,
		Offset: 5,
	}

	userAccounts, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, userAccounts, 5)

	for _, userAccount := range userAccounts {
		require.NotEmpty(t, userAccount)
	}
}

// TestUpdateUserAccountInfo tests the UpdateUserAccountInfo method which updates the details of the user based on the ID
func TestUpdateUserAccountInfo(t *testing.T) {
	userAccount1 := createRandomUser(t)

	arg := UpdateUserAccountParams{
		ID: userAccount1.ID,
		Firstname: userAccount1.Firstname,
		Lastname: userAccount1.Lastname,
		Email: userAccount1.Email,
		Mobile: userAccount1.Mobile,
		Password: userAccount1.Password,
	}

	userAccount2, err := testQueries.UpdateUserAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, userAccount2)

	require.Equal(t, userAccount1.ID, userAccount2.ID)
	require.Equal(t, userAccount1.Firstname, userAccount2.Firstname)
	require.Equal(t, userAccount1.Lastname, userAccount2.Lastname)
	require.Equal(t, userAccount1.Email, userAccount2.Email)
	require.Equal(t, userAccount1.Mobile, userAccount2.Mobile)
	require.Equal(t, userAccount1.Password, userAccount2.Password)

	require.WithinDuration(t, userAccount1.CreatedAt, userAccount2.CreatedAt, time.Second)
}

// TestDeleteUserAccount tests the DeleteUserAccount method which deletes an existing record
func TestDeleteUserAccount(t *testing.T) {
	userAccount1 := createRandomUser(t)

	err := testQueries.DeleteUserAccount(context.Background(), userAccount1.ID)
	require.NoError(t, err)

	userAccount2, err := testQueries.GetUserAccountByID(context.Background(), userAccount1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, userAccount2)
}