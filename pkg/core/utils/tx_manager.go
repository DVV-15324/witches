package utils

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type txContextKey struct{}

var txKey = txContextKey{}

type TxManager struct {
	db *gorm.DB
}

func NewTxManager(db *gorm.DB) *TxManager {
	return &TxManager{db: db}
}

func (t *TxManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx := t.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	txCtx := context.WithValue(ctx, txKey, tx)
	err := fn(txCtx)
	if err != nil {
		if rbErr := tx.Rollback().Error; rbErr != nil {
			return &TransactionError{
				OriginalError: err,
				RollbackError: rbErr,
			}
		}
		return err
	}
	if commitErr := tx.Commit().Error; commitErr != nil {
		return commitErr
	}
	return nil
}

func GetTxFromContext(ctx context.Context) (*gorm.DB, error) {
	tx, ok := ctx.Value(txKey).(*gorm.DB)
	if !ok {
		return nil, errors.New("no transaction found in context")
	}
	return tx, nil
}

type TransactionError struct {
	OriginalError error
	RollbackError error
}

func (e *TransactionError) Error() string {
	return "transaction failed: " + e.OriginalError.Error() + " (rollback error: " + e.RollbackError.Error() + ")"
}
