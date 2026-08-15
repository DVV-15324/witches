package utils

import (
	"context"
	"testing"

	"github.com/DVV-15324/witches/cmd/utils"
	"github.com/stretchr/testify/assert"
)

func getTestConfigContext() *utils.Config {
	return &utils.Config{
		RequestKey: "request-context",
	}
}

func TestSaveAndGetRequestContext(t *testing.T) {
	cfg := getTestConfigContext()
	ctx := context.Background()
	reqCtx := NewRequestContext("sub-123", "tid-456", "device-789", "192.168.1.1", "Mozilla/5.0")
	ctx = SaveRequestContext(ctx, reqCtx, cfg)

	retrieved := GetRequestContext(ctx, cfg)
	assert.Equal(t, reqCtx.Sub, retrieved.Sub)
	assert.Equal(t, reqCtx.Tid, retrieved.Tid)
	assert.Equal(t, reqCtx.DeviceID, retrieved.DeviceID)
	assert.Equal(t, reqCtx.IPAddress, retrieved.IPAddress)
	assert.Equal(t, reqCtx.UserAgent, retrieved.UserAgent)
}

func TestGetSub(t *testing.T) {
	cfg := getTestConfigContext()
	ctx := context.Background()
	reqCtx := NewRequestContext("sub-123", "tid-456", "device-789", "192.168.1.1", "Mozilla/5.0")
	ctx = SaveRequestContext(ctx, reqCtx, cfg)

	assert.Equal(t, "sub-123", GetSub(ctx, cfg))

	emptyCtx := context.Background()
	assert.Empty(t, GetSub(emptyCtx, cfg))
}

func TestGetTid(t *testing.T) {
	cfg := getTestConfigContext()
	ctx := context.Background()
	reqCtx := NewRequestContext("sub-123", "tid-456", "device-789", "192.168.1.1", "Mozilla/5.0")
	ctx = SaveRequestContext(ctx, reqCtx, cfg)

	assert.Equal(t, "tid-456", GetTid(ctx, cfg))

	emptyCtx := context.Background()
	assert.Empty(t, GetTid(emptyCtx, cfg))
}

func TestGetDeviceID(t *testing.T) {
	cfg := getTestConfigContext()
	ctx := context.Background()
	reqCtx := NewRequestContext("sub-123", "tid-456", "device-789", "192.168.1.1", "Mozilla/5.0")
	ctx = SaveRequestContext(ctx, reqCtx, cfg)

	assert.Equal(t, "device-789", GetDeviceID(ctx, cfg))

	emptyCtx := context.Background()
	assert.Empty(t, GetDeviceID(emptyCtx, cfg))
}

func TestGetIPAddress(t *testing.T) {
	cfg := getTestConfigContext()
	ctx := context.Background()
	reqCtx := NewRequestContext("sub-123", "tid-456", "device-789", "192.168.1.1", "Mozilla/5.0")
	ctx = SaveRequestContext(ctx, reqCtx, cfg)

	assert.Equal(t, "192.168.1.1", GetIPAddress(ctx, cfg))

	emptyCtx := context.Background()
	assert.Empty(t, GetIPAddress(emptyCtx, cfg))
}

func TestGetUserAgent(t *testing.T) {
	cfg := getTestConfigContext()
	ctx := context.Background()
	reqCtx := NewRequestContext("sub-123", "tid-456", "device-789", "192.168.1.1", "Mozilla/5.0")
	ctx = SaveRequestContext(ctx, reqCtx, cfg)

	assert.Equal(t, "Mozilla/5.0", GetUserAgent(ctx, cfg))

	emptyCtx := context.Background()
	assert.Empty(t, GetUserAgent(emptyCtx, cfg))
}

func TestGetSessionID(t *testing.T) {
	cfg := getTestConfigContext()
	ctx := context.Background()
	reqCtx := NewRequestContextFull(
		"sub-123", "tid-456", "device-789",
		"192.168.1.1", "Mozilla/5.0",
		"shard-1", "session-123", "req-456",
		"web", "vi-VN", "Asia/Ho_Chi_Minh",
	)
	ctx = SaveRequestContext(ctx, reqCtx, cfg)

	assert.Equal(t, "session-123", GetSessionID(ctx, cfg))

	emptyCtx := context.Background()
	assert.Empty(t, GetSessionID(emptyCtx, cfg))
}

func TestGetPlatform(t *testing.T) {
	cfg := getTestConfigContext()
	ctx := context.Background()
	reqCtx := NewRequestContextFull(
		"sub-123", "tid-456", "device-789",
		"192.168.1.1", "Mozilla/5.0",
		"shard-1", "session-123", "req-456",
		"web", "vi-VN", "Asia/Ho_Chi_Minh",
	)
	ctx = SaveRequestContext(ctx, reqCtx, cfg)

	assert.Equal(t, "web", GetPlatform(ctx, cfg))

	emptyCtx := context.Background()
	assert.Empty(t, GetPlatform(emptyCtx, cfg))
}

func TestGetLocale(t *testing.T) {
	cfg := getTestConfigContext()
	ctx := context.Background()
	reqCtx := NewRequestContextFull(
		"sub-123", "tid-456", "device-789",
		"192.168.1.1", "Mozilla/5.0",
		"shard-1", "session-123", "req-456",
		"web", "vi-VN", "Asia/Ho_Chi_Minh",
	)
	ctx = SaveRequestContext(ctx, reqCtx, cfg)

	assert.Equal(t, "vi-VN", GetLocale(ctx, cfg))

	emptyCtx := context.Background()
	assert.Empty(t, GetLocale(emptyCtx, cfg))
}

func TestGetTimezone(t *testing.T) {
	cfg := getTestConfigContext()
	ctx := context.Background()
	reqCtx := NewRequestContextFull(
		"sub-123", "tid-456", "device-789",
		"192.168.1.1", "Mozilla/5.0",
		"shard-1", "session-123", "req-456",
		"web", "vi-VN", "Asia/Ho_Chi_Minh",
	)
	ctx = SaveRequestContext(ctx, reqCtx, cfg)

	assert.Equal(t, "Asia/Ho_Chi_Minh", GetTimezone(ctx, cfg))

	emptyCtx := context.Background()
	assert.Empty(t, GetTimezone(emptyCtx, cfg))
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
