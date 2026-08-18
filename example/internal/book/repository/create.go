package repository

import (
	"context"
	domainBook "example/internal/shared/domain"
	"example/internal/book/mapping"
)

func (r *BookRepository) Create(ctx context.Context, e *domainBook.Book) (error) {
	model := mapping.FromDomainToModelBook(e)
	err := r.db.WithContext(ctx).Create(model).Error
	return err
}