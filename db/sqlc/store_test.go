package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx1(t *testing.T) {
	store := NewStore(testDB)

	//creating accounts
	//initial account balances is 1000
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	noOfTransfers := 5
	amount := int64(10)

	transferResults := make(chan TransferTxResult)
	transferResultsError := make(chan error)

	for i := 0; i < noOfTransfers; i++ {
		go func() {
			//we cannot check the tests here because it is in a separate go routine form the one our test is running on
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: account1.ID,
				ToAccountId:   account2.ID,
				Amount:        amount,
			})

			transferResultsError <- err
			transferResults <- result
		}()
	}

	//Check for no error in any of the transactions
	for i := 0; i < noOfTransfers; i++ {
		err := <-transferResultsError
		require.NoError(t, err)

		result := <-transferResults
		require.NotEmpty(t, result)

		//Checking the result
		//Checking for transfer in result
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//Checking the entries
		//Checking the fromEntry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)

		//Checking the toEntry
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
	}
}
