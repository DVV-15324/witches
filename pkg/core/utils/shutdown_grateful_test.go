package utils

import (
	"context"

	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShutdownServer_ListenAndServeError(t *testing.T) {
	originalListen := listenAndServe
	originalNotify := notifySignal

	done := make(chan struct{})

	listenAndServe = func(server *http.Server) error {
		defer close(done)
		return assert.AnError
	}

	notifySignal = func(ch chan<- os.Signal, sig ...os.Signal) {
		ch <- os.Interrupt
	}

	ShutdownServer(
		context.Background(),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
		"127.0.0.1",
		"0",
	)

	// đợi goroutine đọc listenAndServe xong.
	<-done

	// Sau đó mới restore global.
	listenAndServe = originalListen
	notifySignal = originalNotify
}

func TestShutdownServerFunc(t *testing.T) {
	server := &http.Server{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := shutdownServerFunc(server, ctx)

	assert.NoError(t, err)
}
func TestShutdownServer_ListenAndServeClosed(t *testing.T) {
	originalListen := listenAndServe
	originalNotify := notifySignal

	done := make(chan struct{})

	listenAndServe = func(server *http.Server) error {
		defer close(done)
		return http.ErrServerClosed
	}

	notifySignal = func(ch chan<- os.Signal, sig ...os.Signal) {
		ch <- os.Interrupt
	}

	ShutdownServer(
		context.Background(),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
		"127.0.0.1",
		"0",
	)

	<-done

	listenAndServe = originalListen
	notifySignal = originalNotify
}

func TestShutdownServer_ListenAndServeNil(t *testing.T) {
	originalListen := listenAndServe
	originalNotify := notifySignal

	done := make(chan struct{})

	listenAndServe = func(server *http.Server) error {
		defer close(done)
		return nil
	}

	notifySignal = func(ch chan<- os.Signal, sig ...os.Signal) {
		ch <- os.Interrupt
	}

	ShutdownServer(
		context.Background(),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
		"127.0.0.1",
		"0",
	)

	<-done

	listenAndServe = originalListen
	notifySignal = originalNotify
}
func TestShutdownServer_Success(t *testing.T) {
	server := &http.Server{}

	assert.NotPanics(t, func() {
		shutdownServer(server)
	})
}
func TestShutdownServer_ShutdownError(t *testing.T) {
	original := shutdownServerFunc
	t.Cleanup(func() {
		shutdownServerFunc = original
	})

	shutdownServerFunc = func(
		server *http.Server,
		ctx context.Context,
	) error {
		return assert.AnError
	}

	shutdownServer(&http.Server{})
}
func TestShutdownServer_ShutdownSuccess(t *testing.T) {
	original := shutdownServerFunc
	t.Cleanup(func() {
		shutdownServerFunc = original
	})

	shutdownServerFunc = func(
		server *http.Server,
		ctx context.Context,
	) error {
		return nil
	}

	shutdownServer(&http.Server{})
}
func TestListenAndServe(t *testing.T) {
	server := &http.Server{
		Addr: "127.0.0.1:0",
	}

	done := make(chan error, 1)

	go func() {
		done <- listenAndServe(server)
	}()

	time.Sleep(100 * time.Millisecond)

	err := server.Close()
	assert.ErrorIs(t, err, nil)

	select {
	case err := <-done:
		assert.ErrorIs(t, err, http.ErrServerClosed)
	case <-time.After(time.Second):
		t.Fatal("listenAndServe did not stop")
	}
}
