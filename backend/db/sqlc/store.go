package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries
type Store struct {
	*Queries // composition (embedding) -> we are extending functionality i.e. transactions, a better way than inheritance
	db       *sql.DB
}

// NewStore creates a new store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// ExecTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := New(tx)
	err = fn(query)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction (tx) err: %v, rollback (rb) err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
