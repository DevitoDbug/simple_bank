package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"simple_bank/util"
	"testing"
)

func createRandomAccountForEntryTest(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
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

func createRandomEntry(t *testing.T) Entry {
	account := createRandomAccountForEntryTest(t)
	args := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}
	createdEntry, err := testQueries.CreateEntry(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, createdEntry)

	require.Equal(t, args.AccountID, createdEntry.AccountID)
	require.Equal(t, args.Amount, createdEntry.Amount)
	return createdEntry
}

func createRandom10Entries(t *testing.T) ([]Entry, Account) {
	account := createRandomAccountForEntryTest(t)
	var createdEntries []Entry
	args := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}
	for i := 0; i < 10; i++ {
		entry, _ := testQueries.CreateEntry(context.Background(), args)
		createdEntries = append(createdEntries, entry)
	}
	return createdEntries, account
}

func TestGetEntry(t *testing.T) {
	entry := createRandomEntry(t)

	fetchedEntry, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedEntry)

	require.Equal(t, entry.ID, fetchedEntry.ID)
	require.Equal(t, entry.AccountID, fetchedEntry.AccountID)
	require.Equal(t, entry.Amount, fetchedEntry.Amount)
	require.Equal(t, entry.CreatedAt, fetchedEntry.CreatedAt)
}

func TestListEntries(t *testing.T) {
	_, account := createRandom10Entries(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(t)
	}

	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	fetchedEntries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedEntries)

	require.Len(t, fetchedEntries, 5)

	for _, entry := range fetchedEntries {
		require.NotEmpty(t, entry)
	}
}
