package test

import (
	"context"
	"net/http"
	"testing"
	"time"

	utils "github.com/DVV-15324/witches/pkg/core/utils"

	"github.com/gin-gonic/gin"
)

func TestHttpShutDown(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	utils.ShutdownServer(
		ctx,
		mux,
		"localhost",
		"8084",
	)
}

func TestGin(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		c.String(200, "hello gin")
	})

	utils.ShutdownServer(
		ctx,
		r,
		"localhost",
		"8085",
	)
}
