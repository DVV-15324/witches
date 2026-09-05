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

var listenAndServeMetrict = http.ListenAndServe

func InstanceMetric(port string, host string, engine *gin.Engine) *gin.Engine {
	if !strings.Contains(port, ":") && !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	mdlw := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	addr := fmt.Sprintf("%s%s", host, port)

	engine.Use(ginmiddleware.Handler("", mdlw))

	go func() {
		log.Printf("Metrics server listening on http://%s", addr)

		if err := listenAndServeMetrict(port, promhttp.Handler()); err != nil {
			log.Printf("Metrics server error: %v", err)
		}
	}()

	return engine
}
