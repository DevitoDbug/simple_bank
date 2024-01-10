package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx(t *testing.T) {
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
			ctx := context.Background()
			//we cannot check the tests here because it is in a separate go routine form the one our test is running on
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountId: account1.ID,
				ToAccountId:   account2.ID,
				Amount:        amount,
			})

			transferResultsError <- err
			transferResults <- result
		}()
	}

	//Check for no error in any of the transactions
	existed := make(map[int]bool)
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

		//Checking account
		//Checking from account
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		//Checking to account
		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		//Checking amount difference in account1
		diff1 := account1.Balance - fromAccount.Balance
		//Checking amount difference in account2
		diff2 := toAccount.Balance - account2.Balance

		require.Equal(t, diff2, diff1)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= noOfTransfers)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	//Check final updated balance
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NotEmpty(t, updatedAccount1)
	require.NoError(t, err)
	require.Equal(t, account1.Balance-amount*int64(noOfTransfers), updatedAccount1.Balance)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NotEmpty(t, updatedAccount2)
	require.NoError(t, err)
	require.Equal(t, account2.Balance+amount*int64(noOfTransfers), updatedAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	//creating accounts
	//initial account balances is 1000
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	noOfTransfers := 20
	amount := int64(10)

	transferResultsError := make(chan error)

	for i := 0; i < noOfTransfers; i++ {
		fromAccountId := account1.ID
		toAccountId := account2.ID

		if i%2 == 0 {
			fromAccountId = account2.ID
			toAccountId = account1.ID
		}

		go func() {
			ctx := context.Background()
			//we cannot check the tests here because it is in a separate go routine form the one our test is running on
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountId: fromAccountId,
				ToAccountId:   toAccountId,
				Amount:        amount,
			})
			transferResultsError <- err
		}()
	}

	for i := 0; i < noOfTransfers; i++ {
		err := <-transferResultsError
		require.NoError(t, err)
	}

	//Check final updated balance
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NotEmpty(t, updatedAccount1)
	require.NoError(t, err)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NotEmpty(t, updatedAccount2)
	require.NoError(t, err)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
