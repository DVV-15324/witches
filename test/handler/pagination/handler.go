package test

import (
	"fmt"
	"net/http"

	"time"

	"github.com/DVV-15324/witches/pkg/core/response"
	"github.com/DVV-15324/witches/pkg/core/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ENTITY

type User struct {
	ID        string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name      string    `json:"name" binding:"required"`
	Email     string    `json:"email" binding:"required,email" gorm:"uniqueIndex"`
	Password  string    `json:"password" binding:"required,min=6"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

// DTO

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Age      int    `json:"age"`
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email" binding:"omitempty,email"`
	Age   int    `json:"age"`
}

// PAGINATION

type PaginationRequest struct {
	Page   int    `form:"page" binding:"omitempty,min=1"`
	Limit  int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Search string `form:"search"`
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

func (p *PaginationRequest) TotalPages(total int64) int {
	limit := p.GetLimit()
	return int((total + int64(limit) - 1) / int64(limit))
}

// REPOSITORY

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindAll(req *PaginationRequest) ([]User, int64, error) {
	var users []User
	var total int64

	query := r.db.Model(&User{})
	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where("name LIKE ? OR email LIKE ?", search, search)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := req.GetOffset()
	limit := req.GetLimit()
	if err := query.Offset(offset).Limit(limit).Order("created_at desc").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) FindByID(id string) (*User, error) {
	var user User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&User{}).Error
}

// HANDLERS (Dùng response của Witches)

var repo *UserRepository

func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database: " + err.Error())
	}
	db.AutoMigrate(&User{})
	return db
}

func setupRepo() {
	if repo == nil {
		db := initDB()
		repo = NewUserRepository(db)
		// Seed 50 users
		for i := 0; i < 50; i++ {
			repo.Create(&User{
				ID:        uuid.New().String(),
				Name:      fmt.Sprintf("User %d", i),
				Email:     fmt.Sprintf("user%d@example.com", i),
				Password:  "123456",
				Age:       20 + i%30,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
		}
	}
}

// GetUsers - Lấy danh sách user với pagination
func GetUsers(c *gin.Context) {
	setupRepo()

	var req PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.WriteError(c, &response.AppError{
			Status:    http.StatusBadRequest,
			Error:     err,
			TimeStamp: time.Now(),
		})
		return
	}

	users, total, err := repo.FindAll(&req)
	if err != nil {
		response.WriteError(c, &response.AppError{
			Status:    http.StatusInternalServerError,
			Error:     err,
			TimeStamp: time.Now(),
		})
		return
	}

	// Tạo pagination response
	pagination := &utils.PaginationResponse{
		Page:       req.GetPage(),
		Limit:      req.GetLimit(),
		Total:      total,
		TotalPages: req.TotalPages(total),
		HasNext:    req.GetPage() < req.TotalPages(total),
		HasPrev:    req.GetPage() > 1,
	}

	// Dùng response của Witches
	response.WriteSuccessWithPagination(c, users, pagination)
}

// GetUserByID - Lấy user theo ID
func GetUserByID(c *gin.Context) {
	setupRepo()
	id := c.Param("id")

	user, err := repo.FindByID(id)
	if err != nil {
		response.WriteError(c, &response.AppError{
			Status:    http.StatusNotFound,
			Error:     fmt.Errorf("user not found"),
			TimeStamp: time.Now(),
		})
		return
	}

	response.WriteSuccess(c, user)
}

// CreateUser - Tạo user mới
func CreateUser(c *gin.Context) {
	setupRepo()

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, &response.AppError{
			Status:    http.StatusBadRequest,
			Error:     err,
			TimeStamp: time.Now(),
		})
		return
	}

	user := &User{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  req.Password,
		Age:       req.Age,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := repo.Create(user); err != nil {
		response.WriteError(c, &response.AppError{
			Status:    http.StatusInternalServerError,
			Error:     err,
			TimeStamp: time.Now(),
		})
		return
	}

	response.WriteSuccess(c, user)
}

// UpdateUser - Cập nhật user
func UpdateUser(c *gin.Context) {
	setupRepo()
	id := c.Param("id")

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, &response.AppError{
			Status:    http.StatusBadRequest,
			Error:     err,
			TimeStamp: time.Now(),
		})
		return
	}

	user, err := repo.FindByID(id)
	if err != nil {
		response.WriteError(c, &response.AppError{
			Status:    http.StatusNotFound,
			Error:     fmt.Errorf("user not found"),
			TimeStamp: time.Now(),
		})
		return
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Age > 0 {
		user.Age = req.Age
	}
	user.UpdatedAt = time.Now()

	if err := repo.Update(user); err != nil {
		response.WriteError(c, &response.AppError{
			Status:    http.StatusInternalServerError,
			Error:     err,
			TimeStamp: time.Now(),
		})
		return
	}

	response.WriteSuccess(c, user)
}

// DeleteUser - Xóa user
func DeleteUser(c *gin.Context) {
	setupRepo()
	id := c.Param("id")

	if err := repo.Delete(id); err != nil {
		response.WriteError(c, &response.AppError{
			Status:    http.StatusNotFound,
			Error:     fmt.Errorf("user not found"),
			TimeStamp: time.Now(),
		})
		return
	}

	response.WriteSuccess(c, gin.H{"message": "User deleted successfully"})
}
