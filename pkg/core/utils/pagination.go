package utils

import (
	"gorm.io/gorm"
)

type PaginationRequest struct {
	Page   int    `form:"page" binding:"omitempty,min=1"`
	Limit  int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Search string `form:"search"`
	Sort   string `form:"sort"`
	Order  string `form:"order"`
}

func (p *PaginationRequest) GetPage() int {
	if p.Page < 1 {
		return 1
	}
	return p.Page
}

func (p *PaginationRequest) GetLimit() int {
	if p.Limit < 1 {
		return 10
	}
	if p.Limit > 100 {
		return 100
	}
	return p.Limit
}

func (p *PaginationRequest) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *PaginationRequest) GetSort(defaultSort string) string {
	if p.Sort == "" {
		return defaultSort
	}
	return p.Sort
}

func (p *PaginationRequest) GetOrder(defaultOrder string) string {
	if p.Order == "" {
		return defaultOrder
	}
	if p.Order != "asc" && p.Order != "desc" {
		return defaultOrder
	}
	return p.Order
}

func (p *PaginationRequest) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit < 1 {
		p.Limit = 10
	}
	if p.Limit > 100 {
		p.Limit = 100
	}
	if p.Sort == "" {
		p.Sort = "created_at"
	}
	if p.Order == "" || (p.Order != "asc" && p.Order != "desc") {
		p.Order = "desc"
	}
}

func (p *PaginationRequest) TotalPages(total int64) int {
	limit := p.GetLimit()
	if limit <= 0 {
		limit = 10
	}
	return int((total + int64(limit) - 1) / int64(limit))
}

func (p *PaginationRequest) HasNext(total int64) bool {
	return p.GetPage() < p.TotalPages(total)
}

func (p *PaginationRequest) HasPrev() bool {
	return p.GetPage() > 1
}

func (p *PaginationRequest) GetNextPage(total int64) int {
	if p.HasNext(total) {
		return p.GetPage() + 1
	}
	return p.GetPage()
}

func (p *PaginationRequest) GetPrevPage() int {
	if p.HasPrev() {
		return p.GetPage() - 1
	}
	return 1
}

func (p *PaginationRequest) GormScope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		p.Normalize()
		return db.Offset(p.GetOffset()).Limit(p.GetLimit())
	}
}

type PaginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
	HasNext    bool  `json:"hasNext"`
	HasPrev    bool  `json:"hasPrev"`
	NextPage   int   `json:"nextPage,omitempty"`
	PrevPage   int   `json:"prevPage,omitempty"`
}

func (p *PaginationRequest) ToPaginationResponse(total int64) PaginationResponse {
	return PaginationResponse{
		Page:       p.GetPage(),
		Limit:      p.GetLimit(),
		Total:      total,
		TotalPages: p.TotalPages(total),
		HasNext:    p.HasNext(total),
		HasPrev:    p.HasPrev(),
		NextPage:   p.GetNextPage(total),
		PrevPage:   p.GetPrevPage(),
	}
}
