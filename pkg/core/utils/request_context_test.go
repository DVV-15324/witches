package utils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// pkg/core/utils/context_test.go
func TestGetSessionID(t *testing.T) {
	ctx := context.Background()
	keyReq := "request-context"

	// Test with context
	reqCtx := NewRequestContextFull(
		"sub-123", "tid-456", "device-789",
		"192.168.1.1", "Mozilla/5.0",
		"shard-1", "session-123", "req-456",
		"web", "vi-VN", "Asia/Ho_Chi_Minh",
	)
	ctx = SaveRequestContext(ctx, reqCtx, keyReq)
	assert.Equal(t, "session-123", GetSessionID(ctx, keyReq))

	// Test without context
	emptyCtx := context.Background()
	assert.Empty(t, GetSessionID(emptyCtx, keyReq))
}

func TestGetIPAddress(t *testing.T) {
	ctx := context.Background()
	keyReq := "request-context"

	reqCtx := NewRequestContext(keyReq, "sub-123", "tid-456", "device-789", "192.168.1.1", "Mozilla/5.0")
	ctx = SaveRequestContext(ctx, reqCtx, keyReq)
	assert.Equal(t, "192.168.1.1", GetIPAddress(ctx, keyReq))

	emptyCtx := context.Background()
	assert.Empty(t, GetIPAddress(emptyCtx, keyReq))
}

func TestGetUserAgent(t *testing.T) {
	ctx := context.Background()
	keyReq := "request-context"

	reqCtx := NewRequestContext(keyReq, "sub-123", "tid-456", "device-789", "192.168.1.1", "Mozilla/5.0")
	ctx = SaveRequestContext(ctx, reqCtx, keyReq)
	assert.Equal(t, "Mozilla/5.0", GetUserAgent(ctx, keyReq))

	emptyCtx := context.Background()
	assert.Empty(t, GetUserAgent(emptyCtx, keyReq))
}

func TestGetPlatform(t *testing.T) {
	ctx := context.Background()
	keyReq := "request-context"

	reqCtx := NewRequestContextFull(
		"sub-123", "tid-456", "device-789",
		"192.168.1.1", "Mozilla/5.0",
		"shard-1", "session-123", "req-456",
		"web", "vi-VN", "Asia/Ho_Chi_Minh",
	)
	ctx = SaveRequestContext(ctx, reqCtx, keyReq)
	assert.Equal(t, "web", GetPlatform(ctx, keyReq))

	emptyCtx := context.Background()
	assert.Empty(t, GetPlatform(emptyCtx, keyReq))
}

func TestGetLocale(t *testing.T) {
	ctx := context.Background()
	keyReq := "request-context"

	reqCtx := NewRequestContextFull(
		"sub-123", "tid-456", "device-789",
		"192.168.1.1", "Mozilla/5.0",
		"shard-1", "session-123", "req-456",
		"web", "vi-VN", "Asia/Ho_Chi_Minh",
	)
	ctx = SaveRequestContext(ctx, reqCtx, keyReq)
	assert.Equal(t, "vi-VN", GetLocale(ctx, keyReq))

	emptyCtx := context.Background()
	assert.Empty(t, GetLocale(emptyCtx, keyReq))
}

func TestGetTimezone(t *testing.T) {
	ctx := context.Background()
	keyReq := "request-context"

	reqCtx := NewRequestContextFull(
		"sub-123", "tid-456", "device-789",
		"192.168.1.1", "Mozilla/5.0",
		"shard-1", "session-123", "req-456",
		"web", "vi-VN", "Asia/Ho_Chi_Minh",
	)
	ctx = SaveRequestContext(ctx, reqCtx, keyReq)
	assert.Equal(t, "Asia/Ho_Chi_Minh", GetTimezone(ctx, keyReq))

	emptyCtx := context.Background()
	assert.Empty(t, GetTimezone(emptyCtx, keyReq))
}

// pkg/core/utils/context_test.go
func TestGetSub(t *testing.T) {
	ctx := context.Background()
	keyReq := "request-context"

	// Test with context
	reqCtx := NewRequestContext(keyReq, "sub-123", "tid-456", "device-789", "192.168.1.1", "Mozilla/5.0")
	ctx = SaveRequestContext(ctx, reqCtx, keyReq)
	assert.Equal(t, "sub-123", GetSub(ctx, keyReq))

	// Test without context
	emptyCtx := context.Background()
	assert.Empty(t, GetSub(emptyCtx, keyReq))
}

func TestGetTid(t *testing.T) {
	ctx := context.Background()
	keyReq := "request-context"

	reqCtx := NewRequestContext(keyReq, "sub-123", "tid-456", "device-789", "192.168.1.1", "Mozilla/5.0")
	ctx = SaveRequestContext(ctx, reqCtx, keyReq)
	assert.Equal(t, "tid-456", GetTid(ctx, keyReq))

	emptyCtx := context.Background()
	assert.Empty(t, GetTid(emptyCtx, keyReq))
}

func TestGetDeviceID(t *testing.T) {
	ctx := context.Background()
	keyReq := "request-context"

	reqCtx := NewRequestContext(keyReq, "sub-123", "tid-456", "device-789", "192.168.1.1", "Mozilla/5.0")
	ctx = SaveRequestContext(ctx, reqCtx, keyReq)
	assert.Equal(t, "device-789", GetDeviceID(ctx, keyReq))

	emptyCtx := context.Background()
	assert.Empty(t, GetDeviceID(emptyCtx, keyReq))
}
