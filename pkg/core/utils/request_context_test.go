package utils

import (
	"context"
	"github.com/DVV-15324/witches/cmd/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getTestConfigContext() *utils.Config {
	return &utils.Config{
		RequestKey: "request-context",
	}
}

func TestSaveAndGetRequestContext(t *testing.T) {
	cfg := getTestConfigContext()
	ctx := context.Background()
	reqCtx := NewRequestContext("sub-123", "tid-456")
	ctx = SaveRequestContext(ctx, reqCtx, cfg)

	retrieved := GetRequestContext(ctx, cfg)
	assert.Equal(t, reqCtx.Sub, retrieved.Sub)
	assert.Equal(t, reqCtx.Tid, retrieved.Tid)
}

func TestGetSub(t *testing.T) {
	cfg := getTestConfigContext()
	ctx := context.Background()
	reqCtx := NewRequestContext("sub-123", "tid-456")
	ctx = SaveRequestContext(ctx, reqCtx, cfg)

	assert.Equal(t, "sub-123", GetSub(ctx, cfg))

	emptyCtx := context.Background()
	assert.Empty(t, GetSub(emptyCtx, cfg))
}

func TestGetTid(t *testing.T) {
	cfg := getTestConfigContext()
	ctx := context.Background()
	reqCtx := NewRequestContext("sub-123", "tid-456")
	ctx = SaveRequestContext(ctx, reqCtx, cfg)

	assert.Equal(t, "tid-456", GetTid(ctx, cfg))

	emptyCtx := context.Background()
	assert.Empty(t, GetTid(emptyCtx, cfg))
}

func TestGetRequestContext_WithInvalidType(t *testing.T) {
	cfg := getTestConfigContext()
	ctx := context.WithValue(context.Background(), cfg.RequestKey, "invalid-value")
	retrieved := GetRequestContext(ctx, cfg)
	assert.NotNil(t, retrieved)
	assert.Empty(t, retrieved.Sub)
	assert.Empty(t, retrieved.Tid)
}

func TestGetRequestContext_WithNil(t *testing.T) {
	cfg := getTestConfigContext()
	ctx := context.Background()
	retrieved := GetRequestContext(ctx, cfg)
	assert.NotNil(t, retrieved)
	assert.Empty(t, retrieved.Sub)
}
