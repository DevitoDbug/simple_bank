package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"simple_bank/util"
	"testing"
)

func createRandomUserForTransfer(t *testing.T, currentBalance int64) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  currentBalance,
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

func creatRandomTransfer(t *testing.T) Transfer {
	user1 := createRandomUserForTransfer(t, 1000)
	user2 := createRandomUserForTransfer(t, 1000)

	arg := CreateTransferParams{
		FromAccountID: user1.ID,
		ToAccountID:   user2.ID,
		Amount:        500,
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	return transfer
}

func creatRandom10TransfersForTwoUsers(t *testing.T) (transferList []Transfer, account1 Account, account2 Account) {
	var transfers []Transfer
	user1 := createRandomUserForTransfer(t, 1000)
	user2 := createRandomUserForTransfer(t, 1000)

	arg := CreateTransferParams{
		FromAccountID: user1.ID,
		ToAccountID:   user2.ID,
		Amount:        500,
	}

	for i := 0; i < 10; i++ {
		transfer, err := testQueries.CreateTransfer(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, transfer)

		require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
		require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
		require.Equal(t, arg.Amount, transfer.Amount)
		transfers = append(transfers, transfer)
	}
	return transfers, user1, user2
}

func TestGetTransfer(t *testing.T) {
	transfer := creatRandomTransfer(t)

	fetchedTransfer, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedTransfer)

	require.Equal(t, transfer.CreatedAt, fetchedTransfer.CreatedAt)
	require.Equal(t, transfer.ID, fetchedTransfer.ID)
	require.Equal(t, transfer.Amount, fetchedTransfer.Amount)
	require.Equal(t, transfer.FromAccountID, fetchedTransfer.FromAccountID)
	require.Equal(t, transfer.ToAccountID, fetchedTransfer.ToAccountID)
}

func TestListTransfers(t *testing.T) {
	_, user1, user2 := creatRandom10TransfersForTwoUsers(t)
	arg := ListTransfersParams{
		FromAccountID: user1.ID,
		ToAccountID:   user2.ID,
		Limit:         5,
		Offset:        5,
	}

	fetchedTransfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedTransfers)

	require.Len(t, fetchedTransfers, 5)
	for _, transfer := range fetchedTransfers {
		require.NotEmpty(t, transfer)
	}
}
