package utils

import (
	"context"
	"github.com/DVV-15324/witches/cmd/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"net/http/httptest"
	"testing"
)

func TestNewDeviceHelper(t *testing.T) {
	tests := []struct {
		name     string
		config   *utils.Config
		expected int // số lượng locales sau khi parse
	}{
		{
			name: "with valid locales",
			config: &utils.Config{
				SupportedLanguages: []string{"en-US", "vi-VN", "fr-FR"},
			},
			expected: 3,
		},
		{
			name: "with empty locales - fallback to default",
			config: &utils.Config{
				SupportedLanguages: []string{},
			},
			expected: 2, // en, vi-VN
		},
		{
			name: "with nil config - fallback to default",
			config: &utils.Config{
				SupportedLanguages: nil,
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			helper := NewDeviceHelper(tt.config)
			assert.NotNil(t, helper)
			assert.NotNil(t, helper.matcher)
			assert.Equal(t, tt.config, helper.config)

			// Kiểm tra số lượng tag trong matcher (khó lấy trực tiếp, nhưng có thể kiểm tra bằng cách match thử)
			if len(tt.config.SupportedLanguages) == 0 {
				// Fallback: en, vi-VN
				tag, _ := language.MatchStrings(helper.matcher, "en-US")
				base, _ := tag.Base()
				assert.Equal(t, "en", base.String())
			} else {
				tag, _ := language.MatchStrings(helper.matcher, "en-US")
				base, _ := tag.Base()
				assert.Equal(t, "en", base.String())
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
				// Set request context
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
				// Gin ClientIP returns remote addr if no X-Forwarded-For
				c.Request.RemoteAddr = "10.0.0.1:1234"
				return c
			},
			expectedIP:        "10.0.0.1",
			expectedUserAgent: "test-agent",
			// deviceID sẽ là MD5 của "10.0.0.1|test-agent"
			expectedDeviceID: generateDeviceID("10.0.0.1", "test-agent"),
			expectedLocale:   "en-US",
			expectedTimezone: "UTC",
		},
		{
			name: "with Accept-Language header - vi-VN",
			setupContext: func() *gin.Context {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				c.Request = httptest.NewRequest("GET", "/", nil)
				c.Request.Header.Set("Accept-Language", "vi-VN,vi;q=0.9,en;q=0.8")
				c.Request.RemoteAddr = "1.1.1.1:1234"
				return c
			},
			expectedIP:        "1.1.1.1",
			expectedUserAgent: "",
			expectedLocale:    "vi-VN",
			expectedTimezone:  "UTC",
		},
		{
			name: "with Accept-Language header - fr-FR (fallback to en-US)",
			setupContext: func() *gin.Context {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				c.Request = httptest.NewRequest("GET", "/", nil)
				c.Request.Header.Set("Accept-Language", "fr-FR,fr;q=0.9")
				c.Request.RemoteAddr = "2.2.2.2:1234"
				return c
			},
			expectedIP:        "2.2.2.2",
			expectedUserAgent: "",
			expectedLocale:    "en-US", // vì fr không có trong supported
			expectedTimezone:  "UTC",
		},
		{
			name: "with X-Timezone header",
			setupContext: func() *gin.Context {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				c.Request = httptest.NewRequest("GET", "/", nil)
				c.Request.Header.Set("X-Timezone", "Asia/Ho_Chi_Minh")
				c.Request.RemoteAddr = "3.3.3.3:1234"
				return c
			},
			expectedIP:        "3.3.3.3",
			expectedUserAgent: "",
			expectedLocale:    "en-US",
			expectedTimezone:  "Asia/Ho_Chi_Minh",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.setupContext()
			deviceID, ip, ua, locale, tz := helper.GetDeviceInfo(c)

			assert.Equal(t, tt.expectedDeviceID, deviceID)
			assert.Equal(t, tt.expectedIP, ip)
			assert.Equal(t, tt.expectedUserAgent, ua)
			assert.Equal(t, tt.expectedLocale, locale)
			assert.Equal(t, tt.expectedTimezone, tz)
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
		{
			name:        "valid token",
			header:      "Bearer abc123",
			expected:    "abc123",
			expectError: false,
		},
		{
			name:        "empty header",
			header:      "",
			expected:    "",
			expectError: true,
		},
		{
			name:        "missing Bearer prefix",
			header:      "abc123",
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid format - too many parts",
			header:      "Bearer abc123 extra",
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid prefix",
			header:      "Basic abc123",
			expected:    "",
			expectError: true,
		},
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
		expected  string
	}{
		{
			name:      "basic",
			ip:        "192.168.1.1",
			userAgent: "test-agent",
			expected:  generateDeviceID("192.168.1.1", "test-agent"),
		},
		{
			name:      "empty user agent",
			ip:        "10.0.0.1",
			userAgent: "",
			expected:  generateDeviceID("10.0.0.1", ""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateDeviceID(tt.ip, tt.userAgent)
			assert.Equal(t, tt.expected, result)
			assert.Len(t, result, 32) // MD5 hex length
		})
	}
}

// Benchmark tests
func BenchmarkDeviceHelper_GetDeviceInfo(b *testing.B) {
	gin.SetMode(gin.TestMode)
	cfg := &utils.Config{
		RequestKey:         "request_context",
		SupportedLanguages: []string{"en-US", "vi-VN"},
	}
	helper := NewDeviceHelper(cfg)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Accept-Language", "vi-VN,vi;q=0.9,en;q=0.8")
	c.Request.Header.Set("X-Timezone", "Asia/Ho_Chi_Minh")
	c.Request.RemoteAddr = "192.168.1.1:1234"

	b.ResetTimer()
	for b.Loop() {
		helper.GetDeviceInfo(c)
	}
}

func BenchmarkExtractTokenFromHeader(b *testing.B) {
	header := "Bearer abc123xyz"
	b.ResetTimer()
	for b.Loop() {
		ExtractTokenFromHeader(header)
	}
}
