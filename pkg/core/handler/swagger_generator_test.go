package handle

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ulule/limiter/v3"
)

func setupTestGenerator() (*SwaggerGenerator, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	gen := NewSwaggerGenerator("Test API", "1.0.0", "localhost:8080", "/api/v1")
	gen.SetEngine(engine)
	gen.storeFactory = mockStoreFactory
	gen.SetRateLimitMiddleware(&mockRateLimitMiddleware{})
	return gen, engine
}

func TestNewSwaggerGenerator(t *testing.T) {
	gen := NewSwaggerGenerator("My API", "2.0.0", "example.com", "/api")

	assert.Equal(t, "2.0", gen.doc.Swagger)
	assert.Equal(t, "My API", gen.doc.Info.Title)
	assert.Equal(t, "2.0.0", gen.doc.Info.Version)
	assert.Equal(t, "example.com", gen.doc.Host)
	assert.Equal(t, "/api", gen.doc.BasePath)
	assert.NotNil(t, gen.doc.Paths)
	assert.NotNil(t, gen.doc.Definitions)
	assert.NotNil(t, gen.doc.SecurityDefinitions)
	assert.NotNil(t, gen.modelParser)
}

func TestSwaggerGenerator_SetEngine(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	gen := NewSwaggerGenerator("Test", "1.0", "localhost", "/api")

	gen.SetEngine(engine)
	assert.Equal(t, engine, gen.engine)
}

func TestSwaggerGenerator_SetRedisClient(t *testing.T) {
	gen := NewSwaggerGenerator("Test", "1.0", "localhost", "/api")
	assert.Nil(t, gen.storeFactory)

	// Test với nil client
	gen.SetRedisClient(nil)

	// Test với mock client (không cần redis thật)
	mockClient := &redis.Client{}
	gen.SetRedisClient(mockClient)
	assert.NotNil(t, gen.storeFactory)
}

func TestSwaggerGenerator_Use(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	gen := NewSwaggerGenerator("Test", "1.0", "localhost", "/api")
	gen.SetEngine(engine)
	gen.storeFactory = mockStoreFactory
	gen.SetRateLimitMiddleware(&mockRateLimitMiddleware{})

	var called bool
	middleware := func(c *gin.Context) {
		called = true
		c.Next()
	}

	gen.Use(middleware)

	assert.Len(t, gen.globalMiddlewares, 1)

	gen.GET("/test").
		Handler(func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		}).
		Build()

	req := httptest.NewRequest("GET", "/api/test", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSwaggerGenerator_AddTag(t *testing.T) {
	gen := NewSwaggerGenerator("Test", "1.0", "localhost", "/api")

	gen.AddTag("users", "User management")
	gen.AddTag("products", "Product management")

	assert.Len(t, gen.doc.Tags, 2)
	assert.Equal(t, "users", gen.doc.Tags[0].Name)
	assert.Equal(t, "User management", gen.doc.Tags[0].Description)
	assert.Equal(t, "products", gen.doc.Tags[1].Name)
}

func TestSwaggerGenerator_AddRoute(t *testing.T) {
	gen := NewSwaggerGenerator("Test", "1.0", "localhost", "/api")

	op := Operation{
		Summary: "Test operation",
		Responses: map[string]Response{
			"200": {Description: "OK"},
		},
	}

	gen.AddRoute("get", "/test", op)

	assert.Contains(t, gen.doc.Paths, "/test")
	assert.Contains(t, gen.doc.Paths["/test"], "get")
	assert.Equal(t, "Test operation", gen.doc.Paths["/test"]["get"].Summary)
}

func TestSwaggerGenerator_RegisterModel(t *testing.T) {
	gen := NewSwaggerGenerator("Test", "1.0", "localhost", "/api")

	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	name := gen.RegisterModel(User{})
	assert.Equal(t, "User", name)

	_, ok := gen.modelParser.schemas["User"]
	assert.True(t, ok)
}

func TestSwaggerGenerator_GenerateJSON(t *testing.T) {

	gen, _ := setupTestGenerator()

	gen.GET("/ping").
		Summary("Ping").
		Handler(func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"pong": true})
		}).
		Build()

	jsonStr := gen.GenerateJSON()
	assert.NotEmpty(t, jsonStr)

	var doc SwaggerDoc
	err := json.Unmarshal([]byte(jsonStr), &doc)
	require.NoError(t, err)
	assert.Equal(t, "2.0", doc.Swagger)
}

func TestSwaggerGenerator_Save(t *testing.T) {
	gen, _ := setupTestGenerator()

	gen.GET("/ping").
		Summary("Ping").
		Handler(func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"pong": true})
		}).
		Build()

	err := gen.Save("swagger_test.json")
	assert.NoError(t, err)
}

func TestSwaggerGenerator_GET(t *testing.T) {
	gen := NewSwaggerGenerator("Test", "1.0", "localhost", "/api")

	builder := gen.GET("/users")
	assert.NotNil(t, builder)
	assert.Equal(t, "get", builder.method)
	assert.Equal(t, "/users", builder.path)
	assert.NotNil(t, builder.op.Responses)
}

