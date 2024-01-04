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

	transferResults := make(chan TransferTxResult)
	transferResultsError := make(chan error)

	noOfTransfers := 5
	amount := int64(10)
	for i := 0; i < noOfTransfers; i++ {
		go func() {
			//we cannot check the tests here because it is in a separate go routine form the one our test is running on
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: account1.ID,
				ToAccountId:   account2.ID,
				Amount:        amount,
			})
			transferResults <- result
			transferResultsError <- err
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
	}

}
