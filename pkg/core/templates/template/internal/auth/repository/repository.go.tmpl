package repository

import (
	"context"

	w_utils "github.com/DVV-15324/witches/pkg/core/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db        *gorm.DB
	redis     *redis.Client
	TxManager w_utils.TxManager
}

func NewAuthRepository(db *gorm.DB, txManager w_utils.TxManager, redis *redis.Client) *AuthRepository {
	return &AuthRepository{
		db:        db,
		TxManager: txManager,
		redis:     redis,
	}
}

func (a *AuthRepository) getDB(ctx context.Context) *gorm.DB {
	tx, err := w_utils.GetTxFromContext(ctx)
	if err == nil {
		return tx // Dùng transaction
	}
	return a.db // Dùng db thường
}
