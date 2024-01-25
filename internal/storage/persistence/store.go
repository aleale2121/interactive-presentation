package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// Store defines all functions to execute db queries and transactions
type Store interface {
	Querier
	GetCurrentPoll(context.Context, uuid.UUID) (CurrentPoll, error)
	UpdateCurrentPollToForwardTx(ctx context.Context, id uuid.UUID) (CurrentPoll, error)
	UpdateCurrentPollToPreviousTx(ctx context.Context, id uuid.UUID) (CurrentPoll, error) 
	VoteCurrentPollTx(context.Context, VoteParams) error

}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	db *sql.DB
	mu sync.Mutex
	*Queries
}

// NewStore creates a new store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// ExecTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
