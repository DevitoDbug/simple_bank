package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"simple_bank/util"
	"testing"
)

func createRandomAccountForTransaction(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  1000,
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID) //Check if the id is set
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	//creating accounts
	//initial account balances is 1000
	account1 := createRandomAccountForTransaction(t)
	account2 := createRandomAccountForTransaction(t)

}
