package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	w_handler "github.com/DVV-15324/witches/pkg/core/handler"
	"github.com/DVV-15324/witches/pkg/core/response"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

func TestHandleSwaggerPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Setup DB
	db := initDB()
	repo = NewUserRepository(db)

	// Seed data
	for i := 0; i < 50; i++ {
		_ = repo.Create(&User{
			ID:        uuid.New().String(),
			Name:      fmt.Sprintf("Test User %d", i),
			Email:     fmt.Sprintf("test%d@example.com", i),
			Password:  "123456",
			Age:       20 + i%30,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	// Swagger Generator
	gen := w_handler.NewSwaggerGenerator(
		"User API",
		"1.0",
		"localhost:8081",
		"/",
	)
	gen.SetEngine(r)

	gen.AddTag("users", "User management")
	gen.RegisterModel(User{})
	gen.RegisterModel(CreateUserRequest{})
	gen.RegisterModel(response.AppResponse{})

	// Routes
	gen.GET("/api/v1/users").
		Summary("List all users").
		Description("Returns a list of all users with pagination").
		Tags("users").
		QueryParam("page", "Page number", false).
		QueryParam("limit", "Items per page", false).
		QueryParam("search", "Search by name or email", false).
		Response(200, response.PaginationResponseWrapper{}, "Success").
		Response(500, response.AppResponse{}, "Internal Server Error").
		Handler(GetUsers).
		Build()

	gen.GET("/api/v1/users/:id").
		Summary("Get user by ID").
		Description("Returns a single user by ID").
		Tags("users").
		PathParam("id", "User ID", true).
		Response(200, User{}, "Success").
		Response(404, response.AppResponse{}, "Not Found").
		Handler(GetUserByID).
		Build()

	gen.POST("/api/v1/users").
		Summary("Create user").
		Description("Create a new user").
		Tags("users").
		Body(CreateUserRequest{}, "User data").
		Response(201, User{}, "Created").
		Response(400, response.AppResponse{}, "Bad Request").
		Handler(CreateUser).
		Build()

	gen.PUT("/api/v1/users/:id").
		Summary("Update user").
		Description("Update an existing user").
		Tags("users").
		PathParam("id", "User ID", true).
		Body(UpdateUserRequest{}, "User data").
		Response(200, User{}, "Updated").
		Response(400, response.AppResponse{}, "Bad Request").
		Response(404, response.AppResponse{}, "Not Found").
		Handler(UpdateUser).
		Build()

	gen.DELETE("/api/v1/users/:id").
		Summary("Delete user").
		Description("Delete a user by ID").
		Tags("users").
		PathParam("id", "User ID", true).
		Response(200, response.AppResponse{}, "Deleted").
		Response(404, response.AppResponse{}, "Not Found").
		Handler(DeleteUser).
		Build()

	// Save swagger.json
	if err := gen.Save("swagger.json"); err != nil {
		t.Logf("Lỗi lưu swagger.json: %v", err)
	} else {
		t.Log("swagger.json generated!")
	}

	r.GET("/swagger/*any", w_handler.SwaggerUI())
	r.GET("/swagger.json", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.String(200, gen.GenerateJSON())
	})

	// ========== TEST CASES ==========

	t.Run("GET /api/v1/users - pagination test", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/users?page=2&limit=5", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		var resp response.PaginationResponseWrapper
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.Status)
		assert.Equal(t, "Success", resp.Message)
		assert.NotNil(t, resp.Data)
		assert.NotNil(t, resp.Pagination)
		assert.Equal(t, 2, resp.Pagination.Page)
		assert.Equal(t, 5, resp.Pagination.Limit)
		assert.Equal(t, int64(50), resp.Pagination.Total)
	})

	t.Run("GET /api/v1/users - search test", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/users?search=User 1&limit=10", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		var resp response.PaginationResponseWrapper
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.Status)
	})

	// THÊM TEST: Invalid query params
	t.Run("GET /api/v1/users - invalid query params", func(t *testing.T) {
		tests := []struct {
			query string
			want  int
		}{
			{"?page=0&limit=10", 200},   // page=0 → default 1
			{"?page=-1&limit=10", 400},  // page=-1 → invalid (min=1)
			{"?page=1&limit=0", 200},    // limit=0 → default 10
			{"?page=1&limit=1000", 400}, // limit=1000 → invalid (max=100)
		}
		for _, tt := range tests {
			req, _ := http.NewRequest("GET", "/api/v1/users"+tt.query, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.want, w.Code)
		}
	})

	// THÊM TEST: Page beyond total
	t.Run("GET /api/v1/users - page beyond total", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/users?page=100&limit=10", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		var resp response.PaginationResponseWrapper
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Empty(t, resp.Data) // Should return empty data
	})

	// THÊM TEST: Search with special characters
	t.Run("GET /api/v1/users - search with special chars", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/users?search=Test%20User%201", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	// THÊM TEST: Empty search
	t.Run("GET /api/v1/users - empty search", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/users?search=", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	// THÊM TEST: GET by ID - success
	t.Run("GET /api/v1/users/:id - success", func(t *testing.T) {
		// Get first user
		req, _ := http.NewRequest("GET", "/api/v1/users?limit=1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var resp response.PaginationResponseWrapper
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		users, ok := resp.Data.([]interface{})
		if !ok || len(users) == 0 {
			t.Skip("No users found")
		}

		user := users[0].(map[string]interface{})
		id := user["id"].(string)

		// Get user by ID
		req2, _ := http.NewRequest("GET", "/api/v1/users/"+id, nil)
		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, req2)

		assert.Equal(t, 200, w2.Code)
	})

	// THÊM TEST: GET by ID - not found
	t.Run("GET /api/v1/users/:id - not found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/users/non-existent-id", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
	})

	// THÊM TEST: POST - create success
	t.Run("POST /api/v1/users - create success", func(t *testing.T) {
		reqBody := `{"name":"New User","email":"newuser@example.com","password":"123456","age":25}`
		req, _ := http.NewRequest("POST", "/api/v1/users", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		var resp response.AppResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.Status)
		assert.NotNil(t, resp.Data)
	})

	// THÊM TEST: POST - invalid data
	t.Run("POST /api/v1/users - invalid data", func(t *testing.T) {
		reqBody := `{"name":"","email":"invalid","password":"","age":-1}`
		req, _ := http.NewRequest("POST", "/api/v1/users", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})

	// THÊM TEST: POST - duplicate email
	t.Run("POST /api/v1/users - duplicate email", func(t *testing.T) {
		reqBody := `{"name":"Duplicate","email":"duplicate@example.com","password":"123456","age":25}`
		req, _ := http.NewRequest("POST", "/api/v1/users", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Try duplicate
		req2, _ := http.NewRequest("POST", "/api/v1/users", strings.NewReader(reqBody))
		req2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, req2)

		assert.NotEqual(t, 200, w2.Code)
	})

	// THÊM TEST: PUT - update success
	t.Run("PUT /api/v1/users/:id - update success", func(t *testing.T) {
		// Get first user
		req, _ := http.NewRequest("GET", "/api/v1/users?limit=1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var resp response.PaginationResponseWrapper
		_ = json.Unmarshal(w.Body.Bytes(), &resp)

		users, ok := resp.Data.([]interface{})
		if !ok || len(users) == 0 {
			t.Skip("No users found")
		}
		user := users[0].(map[string]interface{})
		id := user["id"].(string)

		// Update user
		reqBody := `{"name":"Updated Name","age":99}`
		req2, _ := http.NewRequest("PUT", "/api/v1/users/"+id, strings.NewReader(reqBody))
		req2.Header.Set("Content-Type", "application/json")

		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, req2)

		assert.Equal(t, 200, w2.Code)
	})

	// THÊM TEST: PUT - not found
	t.Run("PUT /api/v1/users/:id - not found", func(t *testing.T) {
		reqBody := `{"name":"Updated"}`
		req, _ := http.NewRequest("PUT", "/api/v1/users/non-existent-id", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
	})

	// THÊM TEST: DELETE - delete success
	t.Run("DELETE /api/v1/users/:id - delete success", func(t *testing.T) {
		// Create user first
		createReq := `{"name":"Delete User","email":"delete@example.com","password":"123456","age":25}`
		req, _ := http.NewRequest("POST", "/api/v1/users", strings.NewReader(createReq))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var resp response.AppResponse
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		user := resp.Data.(map[string]interface{})
		id := user["id"].(string)

		// Delete user
		req2, _ := http.NewRequest("DELETE", "/api/v1/users/"+id, nil)
		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, req2)

		assert.Equal(t, 200, w2.Code)
	})

	// THÊM TEST: DELETE - not found
	t.Run("DELETE /api/v1/users/:id - not found", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/v1/users/non-existent-id", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
	})

	// THÊM TEST: Swagger UI
	t.Run("GET /swagger/index.html - Swagger UI", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/swagger/index.html", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})
	t.Run("GET /api/v1/users - search no results", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/users?search=nonexistent", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		var resp response.PaginationResponseWrapper
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Empty(t, resp.Data)
	})

	t.Run("GET /api/v1/users - invalid page parameter (non-int)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/users?page=abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)
	})

	t.Run("GET /api/v1/users - invalid limit parameter (non-int)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/users?limit=xyz", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)
	})
}
