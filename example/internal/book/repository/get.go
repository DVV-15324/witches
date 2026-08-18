package repository

import (
	"context"
	"example/internal/book/mapping"
	w_utils "github.com/DVV-15324/witches/pkg/core/utils"
	domainBook "example/internal/shared/domain"
	modelBook "example/internal/book/model"
)

func (r *BookRepository) GetByID(ctx context.Context, id int) (*domainBook.Book, error) {
	var e modelBook.Book
	err := r.db.WithContext(ctx).First(&e, id).Error
	domain := mapping.FromModelToDomainBook(&e)
	return domain, err
}

func (r *BookRepository) GetAll(ctx context.Context, req *w_utils.PaginationRequest) ([]*domainBook.Book, int64, error) {
	var list []*modelBook.Book
	var total int64

	query := r.db.WithContext(ctx).Model(&modelBook.Book{})

	// Search filter
	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where("id ILIKE", search, search)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	sort := req.GetSort("created_at")
	order := req.GetOrder("desc")
	query = query.Order(sort + " " + order)

	// Apply pagination
	if err := query.Scopes(req.GormScope()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	//mapping
	domainsMapping := mapping.FromModelToDomainBookList(list)
	return domainsMapping, total, nil
}