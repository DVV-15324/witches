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
		// Get tx from context
		tx, err := GetTxFromContext(txCtx)
		assert.NoError(t, err)
		assert.NotNil(t, tx)

		// Do something with tx
		return nil
	})

	assert.NoError(t, err)
}

func TestTxManager_WithinTransaction_Rollback(t *testing.T) {
	db := setupTestDB(t)
	txManager := NewTxManager(db)

	ctx := context.Background()
	err := txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		// Return error to trigger rollback
		return assert.AnError
	})

	assert.Error(t, err)
}

func TestGetTxFromContext_NoTransaction(t *testing.T) {
	ctx := context.Background()
	tx, err := GetTxFromContext(ctx)
	assert.Error(t, err)
	assert.Equal(t, ErrNoTransaction, err)
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

func TestNoTransactionError(t *testing.T) {
	err := &NoTransactionError{}
	assert.Equal(t, "no transaction found in context", err.Error())
}

// pkg/core/utils/tx_manager_test.go
func TestTxManager_WithinTransaction_WithError(t *testing.T) {
	db := setupTestDB(t)
	txManager := NewTxManager(db)

	ctx := context.Background()
	err := txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		// Return error to trigger rollback
		return assert.AnError
	})
	assert.Error(t, err)
}

func TestTxManager_WithinTransaction_CommitError(t *testing.T) {
	// Test khi commit bị lỗi (khó mock, có thể skip)
	t.Skip("Skipping commit error test - requires complex mock")
}

func TestGetTxFromContext_InvalidType(t *testing.T) {
	// Test khi context có giá trị nhưng không phải *gorm.DB
	ctx := context.WithValue(context.Background(), txKey, "invalid-type")
	tx, err := GetTxFromContext(ctx)
	assert.Error(t, err)
	assert.Equal(t, ErrNoTransaction, err)
	assert.Nil(t, tx)
}
