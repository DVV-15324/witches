package utils

import (
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/stretchr/testify/assert"
)

func resetPrometheusRegistry() {
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
	prometheus.DefaultGatherer = prometheus.DefaultRegisterer.(prometheus.Gatherer)
}
func TestInstanceMetric_PortWithoutColon(t *testing.T) {
	resetPrometheusRegistry()

	original := listenAndServe
	t.Cleanup(func() {
		listenAndServe = original
	})

	called := make(chan string, 1)

	listenAndServe = func(
		addr string,
		handler http.Handler,
	) error {
		called <- addr
		return nil
	}

	gin.SetMode(gin.TestMode)
	engine := gin.New()

	got := InstanceMetric("8083", "localhost", engine)

	assert.Same(t, engine, got)

	select {
	case addr := <-called:
		assert.Equal(t, ":8083", addr)
	case <-time.After(time.Second):
		t.Fatal("ListenAndServe was not called")
	}
}

func TestInstanceMetric_PortWithColon(t *testing.T) {
	resetPrometheusRegistry()

	original := listenAndServe
	t.Cleanup(func() {
		listenAndServe = original
	})

	called := make(chan string, 1)

	listenAndServe = func(
		addr string,
		handler http.Handler,
	) error {
		called <- addr
		return nil
	}

	engine := gin.New()

	got := InstanceMetric(":8083", "localhost", engine)

	assert.Same(t, engine, got)

	select {
	case addr := <-called:
		assert.Equal(t, ":8083", addr)
	case <-time.After(time.Second):
		t.Fatal("ListenAndServe was not called")
	}
}

func TestInstanceMetric_PortWithHost(t *testing.T) {
	resetPrometheusRegistry()

	original := listenAndServe
	t.Cleanup(func() {
		listenAndServe = original
	})

	called := make(chan string, 1)

	listenAndServe = func(
		addr string,
		handler http.Handler,
	) error {
		called <- addr
		return nil
	}

	engine := gin.New()

	got := InstanceMetric("localhost:8083", "127.0.0.1", engine)

	assert.Same(t, engine, got)

	select {
	case addr := <-called:
		assert.Equal(t, "localhost:8083", addr)
	case <-time.After(time.Second):
		t.Fatal("ListenAndServe was not called")
	}
}

func TestInstanceMetric_ListenAndServeError(t *testing.T) {
	resetPrometheusRegistry()

	original := listenAndServe
	t.Cleanup(func() {
		listenAndServe = original
	})

	done := make(chan struct{})

	listenAndServe = func(
		addr string,
		handler http.Handler,
	) error {
		defer close(done)
		return assert.AnError
	}

	engine := gin.New()

	InstanceMetric("8083", "localhost", engine)

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("ListenAndServe was not called")
	}
}
