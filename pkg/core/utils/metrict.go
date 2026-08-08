package utils

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	ginmiddleware "github.com/slok/go-http-metrics/middleware/gin"
)

// InitMetric khởi tạo middleware metrics và chạy metrics server trên một goroutine riêng.
// Trả về gin.Engine đã được bọc middleware .
func InitMetric(port string, host string, engine *gin.Engine) *gin.Engine {
	if !strings.Contains(port, ":") && !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	// 1. Tạo middleware metrics
	mdlw := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})
	addr := fmt.Sprintf("%s%s", host, port)

	// 2. Gắn middleware vào Gin engine
	engine.Use(ginmiddleware.Handler("", mdlw))

	// 3. Chạy metrics server trên goroutine riêng
	go func() {
		log.Printf("Metrics server listening on http://%s", addr)
		if err := http.ListenAndServe(port, promhttp.Handler()); err != nil {
			log.Printf("Metrics server error: %v", err)
		}
	}()
	return engine
}
