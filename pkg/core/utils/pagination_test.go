package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestPaginationRequest_GetPage(t *testing.T) {
	tests := []struct {
		name string
		page int
		want int
	}{
		{
			name: "zero",
			page: 0,
			want: 1,
		},
		{
			name: "negative",
			page: -1,
			want: 1,
		},
		{
			name: "valid",
			page: 3,
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PaginationRequest{Page: tt.page}

			assert.Equal(t, tt.want, p.GetPage())
		})
	}
}

func TestPaginationRequest_GetLimit(t *testing.T) {
	tests := []struct {
		name  string
		limit int
		want  int
	}{
		{
			name:  "zero",
			limit: 0,
			want:  10,
		},
		{
			name:  "negative",
			limit: -1,
			want:  10,
		},
		{
			name:  "valid",
			limit: 20,
			want:  20,
		},
		{
			name:  "max",
			limit: 100,
			want:  100,
		},
		{
			name:  "exceed max",
			limit: 101,
			want:  100,
		},
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
		{
			name:  "page 1 limit 10",
			page:  1,
			limit: 10,
			want:  0,
		},
		{
			name:  "page 2 limit 10",
			page:  2,
			limit: 10,
			want:  10,
		},
		{
			name:  "page 3 limit 5",
			page:  3,
			limit: 5,
			want:  10,
		},
		{
			name:  "invalid page",
			page:  0,
			limit: 10,
			want:  0,
		},
		{
			name:  "invalid limit",
			page:  2,
			limit: 0,
			want:  10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PaginationRequest{
				Page:  tt.page,
				Limit: tt.limit,
			}

			assert.Equal(t, tt.want, p.GetOffset())
		})
	}
}

func TestPaginationRequest_GetSort(t *testing.T) {
	tests := []struct {
		name        string
		sort        string
		defaultSort string
		want        string
	}{
		{
			name:        "empty",
			sort:        "",
			defaultSort: "created_at",
			want:        "created_at",
		},
		{
			name:        "custom",
			sort:        "name",
			defaultSort: "created_at",
			want:        "name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PaginationRequest{Sort: tt.sort}

			assert.Equal(t, tt.want, p.GetSort(tt.defaultSort))
		})
	}
}

func TestPaginationRequest_GetOrder(t *testing.T) {
	tests := []struct {
		name         string
		order        string
		defaultOrder string
		want         string
	}{
		{
			name:         "empty",
			order:        "",
			defaultOrder: "desc",
			want:         "desc",
		},
		{
			name:         "asc",
			order:        "asc",
			defaultOrder: "desc",
			want:         "asc",
		},
		{
			name:         "desc",
			order:        "desc",
			defaultOrder: "asc",
			want:         "desc",
		},
		{
			name:         "invalid",
			order:        "invalid",
			defaultOrder: "desc",
			want:         "desc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PaginationRequest{Order: tt.order}

			assert.Equal(t, tt.want, p.GetOrder(tt.defaultOrder))
		})
	}
}

func TestPaginationRequest_Normalize(t *testing.T) {
	tests := []struct {
		name  string
		input PaginationRequest
		want  PaginationRequest
	}{
		{
			name:  "all defaults",
			input: PaginationRequest{},
			want: PaginationRequest{
				Page:  1,
				Limit: 10,
				Sort:  "created_at",
				Order: "desc",
			},
		},
		{
			name: "negative page and limit",
			input: PaginationRequest{
				Page:  -1,
				Limit: -10,
			},
			want: PaginationRequest{
				Page:  1,
				Limit: 10,
				Sort:  "created_at",
				Order: "desc",
			},
		},
		{
			name: "limit exceeds max",
			input: PaginationRequest{
				Page:  2,
				Limit: 200,
				Sort:  "name",
				Order: "asc",
			},
			want: PaginationRequest{
				Page:  2,
				Limit: 100,
				Sort:  "name",
				Order: "asc",
			},
		},
		{
			name: "invalid order",
			input: PaginationRequest{
				Page:  2,
				Limit: 20,
				Sort:  "name",
				Order: "invalid",
			},
			want: PaginationRequest{
				Page:  2,
				Limit: 20,
				Sort:  "name",
				Order: "desc",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.input

			p.Normalize()

			assert.Equal(t, tt.want, p)
		})
	}
}

