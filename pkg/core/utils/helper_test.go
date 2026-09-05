package utils

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/DVV-15324/witches/cmd/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
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
				c.Request.Header.Set("User-Agent", "mock-agent")
				c.Request = c.Request.WithContext(ctx)
				return c
			},
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
				return c
			},
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
				return c
			},
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
				return c
			},
			expectedUserAgent: "test-ua3",
			expectedLocale:    "en-US",
			expectedTimezone:  "Asia/Ho_Chi_Minh",
		},
		{
			name: "with Accept-Language - vi without region",
			setupContext: func() *gin.Context {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				c.Request = httptest.NewRequest("GET", "/", nil)
				c.Request.Header.Set("Accept-Language", "vi")
				return c
			},
			expectedUserAgent: "",
			expectedLocale:    "vi-VN",
			expectedTimezone:  "UTC",
		},
		{
			name: "with Accept-Language - en without region",
			setupContext: func() *gin.Context {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				c.Request = httptest.NewRequest("GET", "/", nil)
				c.Request.Header.Set("Accept-Language", "en")
				return c
			},
			expectedUserAgent: "",
			expectedLocale:    "en-US",
			expectedTimezone:  "UTC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.setupContext()
			_, userAgent, locale, tz := helper.GetInfo(c)

			assert.Equal(t, tt.expectedUserAgent, userAgent, "UserAgent mismatch")
			assert.Equal(t, tt.expectedLocale, locale, "Locale mismatch")
			assert.Equal(t, tt.expectedTimezone, tz, "Timezone mismatch")
		})
	}
}
func TestHelper_GetInfo_RegionFallback(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		acceptLanguage string
		expectedLocale string
	}{
		{
			name:           "Vietnamese fallback region",
			acceptLanguage: "vi",
			expectedLocale: "vi-VN",
		},
		{
			name:           "Default fallback region",
			acceptLanguage: "en",
			expectedLocale: "en-US",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			helper := NewHelper(&utils.Config{
				SupportedLanguages: []string{"vi", "en"},
			})

			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest("GET", "/", nil)
			c.Request.Header.Set("Accept-Language", tt.acceptLanguage)

			_, _, locale, _ := helper.GetInfo(c)

			assert.Equal(t, tt.expectedLocale, locale)
		})
	}
}
func TestGetRegion(t *testing.T) {
	tests := []struct {
		name     string
		base     language.Base
		region   language.Region
		expected string
	}{
		{
			name:     "Vietnamese without region",
			base:     language.MustParseBase("vi"),
			region:   language.Region{},
			expected: "VN",
		},
		{
			name:     "Other language without region",
			base:     language.MustParseBase("en"),
			region:   language.Region{},
			expected: "US",
		},
		{
			name:     "Existing region",
			base:     language.MustParseBase("en"),
			region:   language.MustParseRegion("GB"),
			expected: "GB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getRegion(tt.base, tt.region)
			assert.Equal(t, tt.expected, got.String())
		})
	}
}
