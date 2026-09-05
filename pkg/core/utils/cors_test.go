package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DVV-15324/witches/cmd/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCors_DefaultConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &utils.Config{}

	router := gin.New()
	router.Use(Cors(cfg))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "same-origin",
		rec.Header().Get("Cross-Origin-Opener-Policy"))
	assert.Equal(t, "*",
		rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST, PUT, DELETE, PATCH, OPTIONS",
		rec.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Content-Type, Authorization",
		rec.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "true",
		rec.Header().Get("Access-Control-Allow-Credentials"))
}

func TestCors_CustomConfig_OPTIONS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &utils.Config{
		CorsAllowOrigins: "https://example.com",
		CorsAllowMethods: "GET, POST",
		CorsAllowHeaders: "X-Custom-Header",
	}

	router := gin.New()
	router.Use(Cors(cfg))
	router.OPTIONS("/test", func(c *gin.Context) {
		t.Fatal("handler should not be called")
	})

	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, "same-origin",
		rec.Header().Get("Cross-Origin-Opener-Policy"))
	assert.Equal(t, "https://example.com",
		rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST",
		rec.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "X-Custom-Header",
		rec.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "true",
		rec.Header().Get("Access-Control-Allow-Credentials"))
}
