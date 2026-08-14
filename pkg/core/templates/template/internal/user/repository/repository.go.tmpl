package repository

import (
	"context"
	w_utils "github.com/DVV-15324/witches/pkg/core/utils"
	"gorm.io/gorm"
)

type UserRepository struct {
	db        *gorm.DB
	TxManager w_utils.TxManager
}

func NewUserRepository(db *gorm.DB, txManager w_utils.TxManager) *UserRepository {
	return &UserRepository{
		db:        db,
		TxManager: txManager,
	}
}

func (u *UserRepository) getDB(ctx context.Context) *gorm.DB {
	tx, err := w_utils.GetTxFromContext(ctx)
	if err == nil {
		return tx // Dùng transaction
	}
	return u.db // Dùng db thường
}
