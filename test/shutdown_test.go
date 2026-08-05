package test

import (
	"context"
	"net/http"
	"testing"

	utils "github.com/DVV-15324/witches/pkg/core/utils"

	"github.com/gin-gonic/gin"
)

func TestHttpShutDown(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	utils.ShutdownServer(
		context.Background(),
		mux,
		"localhost",
		"8080",
	)
}

func TestGin(t *testing.T) {
	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		c.String(200, "hello gin")
	})

	utils.ShutdownServer(
		context.Background(),
		r,
		"localhost",
		"8080",
	)
}
