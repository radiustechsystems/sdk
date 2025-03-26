package test

import (
	"crypto/ecdsa"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"

	"github.com/radiustechsystems/sdk/go/src/radius"
)

func TestBytecodeFromHex(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []byte
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: []byte{},
		},
		{
			name:     "Even-length string",
			input:    "abcdef",
			expected: []byte{0xab, 0xcd, 0xef},
		},
		{
			name:     "Odd-length string",
			input:    "123",
			expected: nil,
		},
		{
			name:     "Mixed case",
			input:    "0xABCdef",
			expected: []byte{0xab, 0xcd, 0xef},
		},
		{
			name:     "Invalid hex",
			input:    "0xZZ",
			expected: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := radius.BytecodeFromHex(tc.input)
			assert.Equal(t, tc.expected, result, "Result should match expected output")
		})
	}
}

func TestGeneratePrivateKey(t *testing.T) {
	key := radius.GeneratePrivateKey()
	t.Logf("Generated private key: %s", key)
	assert.NotNil(t, key, "GeneratePrivateKey should return a non-nil private key")

	pubKey := key.Public()
	assert.NotNil(t, pubKey, "Generated private key should have a public key")

	_, ok := pubKey.(*ecdsa.PublicKey)
	assert.True(t, ok, "Public key should be of type *ecdsa.PublicKey")

	key2 := radius.GeneratePrivateKey()
	assert.NotNil(t, key2, "Second call to GeneratePrivateKey should return a non-nil private key")

	keyBytes := crypto.FromECDSA(key)
	key2Bytes := crypto.FromECDSA(key2)
	assert.Equal(t, 32, len(keyBytes), "Private key bytes should be 32 bytes long")
	assert.Equal(t, 32, len(key2Bytes), "Private key bytes should be 32 bytes long")
	assert.NotEqual(t, keyBytes, key2Bytes, "Two generated keys should not be equal")
}

func TestPadBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		size     int
		expected []byte
	}{
		{
			name:     "Pad empty array to size 5",
			input:    []byte{},
			size:     5,
			expected: []byte{0, 0, 0, 0, 0},
		},
		{
			name:     "Pad smaller array to larger size",
			input:    []byte{1, 2, 3},
			size:     5,
			expected: []byte{0, 0, 1, 2, 3},
		},
		{
			name:     "Pad to same size (no padding needed)",
			input:    []byte{1, 2, 3},
			size:     3,
			expected: []byte{1, 2, 3},
		},
		{
			name:     "Pad to smaller size (truncation)",
			input:    []byte{1, 2, 3, 4, 5},
			size:     3,
			expected: []byte{3, 4, 5}, // Only last 3 bytes are retained
		},
		{
			name:     "Pad to size 0",
			input:    []byte{1, 2, 3},
			size:     0,
			expected: []byte{},
		},
		{
			name:     "Pad binary data",
			input:    []byte{0xFF, 0xAA, 0x55},
			size:     5,
			expected: []byte{0, 0, 0xFF, 0xAA, 0x55},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := radius.PadBytes(tc.input, tc.size)
			assert.Equal(t, tc.size, len(result), "Result size should match expected size")
			assert.Equal(t, tc.expected, result, "Result should match expected output")

			// Ensure the original data wasn't modified
			if len(tc.input) > 0 && len(result) > 0 && &tc.input[0] == &result[0] {
				t.Error("PadBytes modified the original data instead of creating a new slice")
			}
		})
	}
}

func TestPadBytes_EdgeCases(t *testing.T) {
	// Test padding a byte slice to a larger size
	largeSize := 1024
	input := []byte{1, 2, 3}
	result := radius.PadBytes(input, largeSize)
	assert.Equal(t, largeSize, len(result), "Result size should match expected size")

	// Check that the padding is correct
	for i := 0; i < largeSize-len(input); i++ {
		assert.Equal(t, byte(0), result[i], "Expected padding byte at index %d to be 0", i)
	}

	// Check that the original data is preserved at the end
	for i := 0; i < len(input); i++ {
		assert.Equal(t, input[i], result[largeSize-len(input)+i], "Expected byte at index %d to be %d", largeSize-len(input)+i, input[i])
	}

	// Test padding a nil slice
	var nilInput []byte
	nilResult := radius.PadBytes(nilInput, 5)
	assert.Equal(t, 5, len(nilResult), "Expected padded nil slice to have length 5")

	// Check that the padding is correct
	for i := 0; i < 5; i++ {
		assert.Equal(t, byte(0), nilResult[i], "Expected padding byte at index %d to be 0", i)
	}
}
