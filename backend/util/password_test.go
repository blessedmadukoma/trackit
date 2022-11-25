package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashedPassword(t *testing.T) {
	password := RandomString(6)

	hashed, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashed)

	err = PasswordMatch(password, hashed)
	require.NoError(t, err)

	wrongPassword := RandomString(8)
	err = PasswordMatch(wrongPassword, hashed)
	require.Error(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}