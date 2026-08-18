package repository

import (
	"context"
	"example/internal/book/mapping"
	domainBook "example/internal/shared/domain"
)

func (r *BookRepository) Update(ctx context.Context, e *domainBook.Book, id int) error {
	model := mapping.FromDomainToModelBook(e)
	return r.db.WithContext(ctx).Save(model).Error
}