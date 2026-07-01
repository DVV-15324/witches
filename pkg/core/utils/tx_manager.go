// internal/utils/tx.go
package utils

import (
	"context"
	"gorm.io/gorm"
)

// TxManager là interface cho transaction (để dễ mock khi test)
type TxManager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// txManager implements TxManager
type txManager struct {
	db *gorm.DB
}

// NewTxManager khởi tạo tx manager với db connection
func NewTxManager(db *gorm.DB) TxManager {
	return &txManager{db: db}
}

// contextKey là kiểu riêng để tránh xung đột key trong context
type contextKey struct{}

// txKey là key để lưu transaction vào context
var txKey = contextKey{}

// WithinTransaction chạy function trong 1 transaction
// Nếu fn trả về error -> rollback, ngược lại -> commit
func (t *txManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	// Bắt đầu transaction (GORM)
	tx := t.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Lưu transaction vào context để các repository dùng chung
	txCtx := context.WithValue(ctx, txKey, tx)

	// Chạy function (chứa logic nghiệp vụ)
	err := fn(txCtx)

	if err != nil {
		// Có lỗi -> rollback
		if rbErr := tx.Rollback().Error; rbErr != nil {
			// Nếu rollback cũng lỗi -> ghi log hoặc wrap lỗi
			return &TransactionError{
				OriginalError: err,
				RollbackError: rbErr,
			}
		}
		return err
	}

	// Không lỗi -> commit
	if commitErr := tx.Commit().Error; commitErr != nil {
		return commitErr
	}

	return nil
}

// GetTxFromContext lấy transaction từ context (dùng trong repository)
func GetTxFromContext(ctx context.Context) (*gorm.DB, error) {
	tx, ok := ctx.Value(txKey).(*gorm.DB)
	if !ok {
		return nil, ErrNoTransaction
	}
	return tx, nil
}

// TransactionError wrap lỗi rollback và lỗi chính
type TransactionError struct {
	OriginalError error
	RollbackError error
}

func (e *TransactionError) Error() string {
	return "transaction failed: " + e.OriginalError.Error() + " (rollback error: " + e.RollbackError.Error() + ")"
}

// ErrNoTransaction thông báo khi context không chứa transaction
var ErrNoTransaction = &NoTransactionError{}

type NoTransactionError struct{}

func (e *NoTransactionError) Error() string {
	return "no transaction found in context"
}
