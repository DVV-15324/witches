package template

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestServiceNameProcessing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"lowercase", "user", "User"},
		{"uppercase", "USER", "User"},
		{"mixed case", "UsEr", "User"},
		{"with spaces", " user ", "User"},
		{"multiple words", "user profile", "Userprofile"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test logic của AddGoService
			serviceName := strings.TrimSpace(tt.input)
			serviceName = strings.ToLower(serviceName)
			serviceName = strings.ReplaceAll(serviceName, " ", "")
			serviceNameCap := cases.Title(language.English).String(serviceName)
			serviceNameCap = strings.ReplaceAll(serviceNameCap, " ", "")

			assert.Equal(t, tt.expected, serviceNameCap)
		})
	}
}

func TestServiceConfig(t *testing.T) {
	config := ServiceConfig{
		NameCap:    "User",
		Name:       "user",
		FolderName: "user-service",
		ModuleName: "github.com/example/project",
	}

	assert.Equal(t, "github.com/example/project", config.GetMuduleName())
	assert.Equal(t, "user", config.Name)
	assert.Equal(t, "User", config.NameCap)
}
