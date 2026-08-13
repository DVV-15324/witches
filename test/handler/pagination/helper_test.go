package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetLimit(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  int
	}{
		{"default when 0", 0, 10},
		{"default when negative", -5, 10},
		{"cap at 100", 150, 100},
		{"within range", 50, 50},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &PaginationRequest{Limit: tt.input}
			assert.Equal(t, tt.want, req.GetLimit())
		})
	}
}

func TestGetOffset(t *testing.T) {
	req := &PaginationRequest{Page: 3, Limit: 10}
	assert.Equal(t, 20, req.GetOffset())
}

func TestTotalPages(t *testing.T) {
	req := &PaginationRequest{Limit: 10}
	tests := []struct {
		total int64
		want  int
	}{
		{0, 0},
		{5, 1},
		{10, 1},
		{11, 2},
		{100, 10},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, req.TotalPages(tt.total))
	}
}

func TestSetupRepo(t *testing.T) {
	// Reset repo để test lại
	oldRepo := repo
	repo = nil
	setupRepo()
	assert.NotNil(t, repo)
	// Gọi lần 2 – không reset, vẫn dùng repo cũ
	setupRepo()
	assert.NotNil(t, repo)
	// Khôi phục lại (dù không cần)
	repo = oldRepo
}

func TestDelete(t *testing.T) {
	db := initDB()
	repo := NewUserRepository(db)

	// Tạo user
	user := &User{
		ID:        "test-delete-1",
		Name:      "Delete Me",
		Email:     "delete@example.com",
		Password:  "123456",
		Age:       20,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := repo.Create(user)
	assert.NoError(t, err)

	// Delete thành công
	err = repo.Delete("test-delete-1")
	assert.NoError(t, err)

	// Delete lần 2 -> lỗi not found
	err = repo.Delete("test-delete-1")
	assert.Error(t, err)
}

func TestUpdateUser_Success(t *testing.T) {
	db := initDB()
	repo := NewUserRepository(db)

	// Tạo user
	user := &User{
		ID:        "test-update-1",
		Name:      "Original",
		Email:     "original@example.com",
		Password:  "123456",
		Age:       20,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := repo.Create(user)
	assert.NoError(t, err)

	// Update chỉ name
	user.Name = "Updated Name"
	err = repo.Update(user)
	assert.NoError(t, err)

	// Lấy lại và kiểm tra
	updated, err := repo.FindByID("test-update-1")
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", updated.Name)
	assert.Equal(t, "original@example.com", updated.Email) // không đổi
	assert.Equal(t, 20, updated.Age)
}
