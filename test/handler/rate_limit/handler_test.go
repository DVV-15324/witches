package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

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
		Response(429, nil, "Too Many Requests").
		Response(500, response.AppResponse{}, "Internal Server Error").
		Handler(GetUsers).
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

	t.Run("Rate Limit - normal request", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/users?limit=5", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		t.Log("Normal request passes")
	})

	t.Run("Rate Limit - exceed limit", func(t *testing.T) {
		// Gửi nhiều request để vượt quá limit
		for i := 0; i < 10; i++ {
			req, _ := http.NewRequest("GET", "/api/v1/users?limit=5", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code == 429 {
				t.Logf("Rate limit triggered at request %d", i+1)
				return
			}
		}
		t.Error("Rate limit not triggered")
	})
}
