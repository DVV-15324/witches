package utils

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/DVV-15324/witches/cmd/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewHelper(t *testing.T) {
	tests := []struct {
		name   string
		cfg    *utils.Config
		expect int
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
			helper := NewHelper(tt.cfg)
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

func TestHelper_GetInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &utils.Config{
		RequestKey:         "request_context",
		SupportedLanguages: []string{"en-US", "vi-VN"},
	}
	helper := NewHelper(cfg)

	tests := []struct {
		name              string
		setupContext      func() *gin.Context
		expectedJwk       string
		expectedUserAgent string
		expectedLocale    string
		expectedTimezone  string
	}{
		{
			name: "with request context",
			setupContext: func() *gin.Context {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				c.Request = httptest.NewRequest("GET", "/", nil)
				ctx := context.WithValue(c.Request.Context(), cfg.RequestKey, &RequestContext{})
				c.Request.Header.Set("DPoP", "DPoP JWK123")
				c.Request.Header.Set("User-Agent", "mock-agent")
				c.Request = c.Request.WithContext(ctx)
				return c
			},
			expectedJwk:       "DPoP JWK123",
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
				c.Request.Header.Set("DPoP", "DPoP JWK123")
				return c
			},
			expectedJwk:       "DPoP JWK123",
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
				c.Request.Header.Set("DPoP", "DPoP JWK123")
				return c
			},
			expectedJwk:       "DPoP JWK123",
			expectedUserAgent: "test-ua",
			expectedLocale:    "vi-VN",
			expectedTimezone:  "UTC",
		},
		{
			name: "with X-Timezone header",
			setupContext: func() *gin.Context {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				c.Request = httptest.NewRequest("GET", "/", nil)
				c.Request.Header.Set("X-Timezone", "Asia/Ho_Chi_Minh")
				c.Request.Header.Set("User-Agent", "test-ua3")
				c.Request.Header.Set("DPoP", "DPoP JWK123")
				return c
			},
			expectedJwk:       "DPoP JWK123",
			expectedUserAgent: "test-ua3",
			expectedLocale:    "en-US",
			expectedTimezone:  "Asia/Ho_Chi_Minh",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.setupContext()
			DPoPJwk, _, userAgent, locale, tz := helper.GetInfo(c)

			assert.Equal(t, tt.expectedJwk, DPoPJwk, "DPoPJwk mismatch")
			assert.Equal(t, tt.expectedUserAgent, userAgent, "UserAgent mismatch")
			assert.Equal(t, tt.expectedLocale, locale, "Locale mismatch")
			assert.Equal(t, tt.expectedTimezone, tz, "Timezone mismatch")
		})
	}
}

func TestExtractDPoPTokenFromHeader(t *testing.T) {
	tests := []struct {
		name        string
		header      string
		expected    string
		expectError bool
	}{
		{"valid token", "DPoP JWK123", "JWK123", false},
		{"empty header", "", "authorization header is required", true},
		{"missing DPoP prefix", "JWK123", "invalid authorization header format", true},
		{"too many parts", "invalid authorization header format", "", true},
		{"invalid prefix", "invalid authorization header format", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Helper{}
			result, err := h.ExtractDPoPTokenFromHeader(tt.header)
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
