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
	"github.com/redis/go-redis/v9"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/ulule/limiter/v3"
	redis_driver "github.com/ulule/limiter/v3/drivers/store/redis"
)

func resetRedis(mr *miniredis.Miniredis) {
	mr.FlushAll()
}
func TestHandleSwaggerRateLimit(t *testing.T) {

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Setup DB
	db := initDB()
	repo = NewUserRepository(db)

	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer mr.Close()
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

	//  Redis client (in-memory cho test)
	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	durationRateLimitPeriod := time.Duration(60) * time.Second
	rate := limiter.Rate{
		Period: durationRateLimitPeriod,
		Limit:  5,
	}

	store, _ := redis_driver.NewStore(redisClient)
	rateLimiter := limiter.New(store, rate)

	//  Middleware rate limit
	r.Use(func(c *gin.Context) {
		ctx := c.Request.Context()
		limiterCtx, err := rateLimiter.Get(ctx, c.ClientIP())
		if err != nil {
			c.JSON(500, gin.H{"error": "rate limit error"})
			c.Abort()
			return
		}
		if limiterCtx.Reached {
			c.JSON(429, gin.H{"error": "too many requests"})
			c.Abort()
			return
		}
		c.Next()
	})

	// Swagger Generator
	gen := w_handler.NewSwaggerGenerator(
		"User API",
		"1.0",
		"localhost:8082",
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
		RateLimit(rate).
		QueryParam("page", "Page number", false).
		QueryParam("limit", "Items per page", false).
		QueryParam("search", "Search by name or email", false).
		Response(200, response.PaginationResponseWrapper{}, "Success").
		Response(429, nil, "Too Many Requests").
		Response(500, response.AppResponse{}, "Internal Server Error").
		Handler(GetUsers).
		Build()

	gen.GET("/api/v1/users/:id").
		Summary("Get user by ID").
		Description("Returns a single user by ID").
		Tags("users").
		RateLimit(rate).
		PathParam("id", "User ID", true).
		Response(200, User{}, "Success").
		Response(404, response.AppResponse{}, "Not Found").
		Handler(GetUserByID).
		Build()

	gen.POST("/api/v1/users").
		Summary("Create user").
		Description("Create a new user").
		Tags("users").
		RateLimit(rate).
		Body(CreateUserRequest{}, "User data").
		Response(201, User{}, "Created").
		Response(400, response.AppResponse{}, "Bad Request").
		Handler(CreateUser).
		Build()

	gen.PUT("/api/v1/users/:id").
		Summary("Update user").
		Description("Update an existing user").
		Tags("users").
		RateLimit(rate).
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
		RateLimit(rate).
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
	t.Run("GET /api/v1/users - invalid query params", func(t *testing.T) {
		resetRedis(mr)
		tests := []struct {
			query string
			want  int
		}{
			{"?page=0&limit=10", 200},   // page=0 → default to 1
			{"?page=-1&limit=10", 400},  // page=-1 → invalid (binding: min=1)
			{"?page=1&limit=0", 200},    // limit=0 → default to 10
			{"?page=1&limit=1000", 400}, // limit=1000 → exceeds max=100
		}
		for _, tt := range tests {
			req, _ := http.NewRequest("GET", "/api/v1/users"+tt.query, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.want, w.Code)
		}
	})
	// TEST CHO GET /api/v1/users/:id
	t.Run("GET /api/v1/users/:id - success", func(t *testing.T) {
		resetRedis(mr)
		// 1. Tạo user trước
		createReq := `{"name":"Test User","email":"test@example.com","password":"123456","age":25}`
		req, _ := http.NewRequest("POST", "/api/v1/users", strings.NewReader(createReq))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var resp response.AppResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		user := resp.Data.(map[string]interface{})
		id := user["id"].(string)

		// 2. Get user by ID
		req2, _ := http.NewRequest("GET", "/api/v1/users/"+id, nil)
		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, req2)

		assert.Equal(t, 200, w2.Code)
		t.Logf("Get user by ID success: %s", id)
	})

	t.Run("GET /api/v1/users/:id - not found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/users/non-existent-id", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
		t.Log("User not found - correct 404")
	})

	// TEST CHO POST /api/v1/users
	t.Run("POST /api/v1/users - create success", func(t *testing.T) {
		resetRedis(mr)
		reqBody := `{"name":"New User","email":"newuser@example.com","password":"123456","age":25}`
		req, _ := http.NewRequest("POST", "/api/v1/users", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		t.Log("User created successfully")
	})

	t.Run("POST /api/v1/users - duplicate email", func(t *testing.T) {
		resetRedis(mr)
		reqBody := `{"name":"Duplicate","email":"duplicate@example.com","password":"123456","age":25}`
		req, _ := http.NewRequest("POST", "/api/v1/users", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Try create duplicate
		req2, _ := http.NewRequest("POST", "/api/v1/users", strings.NewReader(reqBody))
		req2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, req2)

		// Should return error (400 or 409)
		assert.NotEqual(t, 200, w2.Code)
		t.Log("Duplicate email handled correctly")
	})

	t.Run("POST /api/v1/users - invalid data", func(t *testing.T) {
		resetRedis(mr)
		reqBody := `{"name":"","email":"invalid","password":"","age":-1}`
		req, _ := http.NewRequest("POST", "/api/v1/users", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		t.Log("Invalid data rejected")
	})

	// TEST CHO PUT /api/v1/users/:id
	t.Run("PUT /api/v1/users/:id - update success", func(t *testing.T) {
		resetRedis(mr)
		// Create user first
		createReq := `{"name":"Update User","email":"update@example.com","password":"123456","age":25}`
		req, _ := http.NewRequest("POST", "/api/v1/users", strings.NewReader(createReq))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var resp response.AppResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		user := resp.Data.(map[string]interface{})
		id := user["id"].(string)

		// Update user
		updateReq := `{"name":"Updated Name","age":99}`
		req2, _ := http.NewRequest("PUT", "/api/v1/users/"+id, strings.NewReader(updateReq))
		req2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, req2)

		assert.Equal(t, 200, w2.Code)
		t.Logf("User updated: %s", id)
	})

	t.Run("PUT /api/v1/users/:id - not found", func(t *testing.T) {
		resetRedis(mr)
		reqBody := `{"name":"Updated"}`
		req, _ := http.NewRequest("PUT", "/api/v1/users/non-existent-id", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
		t.Log("Update non-existent user - 404")
	})

	// TEST CHO DELETE /api/v1/users/:id
	t.Run("DELETE /api/v1/users/:id - delete success", func(t *testing.T) {
		resetRedis(mr)
		// Create user first
		createReq := `{"name":"Delete User","email":"delete@example.com","password":"123456","age":25}`
		req, _ := http.NewRequest("POST", "/api/v1/users", strings.NewReader(createReq))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var resp response.AppResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		user := resp.Data.(map[string]interface{})
		id := user["id"].(string)

		// Delete user
		req2, _ := http.NewRequest("DELETE", "/api/v1/users/"+id, nil)
		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, req2)

		assert.Equal(t, 200, w2.Code)
		t.Logf("User deleted: %s", id)
	})

	t.Run("DELETE /api/v1/users/:id - not found", func(t *testing.T) {
		resetRedis(mr)
		req, _ := http.NewRequest("DELETE", "/api/v1/users/non-existent-id", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
		t.Log("Delete non-existent user - 404")
	})

}
