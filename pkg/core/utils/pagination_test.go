package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// pkg/core/utils/pagination_test.go
func TestPaginationRequest_GetSort(t *testing.T) {
	p := &PaginationRequest{}
	assert.Equal(t, "created_at", p.GetSort("created_at"))

	p.Sort = "name"
	assert.Equal(t, "name", p.GetSort("created_at"))

	p.Sort = "email"
	assert.Equal(t, "email", p.GetSort("created_at"))
}

func TestPaginationRequest_GetOrder(t *testing.T) {
	p := &PaginationRequest{}
	assert.Equal(t, "desc", p.GetOrder("desc"))

	p.Order = "asc"
	assert.Equal(t, "asc", p.GetOrder("desc"))

	p.Order = "invalid"
	assert.Equal(t, "desc", p.GetOrder("desc"))

	p.Order = "ASC"
	assert.Equal(t, "desc", p.GetOrder("desc")) // case sensitive
}

func TestPaginationRequest_Normalize(t *testing.T) {
	// Test với giá trị mặc định
	p := &PaginationRequest{
		Page:  0,
		Limit: 0,
		Sort:  "",
		Order: "",
	}
	p.Normalize()
	assert.Equal(t, 1, p.Page)
	assert.Equal(t, 10, p.Limit)
	assert.Equal(t, "created_at", p.Sort)
	assert.Equal(t, "desc", p.Order)

	// Test với limit > 100
	p = &PaginationRequest{
		Page:  2,
		Limit: 200,
		Sort:  "name",
		Order: "asc",
	}
	p.Normalize()
	assert.Equal(t, 2, p.Page)
	assert.Equal(t, 100, p.Limit)
	assert.Equal(t, "name", p.Sort)
	assert.Equal(t, "asc", p.Order)

	// Test với page < 1
	p = &PaginationRequest{
		Page:  -1,
		Limit: 20,
		Sort:  "email",
		Order: "asc",
	}
	p.Normalize()
	assert.Equal(t, 1, p.Page)
	assert.Equal(t, 20, p.Limit)
	assert.Equal(t, "email", p.Sort)
	assert.Equal(t, "asc", p.Order)

	// Test với limit < 1
	p = &PaginationRequest{
		Page:  3,
		Limit: -1,
		Sort:  "name",
		Order: "desc",
	}
	p.Normalize()
	assert.Equal(t, 3, p.Page)
	assert.Equal(t, 10, p.Limit)
	assert.Equal(t, "name", p.Sort)
	assert.Equal(t, "desc", p.Order)
}

func TestPaginationRequest_GetNextPage(t *testing.T) {
	p := &PaginationRequest{Page: 2, Limit: 10}
	assert.Equal(t, 3, p.GetNextPage(100))
	assert.Equal(t, 2, p.GetNextPage(15)) // Last page
	assert.Equal(t, 2, p.GetNextPage(0))
}

func TestPaginationRequest_GetPrevPage(t *testing.T) {
	p := &PaginationRequest{Page: 2, Limit: 10}
	assert.Equal(t, 1, p.GetPrevPage())

	p2 := &PaginationRequest{Page: 1, Limit: 10}
	assert.Equal(t, 1, p2.GetPrevPage())

	p3 := &PaginationRequest{Page: 0, Limit: 10}
	assert.Equal(t, 1, p3.GetPrevPage())
}

// pkg/core/utils/pagination_test.go
func TestPaginationRequest_GetLimit(t *testing.T) {
	tests := []struct {
		name  string
		limit int
		want  int
	}{
		{"valid limit", 20, 20},
		{"min limit", 1, 1},
		{"zero limit", 0, 10},
		{"negative limit", -1, 10},
		{"max limit", 100, 100},
		{"exceed max", 200, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PaginationRequest{Limit: tt.limit}
			assert.Equal(t, tt.want, p.GetLimit())
		})
	}
}

func TestPaginationRequest_GetOffset(t *testing.T) {
	tests := []struct {
		name  string
		page  int
		limit int
		want  int
	}{
		{"page 1 limit 10", 1, 10, 0},
		{"page 2 limit 10", 2, 10, 10},
		{"page 3 limit 5", 3, 5, 10},
		{"page 0 limit 10", 0, 10, 0},
		{"page 1 limit 0", 1, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PaginationRequest{Page: tt.page, Limit: tt.limit}
			assert.Equal(t, tt.want, p.GetOffset())
		})
	}
}

func TestPaginationRequest_ToPaginationResponse(t *testing.T) {
	// Test với page 1
	p1 := &PaginationRequest{Page: 1, Limit: 10}
	resp1 := p1.ToPaginationResponse(100)
	assert.Equal(t, 1, resp1.Page)
	assert.Equal(t, 10, resp1.Limit)
	assert.Equal(t, int64(100), resp1.Total)
	assert.Equal(t, 10, resp1.TotalPages)
	assert.True(t, resp1.HasNext)
	assert.False(t, resp1.HasPrev)

	// Test với page 10 (last page)
	p2 := &PaginationRequest{Page: 10, Limit: 10}
	resp2 := p2.ToPaginationResponse(100)
	assert.Equal(t, 10, resp2.Page)
	assert.False(t, resp2.HasNext)
	assert.True(t, resp2.HasPrev)

	// Test với total = 0
	p3 := &PaginationRequest{Page: 1, Limit: 10}
	resp3 := p3.ToPaginationResponse(0)
	assert.Equal(t, 0, resp3.TotalPages)
	assert.False(t, resp3.HasNext)
	assert.False(t, resp3.HasPrev)
}

// pkg/core/utils/pagination_test.go
func TestPaginationRequest_TotalPages(t *testing.T) {
	tests := []struct {
		name  string
		limit int
		total int64
		want  int
	}{
		{"exact division", 10, 100, 10},
		{"not exact", 10, 95, 10},
		{"empty", 10, 0, 0},
		{"limit 0", 0, 100, 10}, // default limit 10
		{"limit 0 total 0", 0, 0, 0},
		{"limit 100 total 1000", 100, 1000, 10},
		{"limit 100 total 999", 100, 999, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PaginationRequest{Limit: tt.limit}
			assert.Equal(t, tt.want, p.TotalPages(tt.total))
		})
	}
}

func TestPaginationRequest_GormScope(t *testing.T) {
	// Test GormScope với các giá trị khác nhau
	p := &PaginationRequest{Page: 2, Limit: 10}
	scope := p.GormScope()
	assert.NotNil(t, scope)

	// Test với page và limit default
	p2 := &PaginationRequest{Page: 0, Limit: 0}
	scope2 := p2.GormScope()
	assert.NotNil(t, scope2)

	// Test với limit > 100
	p3 := &PaginationRequest{Page: 1, Limit: 200}
	scope3 := p3.GormScope()
	assert.NotNil(t, scope3)
}
