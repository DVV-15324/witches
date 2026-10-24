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

	cfg := &utils.Config{} // Cấu hình rỗng (mặc định)

	router := gin.New()
	router.Use(Cors(cfg))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	req.Header.Set("Origin", "http://localhost:60560")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "same-origin", rec.Header().Get("Cross-Origin-Opener-Policy"))

	assert.Equal(t, "http://localhost:60560", rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "Origin", rec.Header().Get("Vary")) // Bắt buộc phải có Vary khi Echo Origin

	assert.Equal(t, "GET, POST, PUT, DELETE, PATCH, OPTIONS", rec.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Content-Type, Authorization", rec.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "true", rec.Header().Get("Access-Control-Allow-Credentials"))
}

func TestCors_CustomConfig_OPTIONS_Allowed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &utils.Config{
		CorsAllowOrigins: "https://example.com,http://localhost:3000",
		CorsAllowMethods: "GET, POST",
		CorsAllowHeaders: "X-Custom-Header",
	}

	router := gin.New()
	router.Use(Cors(cfg))
	router.OPTIONS("/test", func(c *gin.Context) {
		t.Fatal("handler should not be called for OPTIONS request")
	})

	req := httptest.NewRequest(http.MethodOptions, "/test", nil)

	req.Header.Set("Origin", "https://example.com")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, "same-origin", rec.Header().Get("Cross-Origin-Opener-Policy"))

	assert.Equal(t, "https://example.com", rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST", rec.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "X-Custom-Header", rec.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "true", rec.Header().Get("Access-Control-Allow-Credentials"))
}

func TestCors_CustomConfig_OPTIONS_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &utils.Config{
		CorsAllowOrigins: "https://example.com",
	}

	router := gin.New()
	router.Use(Cors(cfg))
	router.OPTIONS("/test", func(c *gin.Context) {
		t.Fatal("handler should not be called")
	})

	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	req.Header.Set("Origin", "https://hacker-site.com")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)

	assert.Empty(t, rec.Header().Get("Access-Control-Allow-Origin"))
}
