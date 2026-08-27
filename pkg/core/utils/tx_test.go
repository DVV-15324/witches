package utils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	return db
}

func TestNewTxManager(t *testing.T) {
	db := setupTestDB(t)
	txManager := NewTxManager(db)
	assert.NotNil(t, txManager)
}

func TestTxManager_WithinTransaction_Success(t *testing.T) {
	db := setupTestDB(t)
	txManager := NewTxManager(db)

	ctx := context.Background()
	err := txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		tx, err := GetTxFromContext(txCtx)
		assert.NoError(t, err)
		assert.NotNil(t, tx)
		return nil
	})

	assert.NoError(t, err)
}

func TestTxManager_WithinTransaction_Rollback(t *testing.T) {
	db := setupTestDB(t)
	txManager := NewTxManager(db)

	ctx := context.Background()
	err := txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		return assert.AnError
	})

	assert.Error(t, err)
}

func TestGetTxFromContext_NoTransaction(t *testing.T) {
	ctx := context.Background()
	tx, err := GetTxFromContext(ctx)
	assert.Error(t, err)
	assert.Equal(t, "no transaction found in context", err.Error())
	assert.Nil(t, tx)
}

func TestTransactionError(t *testing.T) {
	originalErr := assert.AnError
	rollbackErr := gorm.ErrInvalidDB

	txErr := &TransactionError{
		OriginalError: originalErr,
		RollbackError: rollbackErr,
	}

	errMsg := txErr.Error()
	assert.Contains(t, errMsg, "transaction failed")
	assert.Contains(t, errMsg, originalErr.Error())
	assert.Contains(t, errMsg, rollbackErr.Error())
}

func TestTxManager_WithinTransaction_WithError(t *testing.T) {
	db := setupTestDB(t)
	txManager := NewTxManager(db)

	ctx := context.Background()
	err := txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		return assert.AnError
	})
	assert.Error(t, err)
}
