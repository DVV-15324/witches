package utils

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRandomStr(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		wantLen int
	}{
		{"length 10", 10, 20},
		{"length 5", 5, 10},
		{"length 1", 1, 2},
		{"length 0", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str, err := RandomStr(tt.length)
			assert.NoError(t, err)
			if tt.length > 0 {
				assert.Len(t, str, tt.wantLen)
				_, err := hex.DecodeString(str)
				assert.NoError(t, err)
				str2, err := RandomStr(tt.length)
				assert.NoError(t, err)
				assert.NotEqual(t, str, str2)
			} else {
				assert.Empty(t, str)
			}
		})
	}
}
func TestHash_GenerateFromPassword(t *testing.T) {
	h := &Hash{}
	salt := "salt123"
	password := "mypassword"
	hashed, err := h.GenerateFromPassword(password, salt)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashed)
	hashed2, err := h.GenerateFromPassword(password, salt)
	assert.NoError(t, err)
	assert.NotEqual(t, hashed, hashed2)
}

func TestHash_CompareHashAndPassword(t *testing.T) {
	h := &Hash{}
	salt := "salt123"
	password := "mypassword"

	// Generate hash
	hashed, err := h.GenerateFromPassword(password, salt)
	require.NoError(t, err)

	tests := []struct {
		name     string
		hashed   string
		password string
		salt     string
		expected bool
	}{
		{
			name:     "correct password",
			hashed:   hashed,
			password: password,
			salt:     salt,
			expected: true,
		},
		{
			name:     "wrong password",
			hashed:   hashed,
			password: "wrongpassword",
			salt:     salt,
			expected: false,
		},
		{
			name:     "wrong salt",
			hashed:   hashed,
			password: password,
			salt:     "wrongsalt",
			expected: false,
		},
		{
			name:     "empty password",
			hashed:   hashed,
			password: "",
			salt:     salt,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := h.CompareHashAndPassword(tt.hashed, tt.password, tt.salt)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHash_GenerateFromPassword_Error(t *testing.T) {
	h := &Hash{}
	longPassword := string(make([]byte, 10000))
	_, err := h.GenerateFromPassword(longPassword, "salt")
	t.Logf("Long password test: %v", err)
}
