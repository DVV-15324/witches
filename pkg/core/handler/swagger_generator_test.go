package handle

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

// ==================== TESTS ====================

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
	// Không test redis thật
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

func TestSwaggerGenerator_AddBearerAuth(t *testing.T) {
	gen := NewSwaggerGenerator("Test", "1.0", "localhost", "/api")

	gen.AddBearerAuth("BearerAuth")

	assert.Contains(t, gen.doc.SecurityDefinitions, "BearerAuth")
	assert.Equal(t, "apiKey", gen.doc.SecurityDefinitions["BearerAuth"].Type)
	assert.Equal(t, "header", gen.doc.SecurityDefinitions["BearerAuth"].In)
	assert.Equal(t, "Authorization", gen.doc.SecurityDefinitions["BearerAuth"].Name)
}

func TestSwaggerGenerator_AddBearerAuthWithDescription(t *testing.T) {
	gen := NewSwaggerGenerator("Test", "1.0", "localhost", "/api")

	gen.AddBearerAuthWithDescription("BearerAuth", "Custom description")

	assert.Contains(t, gen.doc.SecurityDefinitions, "BearerAuth")
	assert.Equal(t, "Custom description", gen.doc.SecurityDefinitions["BearerAuth"].Description)
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
	// ✅ Dùng setup đầy đủ
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
	// ✅ Dùng setup đầy đủ
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
