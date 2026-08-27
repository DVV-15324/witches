package handle

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func mockStoreFactory() (limiter.Store, error) {
	return memory.NewStore(), nil
}
func mockStoreFactoryError() (limiter.Store, error) {
	return nil, assert.AnError
}

type mockRateLimitMiddleware struct{}

func (m *mockRateLimitMiddleware) CreateRateLimitMiddleware(rateLimit *limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func setupTest() (*SwaggerGenerator, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	gen := NewSwaggerGenerator("Test API", "1.0.0", "localhost:8080", "/api/v1")
	gen.SetEngine(engine)
	gen.storeFactory = mockStoreFactory

	gen.SetRateLimitMiddleware(&mockRateLimitMiddleware{})
	return gen, engine
}

func TestRouteBuilder_Build_GET(t *testing.T) {
	gen, engine := setupTest()

	builder := gen.GET("/users/:id").
		Summary("Get user by ID").
		Description("Returns user details").
		Tags("users").
		PathParam("id", "User ID", true).
		QueryParam("include", "Include fields", false).
		Response(200, map[string]interface{}{"id": 1, "name": "John"}, "Success").
		Response(404, nil, "Not found").
		Handler(func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"id": 1, "name": "John"})
		})

	builder.Build()

	req := httptest.NewRequest("GET", "/api/v1/users/1", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "John")

	assert.NotNil(t, gen.doc.Paths["/users/:id"])
	assert.Contains(t, gen.doc.Paths["/users/:id"], "get")
}

func TestRouteBuilder_Build_POST(t *testing.T) {
	gen, engine := setupTest()

	type CreateUserRequest struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	builder := gen.POST("/users").
		Summary("Create user").
		Tags("users").
		Body(CreateUserRequest{}, "User data").
		Response(201, map[string]interface{}{"id": 1}, "Created").
		Response(400, nil, "Bad request").
		Handler(func(c *gin.Context) {
			c.JSON(http.StatusCreated, gin.H{"id": 1})
		})

	builder.Build()

	req := httptest.NewRequest("POST", "/api/v1/users", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestRouteBuilder_RateLimit(t *testing.T) {
	gen, engine := setupTest()

	rate := limiter.Rate{
		Period: 1 * time.Second,
		Limit:  1,
	}

	builder := gen.GET("/limited").
		Summary("Limited endpoint").
		RateLimit(rate).
		Handler(func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

	builder.Build()

	req := httptest.NewRequest("GET", "/api/v1/limited", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

}

func TestRouteBuilder_RateLimit_NilStore(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	gen := NewSwaggerGenerator("Test", "1.0", "localhost", "/api/v1")
	gen.SetEngine(engine)
	gen.storeFactory = func() (limiter.Store, error) {
		return nil, nil
	}

	builder := gen.GET("/test").
		RateLimit(limiter.Rate{Period: 1 * time.Second, Limit: 1}).
		Handler(func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})

	assert.NotPanics(t, func() {
		builder.Build()
	})
}

func TestRouteBuilder_RateLimit_StoreError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	gen := NewSwaggerGenerator("Test", "1.0", "localhost", "/api/v1")
	gen.SetEngine(engine)
	gen.storeFactory = mockStoreFactoryError

	builder := gen.GET("/test").
		RateLimit(limiter.Rate{Period: 1 * time.Second, Limit: 1}).
		Handler(func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})

	assert.NotPanics(t, func() {
		builder.Build()
	})
}

func TestRouteBuilder_UseMiddleware(t *testing.T) {
	gen, engine := setupTest()

	var middlewareCalled bool
	middleware := func(c *gin.Context) {
		middlewareCalled = true
		c.Next()
	}

	builder := gen.GET("/middleware").
		Use(middleware).
		Handler(func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})

	builder.Build()

	req := httptest.NewRequest("GET", "/api/v1/middleware", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	assert.True(t, middlewareCalled)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRouteBuilder_HeaderParam(t *testing.T) {
	gen, _ := setupTest()

	builder := gen.GET("/test").
		HeaderParam("X-API-Key", "API key", true).
		Handler(func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})

	builder.Build()

	pathItem := gen.doc.Paths["/test"]
	assert.NotNil(t, pathItem)
	op := pathItem["get"]
	assert.NotNil(t, op)

	var hasHeaderParam bool
	for _, param := range op.Parameters {
		if param.In == "header" && param.Name == "X-API-Key" {
			hasHeaderParam = true
			assert.True(t, param.Required)
		}
	}
	assert.True(t, hasHeaderParam)
}

func TestRouteBuilder_BodyModel(t *testing.T) {
	gen, _ := setupTest()

	type TestModel struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	builder := gen.POST("/test").
		Body(TestModel{}, "Test body").
		Handler(func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})

	builder.Build()

	pathItem := gen.doc.Paths["/test"]
	assert.NotNil(t, pathItem)
	op := pathItem["post"]
	assert.NotNil(t, op)

	var hasBodyParam bool
	for _, param := range op.Parameters {
		if param.In == "body" {
			hasBodyParam = true
			assert.NotNil(t, param.Schema)
			assert.Contains(t, param.Schema.Ref, "TestModel")
		}
	}
	assert.True(t, hasBodyParam)

	_, ok := gen.modelParser.schemas["TestModel"]
	assert.True(t, ok)
}

