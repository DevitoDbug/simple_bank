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
	FromAccountId Account `json:"fromAccountId"`
	ToAccountId   Account `json:"toAccountId"`
	Amount        int64   `json:"amount"`
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

}
