package test

import (
	"context"
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
		repo.Create(&User{
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
		"localhost:8080",
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

	// TEST CASES

	t.Run("GET /api/v1/users - pagination test", func(t *testing.T) {
		// Test page 2, limit 5
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
		t.Logf("Pagination: Page=%d, Limit=%d, Total=%d", resp.Pagination.Page, resp.Pagination.Limit, resp.Pagination.Total)
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
		assert.NotNil(t, resp.Data)
		t.Logf("Search results: %d items", resp.Pagination.Total)
	})

	t.Run("GET /api/v1/users/:id", func(t *testing.T) {
		// Get first user
		req, _ := http.NewRequest("GET", "/api/v1/users?limit=1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var resp response.PaginationResponseWrapper
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		// Kiểm tra Data không nil
		if resp.Data == nil {
			t.Skip("No users found (Data is nil)")
			return
		}

		// Kiểm tra Data là slice
		users, ok := resp.Data.([]interface{})
		if !ok || len(users) == 0 {
			t.Skip("No users found or Data is not a slice")
			return
		}

		user := users[0].(map[string]interface{})
		id := user["id"].(string)

		// Get user by ID
		req2, _ := http.NewRequest("GET", "/api/v1/users/"+id, nil)
		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, req2)

		assert.Equal(t, 200, w2.Code)

		var resp2 response.AppResponse
		err = json.Unmarshal(w2.Body.Bytes(), &resp2)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp2.Status)
		assert.NotNil(t, resp2.Data)
		t.Logf("User found: %s", id)
	})

	t.Run("GET /api/v1/users/:id", func(t *testing.T) {
		// Get first user
		req, _ := http.NewRequest("GET", "/api/v1/users?limit=1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var resp response.PaginationResponseWrapper
		json.Unmarshal(w.Body.Bytes(), &resp)

		users := resp.Data.([]interface{})
		if len(users) == 0 {
			t.Skip("No users found")
		}
		user := users[0].(map[string]interface{})
		id := user["id"].(string)

		// Get user by ID
		req2, _ := http.NewRequest("GET", "/api/v1/users/"+id, nil)
		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, req2)

		assert.Equal(t, 200, w2.Code)

		var resp2 response.AppResponse
		err := json.Unmarshal(w2.Body.Bytes(), &resp2)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp2.Status)
		assert.NotNil(t, resp2.Data)
		t.Logf("User found: %s", id)
	})

	t.Run("PUT /api/v1/users/:id", func(t *testing.T) {
		// Get first user
		req, _ := http.NewRequest("GET", "/api/v1/users?limit=1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var resp response.PaginationResponseWrapper
		json.Unmarshal(w.Body.Bytes(), &resp)

		users := resp.Data.([]interface{})
		if len(users) == 0 {
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

		var resp2 response.AppResponse
		err := json.Unmarshal(w2.Body.Bytes(), &resp2)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp2.Status)
		assert.NotNil(t, resp2.Data)
		t.Logf("User updated: %s", id)
	})

	t.Run("GET /swagger/index.html - Swagger UI", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/swagger/index.html", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		t.Log("Swagger UI loaded")
	})

	srv := &http.Server{Addr: ":8080", Handler: r}
	go srv.ListenAndServe()
	defer srv.Shutdown(context.Background())
}
