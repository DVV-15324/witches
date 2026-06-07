package test

import (
	"context"
	"net/http"
	"testing"

	"core-v/pkg/core/handle/server"

	"github.com/gin-gonic/gin"
)

func TestHttpShutDown(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	server.ShadownServer(
		context.Background(),
		mux,
		8080,
	)
}

func TestGin(t *testing.T) {
	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello gin",
		})
	})

	server.ShadownServer(
		context.Background(),
		r,
		8080,
	)
}
