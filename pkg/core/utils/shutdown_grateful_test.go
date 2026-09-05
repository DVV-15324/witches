package utils

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShutdownServer(t *testing.T) {
	originalNotify := notifySignal
	originalShutdown := shutdownServerFunc

	t.Cleanup(func() {
		notifySignal = originalNotify
		shutdownServerFunc = originalShutdown
	})

	notifySignal = func(c chan<- os.Signal, sig ...os.Signal) {
		c <- os.Interrupt
	}

	shutdownServerFunc = func(
		server *http.Server,
		ctx context.Context,
	) error {
		return nil
	}

	assert.NotPanics(t, func() {
		ShutdownServer(
			context.Background(),
			http.NewServeMux(),
			"127.0.0.1",
			"0",
		)
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

	server := &http.Server{}

	shutdownServer(server)
}

func TestShutdownServer_Success(t *testing.T) {
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

	server := &http.Server{}

	shutdownServer(server)
}
