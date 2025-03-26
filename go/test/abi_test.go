package test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/radiustechsystems/sdk/go/src/radius"
)

func TestABI(t *testing.T) {
	t.Run("NewABI accepts a JSON string", func(t *testing.T) {
		abi, err := radius.NewABI(SimpleStorageABI)
		require.NoError(t, err, "Failed to create ABI")
		require.NotNil(t, abi, "ABI should not be nil")
	})

	t.Run("Pack method encodes arguments", func(t *testing.T) {
		abi, err := radius.NewABI(SimpleStorageABI)
		require.NoError(t, err, "Failed to create ABI")
		require.NotNil(t, abi, "ABI should not be nil")

		// Calculate function selector (first 4 bytes of the hash of the function signature)
		selector := crypto.Keccak256([]byte("set(uint256)"))[:4]

		// Create padded uint256 value (32 bytes, right-aligned)
		value := big.NewInt(42)
		paddedValue := make([]byte, 32)
		value.FillBytes(paddedValue)

		// Construct expected signature + encoded parameter
		expected := append(selector, paddedValue...)

		packed, err := abi.Pack("set", value)
		require.NoError(t, err, "Failed to pack arguments")
		require.NotNil(t, packed, "Packed arguments should not be nil")
		assert.Equal(t, expected, packed, "Unexpected packed arguments")
	})

	t.Run("Unpack method decodes return values", func(t *testing.T) {
		abi, err := radius.NewABI(SimpleStorageABI)
		require.NoError(t, err, "Failed to create ABI")
		require.NotNil(t, abi, "ABI should not be nil")

		// Return values are encoded as 32-byte chunks
		value := big.NewInt(42)
		paddedValue := make([]byte, 32)
		value.FillBytes(paddedValue)

		// Unpack should return an array of interface{} values
		expected := []interface{}{value}

		unpacked, err := abi.Unpack("get", paddedValue)
		require.NoError(t, err, "Failed to unpack return values")
		require.NotNil(t, unpacked, "Unpacked return values should not be nil")
		assert.Equal(t, expected, unpacked, "Unexpected unpacked return values")
	})
}
