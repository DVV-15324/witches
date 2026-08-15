//go:build integration

package utils

import (
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestInitMetric(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	gin.SetMode(gin.TestMode)
	engine := gin.Default()
	engine.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	port := "8083"
	host := "localhost"
	InitMetric(port, host, engine)

	time.Sleep(200 * time.Millisecond)

	// Test metrics endpoint
	resp, err := http.Get("http://localhost:8083/metrics")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	resp.Body.Close()
}
