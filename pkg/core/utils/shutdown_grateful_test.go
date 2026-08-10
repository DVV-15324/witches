// pkg/core/utils/shutdown_grateful_test.go
//go:build integration

package utils

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestShutdownServer(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	gin.SetMode(gin.TestMode)
	engine := gin.Default()
	engine.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		ShutdownServer(ctx, engine, "localhost", "8081")
	}()

	time.Sleep(200 * time.Millisecond)

	resp, err := http.Get("http://localhost:8081/test")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	resp.Body.Close()

	cancel()
	time.Sleep(100 * time.Millisecond)
}