func TestSwaggerGenerator_POST(t *testing.T) {
	gen := NewSwaggerGenerator("Test", "1.0", "localhost", "/api")

	builder := gen.POST("/users")
	assert.NotNil(t, builder)
	assert.Equal(t, "post", builder.method)
	assert.Equal(t, "/users", builder.path)
	assert.Contains(t, builder.op.Consumes, "application/json")
}

func TestSwaggerGenerator_PUT(t *testing.T) {
	gen := NewSwaggerGenerator("Test", "1.0", "localhost", "/api")

	builder := gen.PUT("/users/:id")
	assert.NotNil(t, builder)
	assert.Equal(t, "put", builder.method)
	assert.Equal(t, "/users/:id", builder.path)
}

func TestSwaggerGenerator_DELETE(t *testing.T) {
	gen := NewSwaggerGenerator("Test", "1.0", "localhost", "/api")

	builder := gen.DELETE("/users/:id")
	assert.NotNil(t, builder)
	assert.Equal(t, "delete", builder.method)
}

func TestSwaggerGenerator_PATCH(t *testing.T) {
	gen := NewSwaggerGenerator("Test", "1.0", "localhost", "/api")

	builder := gen.PATCH("/users/:id")
	assert.NotNil(t, builder)
	assert.Equal(t, "patch", builder.method)
}

func TestSwaggerGenerator_OPTIONS(t *testing.T) {
	gen := NewSwaggerGenerator("Test", "1.0", "localhost", "/api")

	builder := gen.OPTIONS("/users")
	assert.NotNil(t, builder)
	assert.Equal(t, "options", builder.method)
}

// ==================== BENCHMARKS ====================

func BenchmarkSwaggerGenerator_GenerateJSON(b *testing.B) {
	gen, _ := setupTestGenerator()

	gen.GET("/ping").
		Summary("Ping").
		Handler(func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"pong": true})
		}).
		Build()

	b.ResetTimer()
	for b.Loop() {
		gen.GenerateJSON()
	}
}
func TestSwaggerGenerator_RateLimitMiddlewareNil(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	gen := NewSwaggerGenerator("Test", "1.0", "localhost", "/api/v1")
	gen.SetEngine(engine)
	gen.storeFactory = mockStoreFactory
	// Không set rateLimitMiddleware

	builder := gen.GET("/test").
		RateLimit(limiter.Rate{Period: 1 * time.Second, Limit: 1}).
		Handler(func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})

	// Không panic
	assert.NotPanics(t, func() {
		builder.Build()
	})
}
func TestSwaggerGenerator_RateLimit_GeneratorNil(t *testing.T) {
	var gen *SwaggerGenerator

	b := &RouteBuilder{
		gen: gen,
	}

	got := b.RateLimit(limiter.Rate{
		Period: time.Minute,
		Limit:  10,
	})

	assert.Same(t, b, got)
	assert.Nil(t, b.rateLimit)
}
func TestSwaggerGenerator_SetRedisClient_Nil(t *testing.T) {
	gen := NewSwaggerGenerator(
		"test",
		"1.0",
		"localhost",
		"/api",
	)

	gen.SetRedisClient(nil)

	require.NotNil(t, gen.storeFactory)

	store, err := gen.storeFactory()

	assert.Error(t, err)
	assert.Nil(t, store)
	assert.Equal(t, "redis client not configured", err.Error())
}
func TestSwaggerGenerator_GenerateJSON_WithSchemas(t *testing.T) {
	gen := NewSwaggerGenerator("Test API", "1.0", "localhost", "/api")

	type User struct {
		ID int `json:"id"`
	}

	gen.RegisterModel(User{})

	data := gen.GenerateJSON()

	var doc SwaggerDoc
	require.NoError(t, json.Unmarshal([]byte(data), &doc))

	assert.Contains(t, doc.Definitions, "User")
}
func TestSwaggerGenerator_GenerateJSON_EmptySchemas(t *testing.T) {
	gen := NewSwaggerGenerator("Test API", "1.0", "localhost", "/api")

	data := gen.GenerateJSON()

	assert.NotEmpty(t, data)
	assert.NotNil(t, gen.doc.Definitions)
	assert.Empty(t, gen.doc.Definitions)
}
func TestSwaggerGenerator_Save_GetwdError(t *testing.T) {
	original := getwd
	t.Cleanup(func() { getwd = original })

	getwd = func() (string, error) {
		return "", errors.New("getwd failed")
	}

	gen := NewSwaggerGenerator("Test", "1.0", "localhost", "/api")

	err := gen.Save("swagger.json")

	assert.EqualError(t, err, "getwd failed")
}
func TestSwaggerGenerator_Save_MkdirError(t *testing.T) {
	originalGetwd := getwd
	originalMkdirAll := mkdirAll
	t.Cleanup(func() {
		getwd = originalGetwd
		mkdirAll = originalMkdirAll
	})

	getwd = func() (string, error) {
		return "/tmp/test", nil
	}

	mkdirAll = func(string, os.FileMode) error {
		return errors.New("mkdir failed")
	}

	gen := NewSwaggerGenerator("Test", "1.0", "localhost", "/api")

	err := gen.Save("swagger.json")

	assert.EqualError(t, err, "mkdir failed")
}
