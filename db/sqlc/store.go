package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all the functions to execute queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function thad does db transaction
func (s *Store) execTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("query execution error: %v\n rollback error: %v\n", err, rbErr)
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

type TransferTxParams struct {
	FromAccountId int64 `json:"fromAccountId"`
	ToAccountId   int64 `json:"toAccountId"`
	Amount        int64 `json:"amount"`
}
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"fromAccount"`
	ToAccount   Account  `json:"toAccount"`
	FromEntry   Entry    `json:"fromEntry"`
	ToEntry     Entry    `json:"toEntry"`
}

// TransferTx transfers money from one account to another account
// Creates a transfer record
// Add account entries
// Update account balance
func (s *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	createTransferParams := CreateTransferParams{
		FromAccountID: arg.FromAccountId,
		ToAccountID:   arg.ToAccountId,
		Amount:        arg.Amount,
	}

	createSenderEntryParams := CreateEntryParams{
		AccountID: arg.FromAccountId,
		Amount:    -arg.Amount,
	}

	createRecipientEntryParams := CreateEntryParams{
		AccountID: arg.ToAccountId,
		Amount:    arg.Amount,
	}

	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		// Creates a transfer record
		result.Transfer, err = q.CreateTransfer(ctx, createTransferParams)
		if err != nil {
			return err
		}

		//Add to and from accounts
		// Add account entries
		result.FromEntry, err = q.CreateEntry(ctx, createSenderEntryParams)
		if err != nil {
			return err
		}
		result.ToEntry, err = q.CreateEntry(ctx, createRecipientEntryParams)
		if err != nil {
			return err
		}

		if arg.FromAccountId > arg.ToAccountId {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountId, -arg.Amount, arg.ToAccountId, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountId, arg.Amount, arg.FromAccountId, -arg.Amount)
		}
		return nil
	})
	if err != nil {
		return result, err
	}
	return result, nil
}

func addMoney(
	ctx context.Context,
	q *Queries,
	toAccountId int64,
	amountTo int64,
	fromAccountId int64,
	amountFrom int64,
) (toAccount Account, fromAccount Account, err error) {
	toAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     toAccountId,
		Amount: amountTo,
	})
	if err != nil {
		return
	}

	fromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     fromAccountId,
		Amount: amountFrom,
	})
	if err != nil {
		return
	}

	return
}
