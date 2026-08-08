package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	ginmiddleware "github.com/slok/go-http-metrics/middleware/gin"
)

// InitMetric khởi tạo middleware metrics và chạy metrics server trên một goroutine riêng.
// Trả về gin.Engine đã được bọc middleware .
func InitMetric(metricsAddr string, engine *gin.Engine) *gin.Engine {
	// 1. Tạo middleware metrics
	mdlw := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	// 2. Gắn middleware vào Gin engine
	engine.Use(ginmiddleware.Handler("", mdlw))

	// 3. Chạy metrics server trên goroutine riêng
	go func() {
		log.Printf("Metrics server listening on %s", metricsAddr)
		if err := http.ListenAndServe(metricsAddr, promhttp.Handler()); err != nil {
			log.Printf("Metrics server error: %v", err)
		}
	}()
	return engine
}
