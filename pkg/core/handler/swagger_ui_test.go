package handle

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSwaggerUI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	engine.GET("/swagger", SwaggerUI())

	req := httptest.NewRequest("GET", "/swagger", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "text/html")
	assert.Contains(t, w.Body.String(), "swagger-ui")
	assert.Contains(t, w.Body.String(), "Swagger UI")
}
