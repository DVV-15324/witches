package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapPtrSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    []*int
		mapper   func(*int) *string
		expected []*string
	}{
		{
			name:  "normal case",
			input: []*int{ptr(1), ptr(2), ptr(3)},
			mapper: func(i *int) *string {
				if i == nil {
					return nil
				}
				s := "num-" + string(rune('0'+*i))
				return &s
			},
			expected: []*string{ptr("num-1"), ptr("num-2"), ptr("num-3")},
		},
		{
			name:     "nil input",
			input:    nil,
			mapper:   func(i *int) *string { return nil },
			expected: nil,
		},
		{
			name:     "empty slice",
			input:    []*int{},
			mapper:   func(i *int) *string { return nil },
			expected: []*string{},
		},
		{
			name:  "with nil elements",
			input: []*int{ptr(1), nil, ptr(3)},
			mapper: func(i *int) *string {
				if i == nil {
					return nil
				}
				s := "val-" + string(rune('0'+*i))
				return &s
			},
			expected: []*string{ptr("val-1"), nil, ptr("val-3")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapPtrSlice(tt.input, tt.mapper)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapPtrSlice_WithStruct(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}
	type UserDTO struct {
		ID   int
		Name string
	}

	input := []*User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}
	mapper := func(u *User) *UserDTO {
		if u == nil {
			return nil
		}
		return &UserDTO{
			ID:   u.ID,
			Name: u.Name,
		}
	}
	expected := []*UserDTO{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}

	result := MapPtrSlice(input, mapper)
	assert.Equal(t, expected, result)
}

// Helper để tạo pointer
func ptr[T any](v T) *T {
	return &v
}
