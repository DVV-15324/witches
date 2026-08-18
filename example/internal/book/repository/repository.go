package repository

import (
	"context"

	u_tx "github.com/DVV-15324/witches/pkg/core/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type BookRepository struct {
	db        *gorm.DB
	redis     *redis.Client
	TxManager u_tx.TxManager
}

func NewBookRepository(db *gorm.DB, txManager u_tx.TxManager, redis *redis.Client) *BookRepository {
	return &BookRepository{
		db:        db,
		TxManager: txManager,
		redis:     redis,
	}
}

func (a *BookRepository) getDB(ctx context.Context) *gorm.DB {
	tx, err := u_tx.GetTxFromContext(ctx)
	if err == nil {
		return tx // Dùng transaction
	}
	return a.db // Dùng db thường
}