func TestPaginationRequest_TotalPages(t *testing.T) {
	tests := []struct {
		name  string
		limit int
		total int64
		want  int
	}{
		{
			name:  "exact division",
			limit: 10,
			total: 100,
			want:  10,
		},
		{
			name:  "not exact",
			limit: 10,
			total: 101,
			want:  11,
		},
		{
			name:  "empty",
			limit: 10,
			total: 0,
			want:  0,
		},
		{
			name:  "invalid limit falls back",
			limit: 0,
			total: 25,
			want:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PaginationRequest{Limit: tt.limit}

			assert.Equal(t, tt.want, p.TotalPages(tt.total))
		})
	}
}

func TestPaginationRequest_HasNext(t *testing.T) {
	tests := []struct {
		name  string
		page  int
		limit int
		total int64
		want  bool
	}{
		{
			name:  "has next",
			page:  1,
			limit: 10,
			total: 25,
			want:  true,
		},
		{
			name:  "last page",
			page:  3,
			limit: 10,
			total: 25,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PaginationRequest{
				Page:  tt.page,
				Limit: tt.limit,
			}

			assert.Equal(t, tt.want, p.HasNext(tt.total))
		})
	}
}

func TestPaginationRequest_HasPrev(t *testing.T) {
	assert.False(t, (&PaginationRequest{Page: 1}).HasPrev())
	assert.True(t, (&PaginationRequest{Page: 2}).HasPrev())
	assert.False(t, (&PaginationRequest{Page: 0}).HasPrev())
}

func TestPaginationRequest_GetNextPage(t *testing.T) {
	tests := []struct {
		name  string
		page  int
		limit int
		total int64
		want  int
	}{
		{
			name:  "has next",
			page:  1,
			limit: 10,
			total: 25,
			want:  2,
		},
		{
			name:  "no next",
			page:  3,
			limit: 10,
			total: 25,
			want:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PaginationRequest{
				Page:  tt.page,
				Limit: tt.limit,
			}

			assert.Equal(t, tt.want, p.GetNextPage(tt.total))
		})
	}
}

func TestPaginationRequest_GetPrevPage(t *testing.T) {
	assert.Equal(
		t,
		1,
		(&PaginationRequest{Page: 1}).GetPrevPage(),
	)

	assert.Equal(
		t,
		1,
		(&PaginationRequest{Page: 0}).GetPrevPage(),
	)

	assert.Equal(
		t,
		2,
		(&PaginationRequest{Page: 3}).GetPrevPage(),
	)
}

func TestPaginationRequest_GormScope(t *testing.T) {
	db, err := gorm.Open(
		sqlite.Open(":memory:"),
		&gorm.Config{},
	)
	require.NoError(t, err)

	p := &PaginationRequest{
		Page:  2,
		Limit: 20,
	}

	stmt := db.Session(&gorm.Session{
		DryRun: true,
	}).Scopes(p.GormScope()).Find(&[]struct{}{}).Statement

	sql := stmt.SQL.String()

	assert.Contains(t, sql, "LIMIT 20")
	assert.Contains(t, sql, "OFFSET 20")
}
func TestPaginationRequest_ToPaginationResponse(t *testing.T) {
	p := &PaginationRequest{
		Page:  2,
		Limit: 10,
	}

	got := p.ToPaginationResponse(25)

	want := PaginationResponse{
		Page:       2,
		Limit:      10,
		Total:      25,
		TotalPages: 3,
		HasNext:    true,
		HasPrev:    true,
		NextPage:   3,
		PrevPage:   1,
	}

	assert.Equal(t, want, got)
}
func TestPaginationRequest_GormScope_Normalize(t *testing.T) {
	db, err := gorm.Open(
		sqlite.Open(":memory:"),
		&gorm.Config{},
	)
	require.NoError(t, err)

	p := &PaginationRequest{
		Page:  0,
		Limit: 0,
	}

	db.Session(&gorm.Session{
		DryRun: true,
	}).
		Scopes(p.GormScope()).
		Find(&[]struct{}{})

	assert.Equal(t, 1, p.Page)
	assert.Equal(t, 10, p.Limit)
	assert.Equal(t, "created_at", p.Sort)
	assert.Equal(t, "desc", p.Order)
}
