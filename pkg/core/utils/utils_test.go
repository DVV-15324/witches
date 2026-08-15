package utils

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUID_ToBase58AndDecode(t *testing.T) {
	tests := []struct {
		name     string
		localID  uint32
		objectID uint
		bits     int64
	}{
		{"bits 26", 123, 1, 26},
		{"bits 26 large local", 123456, 5, 26},
		{"bits 32", 999, 3, 32},
		{"bits 32 max object", 456, 100, 32},
		{"bits 20", 789, 7, 20},
		{"bits 10", 1, 1, 10},
		{"bits 10 local max", 1023, 1, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uid := NewUID(tt.localID, tt.objectID)
			encoded := uid.ToBase58(tt.bits)
			assert.NotEmpty(t, encoded)

			decoded := DecodeFromBase58(encoded, tt.bits)
			require.NotNil(t, decoded)
			assert.Equal(t, tt.localID, decoded.LocalID)
			assert.Equal(t, tt.objectID, decoded.ObjectID)
		})
	}
}

func TestUID_ToBase58_DifferentBits(t *testing.T) {
	uid := NewUID(100, 2)

	enc26 := uid.ToBase58(26)
	enc32 := uid.ToBase58(32)
	enc20 := uid.ToBase58(20)

	assert.NotEqual(t, enc26, enc32)
	assert.NotEqual(t, enc26, enc20)
	assert.NotEqual(t, enc32, enc20)

	dec26 := DecodeFromBase58(enc26, 26)
	dec32 := DecodeFromBase58(enc32, 32)
	dec20 := DecodeFromBase58(enc20, 20)

	assert.Equal(t, uid.LocalID, dec26.LocalID)
	assert.Equal(t, uid.ObjectID, dec26.ObjectID)
	assert.Equal(t, uid.LocalID, dec32.LocalID)
	assert.Equal(t, uid.ObjectID, dec32.ObjectID)
	assert.Equal(t, uid.LocalID, dec20.LocalID)
	assert.Equal(t, uid.ObjectID, dec20.ObjectID)
}

func TestUID_Decode_Invalid(t *testing.T) {
	decoded := DecodeFromBase58("", 26)
	assert.NotNil(t, decoded)
	assert.Equal(t, uint32(0), decoded.LocalID)
	assert.Equal(t, uint(0), decoded.ObjectID)

	decoded = DecodeFromBase58("invalid!!!", 26)
	assert.NotNil(t, decoded)
	assert.Equal(t, uint32(0), decoded.LocalID)
	assert.Equal(t, uint(0), decoded.ObjectID)
}

func BenchmarkUID_ToBase58(b *testing.B) {
	uid := NewUID(123456, 5)
	b.ResetTimer()
	for b.Loop() {
		uid.ToBase58(26)
	}

}

func BenchmarkUID_DecodeFromBase58(b *testing.B) {
	uid := NewUID(123456, 5)
	encoded := uid.ToBase58(26)
	b.ResetTimer()
	for b.Loop() {
		DecodeFromBase58(encoded, 26)
	}

}