func TestRouteBuilder_PathParam(t *testing.T) {
	gen, _ := setupTest()

	builder := gen.GET("/users/:id").
		PathParam("id", "User ID", true).
		Handler(func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"id": c.Param("id")})
		})

	builder.Build()

	pathItem := gen.doc.Paths["/users/:id"]
	assert.NotNil(t, pathItem)
	op := pathItem["get"]
	assert.NotNil(t, op)

	var hasPathParam bool
	for _, param := range op.Parameters {
		if param.In == "path" && param.Name == "id" {
			hasPathParam = true
			assert.True(t, param.Required)
		}
	}
	assert.True(t, hasPathParam)
}

func TestRouteBuilder_QueryParam(t *testing.T) {
	gen, _ := setupTest()

	builder := gen.GET("/search").
		QueryParam("q", "Search query", true).
		QueryParam("limit", "Result limit", false).
		Handler(func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})

	builder.Build()

	pathItem := gen.doc.Paths["/search"]
	assert.NotNil(t, pathItem)
	op := pathItem["get"]
	assert.NotNil(t, op)

	var queryParams []Parameter
	for _, param := range op.Parameters {
		if param.In == "query" {
			queryParams = append(queryParams, param)
		}
	}
	assert.Len(t, queryParams, 2)
	assert.Equal(t, "q", queryParams[0].Name)
	assert.True(t, queryParams[0].Required)
	assert.Equal(t, "limit", queryParams[1].Name)
	assert.False(t, queryParams[1].Required)
}

func TestRouteBuilder_ResponseWithStringModel(t *testing.T) {
	gen, _ := setupTest()

	builder := gen.GET("/text").
		Response(200, "text/plain", "Plain text response").
		Handler(func(c *gin.Context) {
			c.String(http.StatusOK, "ok")
		})

	builder.Build()

	pathItem := gen.doc.Paths["/text"]
	assert.NotNil(t, pathItem)
	op := pathItem["get"]
	assert.NotNil(t, op)

	resp, ok := op.Responses["200"]
	assert.True(t, ok)
	assert.Equal(t, "Plain text response", resp.Description)
}

func TestRouteBuilder_AllMethods(t *testing.T) {
	gen, engine := setupTest()

	methods := []struct {
		method string
		fn     func(string) *RouteBuilder
	}{
		{"GET", gen.GET},
		{"POST", gen.POST},
		{"PUT", gen.PUT},
		{"DELETE", gen.DELETE},
		{"PATCH", gen.PATCH},
		{"OPTIONS", gen.OPTIONS},
	}

	for _, m := range methods {
		t.Run(m.method, func(t *testing.T) {
			builder := m.fn("/test")
			builder.Handler(func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"method": m.method})
			})
			builder.Build()

			req := httptest.NewRequest(m.method, "/api/v1/test", nil)
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Contains(t, w.Body.String(), m.method)
		})
	}
}

func TestRouteBuilder_SwaggerGeneration(t *testing.T) {
	gen, _ := setupTest()

	type UserResponse struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	gen.GET("/users/:id").
		Summary("Get user").
		Description("Get user by ID").
		Tags("users", "v1").
		PathParam("id", "User ID", true).
		Response(200, UserResponse{}, "Success").
		Handler(func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"id": 1, "name": "John"})
		}).
		Build()

	doc := gen.doc
	assert.NotNil(t, doc.Paths["/users/:id"])
	assert.Contains(t, doc.Paths["/users/:id"], "get")

	op := doc.Paths["/users/:id"]["get"]
	assert.Equal(t, "Get user", op.Summary)
	assert.Equal(t, "Get user by ID", op.Description)
	assert.Contains(t, op.Tags, "users")
	assert.Contains(t, op.Tags, "v1")

	_, ok := gen.modelParser.schemas["UserResponse"]
	assert.True(t, ok)
}

func BenchmarkRouteBuilder_Build(b *testing.B) {
	gen, _ := setupTest()

	b.ResetTimer()
	for b.Loop() {
		builder := gen.GET("/test").
			Summary("Test").
			Handler(func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"ok": true})
			})
		builder.Build()
	}
}
