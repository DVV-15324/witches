package utils

import (
	"context"
	"testing"

	wcmd_utils "github.com/DVV-15324/witches/cmd/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewRequestContext(t *testing.T) {
	reqCtx := NewRequestContext("user-123", "tid-456")

	assert.Equal(t, "user-123", reqCtx.Sub)
	assert.Equal(t, "tid-456", reqCtx.Tid)
}

func TestNewRequestContextFull(t *testing.T) {
	reqCtx := NewRequestContextFull("user-123", "tid-456")

	assert.Equal(t, "user-123", reqCtx.Sub)
	assert.Equal(t, "tid-456", reqCtx.Tid)
}

func TestSaveRequestContext(t *testing.T) {
	cfg := &wcmd_utils.Config{
		RequestKey: "request-context",
	}

	reqCtx := NewRequestContext("user-123", "tid-456")

	ctx := SaveRequestContext(
		context.Background(),
		reqCtx,
		cfg,
	)

	got := GetRequestContext(ctx, cfg)

	assert.Equal(t, reqCtx, got)
}

func TestGetRequestContext(t *testing.T) {
	cfg := &wcmd_utils.Config{
		RequestKey: "request-context",
	}

	t.Run("nil value", func(t *testing.T) {
		ctx := context.Background()

		got := GetRequestContext(ctx, cfg)

		assert.Equal(t, &RequestContext{}, got)
	})

	t.Run("invalid type", func(t *testing.T) {
		ctx := context.WithValue(
			context.Background(),
			cfg.RequestKey,
			"invalid-request-context",
		)

		got := GetRequestContext(ctx, cfg)

		assert.Equal(t, &RequestContext{}, got)
	})

	t.Run("valid request context", func(t *testing.T) {
		reqCtx := NewRequestContext("user-123", "tid-456")

		ctx := context.WithValue(
			context.Background(),
			cfg.RequestKey,
			reqCtx,
		)

		got := GetRequestContext(ctx, cfg)

		assert.Equal(t, reqCtx, got)
	})
}

func TestGetSub(t *testing.T) {
	cfg := &wcmd_utils.Config{
		RequestKey: "request-context",
	}

	reqCtx := NewRequestContext("user-123", "tid-456")
	ctx := SaveRequestContext(context.Background(), reqCtx, cfg)

	assert.Equal(t, "user-123", GetSub(ctx, cfg))
}

func TestGetTid(t *testing.T) {
	cfg := &wcmd_utils.Config{
		RequestKey: "request-context",
	}

	reqCtx := NewRequestContext("user-123", "tid-456")
	ctx := SaveRequestContext(context.Background(), reqCtx, cfg)

	assert.Equal(t, "tid-456", GetTid(ctx, cfg))
}
