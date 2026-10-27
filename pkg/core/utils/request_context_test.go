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
func TestSetExtra(t *testing.T) {
	reqCtx := NewRequestContext("user-123", "tid-456")

	type myExtra struct {
		TenantID string
	}

	reqCtx.SetExtra(&myExtra{TenantID: "tenant-1"})

	got, ok := reqCtx.RQCxt.(*myExtra)
	assert.True(t, ok)
	assert.Equal(t, "tenant-1", got.TenantID)
}

func TestGetExtra(t *testing.T) {
	cfg := &wcmd_utils.Config{
		RequestKey: "request-context",
	}

	type myExtra struct {
		TenantID string
	}

	t.Run("nil RQCxt", func(t *testing.T) {
		reqCtx := NewRequestContext("user-123", "tid-456")
		ctx := SaveRequestContext(context.Background(), reqCtx, cfg)

		got, ok := GetExtra[*myExtra](ctx, cfg)

		assert.False(t, ok)
		assert.Nil(t, got)
	})

	t.Run("wrong type", func(t *testing.T) {
		reqCtx := NewRequestContext("user-123", "tid-456")
		reqCtx.SetExtra("not-a-struct")
		ctx := SaveRequestContext(context.Background(), reqCtx, cfg)

		got, ok := GetExtra[*myExtra](ctx, cfg)

		assert.False(t, ok)
		assert.Nil(t, got)
	})

	t.Run("valid type", func(t *testing.T) {
		reqCtx := NewRequestContext("user-123", "tid-456")
		reqCtx.SetExtra(&myExtra{TenantID: "tenant-1"})
		ctx := SaveRequestContext(context.Background(), reqCtx, cfg)

		got, ok := GetExtra[*myExtra](ctx, cfg)

		assert.True(t, ok)
		assert.Equal(t, "tenant-1", got.TenantID)
	})

	t.Run("no request context in ctx", func(t *testing.T) {
		ctx := context.Background()

		got, ok := GetExtra[*myExtra](ctx, cfg)

		assert.False(t, ok)
		assert.Nil(t, got)
	})
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
