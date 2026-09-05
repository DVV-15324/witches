package utils

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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
func TestTxManager_WithinTransaction_BeginError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(`select sqlite_version\(\)`).
		WillReturnRows(
			sqlmock.NewRows([]string{"sqlite_version"}).
				AddRow("3.45.0"),
		)

	gormDB, err := gorm.Open(
		sqlite.Dialector{Conn: db},
		&gorm.Config{},
	)
	require.NoError(t, err)

	mock.ExpectBegin().WillReturnError(assert.AnError)

	txManager := NewTxManager(gormDB)

	err = txManager.WithinTransaction(
		context.Background(),
		func(ctx context.Context) error {
			t.Fatal("callback should not be called")
			return nil
		},
	)

	assert.ErrorIs(t, err, assert.AnError)
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestTxManager_WithinTransaction_RollbackError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(`select sqlite_version\(\)`).
		WillReturnRows(
			sqlmock.NewRows([]string{"sqlite_version"}).
				AddRow("3.45.0"),
		)

	gormDB, err := gorm.Open(
		sqlite.Dialector{Conn: db},
		&gorm.Config{},
	)
	require.NoError(t, err)

	mock.ExpectBegin()
	mock.ExpectRollback().WillReturnError(gorm.ErrInvalidDB)

	txManager := NewTxManager(gormDB)

	err = txManager.WithinTransaction(
		context.Background(),
		func(ctx context.Context) error {
			return assert.AnError
		},
	)

	var txErr *TransactionError
	require.ErrorAs(t, err, &txErr)

	assert.Equal(t, assert.AnError, txErr.OriginalError)
	assert.Equal(t, gorm.ErrInvalidDB, txErr.RollbackError)
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestTxManager_WithinTransaction_CommitError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(`select sqlite_version\(\)`).
		WillReturnRows(
			sqlmock.NewRows([]string{"sqlite_version"}).
				AddRow("3.45.0"),
		)

	gormDB, err := gorm.Open(
		sqlite.Dialector{Conn: db},
		&gorm.Config{},
	)
	require.NoError(t, err)

	mock.ExpectBegin()
	mock.ExpectCommit().WillReturnError(gorm.ErrInvalidDB)

	txManager := NewTxManager(gormDB)

	err = txManager.WithinTransaction(
		context.Background(),
		func(ctx context.Context) error {
			return nil
		},
	)

	assert.ErrorIs(t, err, gorm.ErrInvalidDB)
	assert.NoError(t, mock.ExpectationsWereMet())
}
