package utils

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/DVV-15324/witches/cmd/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewDeviceHelper(t *testing.T) {
	tests := []struct {
		name   string
		cfg    *utils.Config
		expect int // không dùng, chỉ để kiểm tra cơ bản
	}{
		{
			name: "with valid locales",
			cfg: &utils.Config{
				RequestKey:         "request_context",
				SupportedLanguages: []string{"en-US", "vi-VN", "fr-FR"},
			},
		},
		{
			name: "with empty locales - fallback to default",
			cfg: &utils.Config{
				RequestKey:         "request_context",
				SupportedLanguages: []string{},
			},
		},
		{
			name: "with nil config - fallback to default",
			cfg: &utils.Config{
				RequestKey:         "request_context",
				SupportedLanguages: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			helper := NewDeviceHelper(tt.cfg)
			assert.NotNil(t, helper)
			assert.NotNil(t, helper.matcher)
			if tt.cfg != nil {
				assert.Equal(t, tt.cfg, helper.config)
			} else {
				assert.Nil(t, helper.config)
			}
		})
	}
}

func TestDeviceHelper_GetDeviceInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &utils.Config{
		RequestKey:         "request_context",
		SupportedLanguages: []string{"en-US", "vi-VN"},
	}
	helper := NewDeviceHelper(cfg)

	tests := []struct {
		name              string
		setupContext      func() *gin.Context
		expectedDeviceID  string
		expectedIP        string
		expectedUserAgent string
		expectedLocale    string
		expectedTimezone  string
	}{
		{
			name: "with request context",
			setupContext: func() *gin.Context {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				c.Request = httptest.NewRequest("GET", "/", nil)
				ctx := context.WithValue(c.Request.Context(), cfg.RequestKey, &RequestContext{
					DeviceID:  "mock-device-id",
					IPAddress: "192.168.1.1",
					UserAgent: "mock-agent",
				})
				c.Request = c.Request.WithContext(ctx)
				return c
			},
			expectedDeviceID:  "mock-device-id",
			expectedIP:        "192.168.1.1",
			expectedUserAgent: "mock-agent",
			expectedLocale:    "en-US",
			expectedTimezone:  "UTC",
		},
		{
			name: "without request context - generate from IP and User-Agent",
			setupContext: func() *gin.Context {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				c.Request = httptest.NewRequest("GET", "/", nil)
				c.Request.Header.Set("User-Agent", "test-agent")
				c.Request.RemoteAddr = "10.0.0.1:1234"
				return c
			},
			expectedDeviceID:  generateDeviceID("10.0.0.1", "test-agent"),
			expectedIP:        "10.0.0.1",
			expectedUserAgent: "test-agent",
			expectedLocale:    "en-US",
			expectedTimezone:  "UTC",
		},
		{
			name: "with Accept-Language header - vi-VN",
			setupContext: func() *gin.Context {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				c.Request = httptest.NewRequest("GET", "/", nil)
				c.Request.Header.Set("Accept-Language", "vi-VN,vi;q=0.9,en;q=0.8")
				c.Request.Header.Set("User-Agent", "test-ua")
				c.Request.RemoteAddr = "1.1.1.1:1234"
				return c
			},
			expectedDeviceID:  generateDeviceID("1.1.1.1", "test-ua"),
			expectedIP:        "1.1.1.1",
			expectedUserAgent: "test-ua",
			expectedLocale:    "vi-VN",
			expectedTimezone:  "UTC",
		},
		{
			name: "with Accept-Language header - fr-FR (fallback to en-US)",
			setupContext: func() *gin.Context {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				c.Request = httptest.NewRequest("GET", "/", nil)
				c.Request.Header.Set("Accept-Language", "fr-FR,fr;q=0.9")
				c.Request.Header.Set("User-Agent", "test-ua2")
				c.Request.RemoteAddr = "2.2.2.2:1234"
				return c
			},
			expectedDeviceID:  generateDeviceID("2.2.2.2", "test-ua2"),
			expectedIP:        "2.2.2.2",
			expectedUserAgent: "test-ua2",
			expectedLocale:    "en-US", // fallback vì fr không có trong supported
			expectedTimezone:  "UTC",
		},
		{
			name: "with X-Timezone header",
			setupContext: func() *gin.Context {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				c.Request = httptest.NewRequest("GET", "/", nil)
				c.Request.Header.Set("X-Timezone", "Asia/Ho_Chi_Minh")
				c.Request.Header.Set("User-Agent", "test-ua3")
				c.Request.RemoteAddr = "3.3.3.3:1234"
				return c
			},
			expectedDeviceID:  generateDeviceID("3.3.3.3", "test-ua3"),
			expectedIP:        "3.3.3.3",
			expectedUserAgent: "test-ua3",
			expectedLocale:    "en-US",
			expectedTimezone:  "Asia/Ho_Chi_Minh",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.setupContext()
			deviceID, ip, ua, locale, tz := helper.GetDeviceInfo(c)

			assert.Equal(t, tt.expectedDeviceID, deviceID, "deviceID mismatch")
			assert.Equal(t, tt.expectedIP, ip, "IP mismatch")
			assert.Equal(t, tt.expectedUserAgent, ua, "User-Agent mismatch")
			assert.Equal(t, tt.expectedLocale, locale, "Locale mismatch")
			assert.Equal(t, tt.expectedTimezone, tz, "Timezone mismatch")
		})
	}
}

func TestExtractTokenFromHeader(t *testing.T) {
	tests := []struct {
		name        string
		header      string
		expected    string
		expectError bool
	}{
		{"valid token", "Bearer abc123", "abc123", false},
		{"empty header", "", "", true},
		{"missing Bearer prefix", "abc123", "", true},
		{"too many parts", "Bearer abc123 extra", "", true},
		{"invalid prefix", "Basic abc123", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExtractTokenFromHeader(tt.header)
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestGenerateDeviceID(t *testing.T) {
	tests := []struct {
		name      string
		ip        string
		userAgent string
	}{
		{"basic", "192.168.1.1", "test-agent"},
		{"empty user agent", "10.0.0.1", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := generateDeviceID(tt.ip, tt.userAgent)
			assert.Len(t, id, 32) // MD5 hex length
			// Kiểm tra tính nhất quán
			id2 := generateDeviceID(tt.ip, tt.userAgent)
			assert.Equal(t, id, id2)
		})
	}
}
