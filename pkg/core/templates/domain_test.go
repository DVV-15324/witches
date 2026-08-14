package template

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestDomainNameProcessing(t *testing.T) {
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
			// Test logic của AddGodomain
			domainName := strings.TrimSpace(tt.input)
			domainName = strings.ToLower(domainName)
			domainName = strings.ReplaceAll(domainName, " ", "")
			domainNameCap := cases.Title(language.English).String(domainName)
			domainNameCap = strings.ReplaceAll(domainNameCap, " ", "")

			assert.Equal(t, tt.expected, domainNameCap)
		})
	}
}

func TestDomainConfig(t *testing.T) {
	config := DomainConfig{
		NameCap:    "User",
		Name:       "user",
		FolderName: "user",
		ModuleName: "github.com/example/project",
	}

	assert.Equal(t, "github.com/example/project", config.GetMuduleName())
	assert.Equal(t, "user", config.Name)
	assert.Equal(t, "User", config.NameCap)
}
