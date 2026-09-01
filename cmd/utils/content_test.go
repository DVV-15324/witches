package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateContent(t *testing.T) {
	content := CreateContent()
	assert.NotEmpty(t, content)
	assert.Contains(t, content, "APP_PORT=8080")
	assert.Contains(t, content, "DB_DRIVER=your_driver")
	assert.Contains(t, content, "REDIS_HOST=localhost")
	assert.Contains(t, content, "UID_BITS=26")
	assert.Contains(t, content, "REQUEST_KEY=request_context")
	assert.Contains(t, content, "LOG_PATH=./logs")
	assert.Contains(t, content, "CORS_ALLOW_ORIGINS=*")
	assert.Contains(t, content, "SUPPORTED_LANGUAGES=en-US,vi-VN")
}
