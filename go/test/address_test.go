package test

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/radiustechsystems/sdk/go/src/radius"
)

func TestNewAddress(t *testing.T) {
	t.Run("Creates address from bytes", func(t *testing.T) {
		hexStr := "0102030405060708090a0b0c0d0e0f1011121314"
		bytes, err := hex.DecodeString(hexStr)
		require.NoError(t, err, "Failed to decode hex string")

		addr := radius.NewAddress(bytes)
		expected := common.BytesToAddress(bytes)
		assert.Equal(t, expected, addr, "Address should match expected value")
		assert.Equal(t, expected.Hex(), addr.Hex(), "Checksum address should match expected")
	})

	t.Run("Handles different length byte arrays", func(t *testing.T) {
		shortBytes, err := hex.DecodeString("010203")
		require.NoError(t, err)

		shortAddr := radius.NewAddress(shortBytes)
		shortExpected := common.BytesToAddress(shortBytes)
		assert.Equal(t, shortExpected.Hex(), shortAddr.Hex(), "Checksum address should match expected")

		longStr := "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"
		longBytes, err := hex.DecodeString(longStr)
		require.NoError(t, err)

		longAddr := radius.NewAddress(longBytes)
		longExpected := common.BytesToAddress(longBytes)
		assert.Equal(t, longExpected.Hex(), longAddr.Hex(), "Checksum address should match expected")
	})
}

func TestNewAddressFromHex(t *testing.T) {
	t.Run("Creates address from hex string with 0x prefix", func(t *testing.T) {
		hexStr := "0xdAC17F958D2ee523a2206206994597C13D831ec7"
		addr := radius.NewAddressFromHex(hexStr)
		expected := common.HexToAddress(hexStr)
		assert.Equal(t, expected, addr, "Address should match expected value")
		assert.Equal(t, expected.Hex(), addr.Hex(), "Checksum address should match expected")
	})

	t.Run("Creates address from hex string without 0x prefix", func(t *testing.T) {
		hexStr := "A0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"
		addr := radius.NewAddressFromHex(hexStr)
		expected := common.HexToAddress(hexStr)
		assert.Equal(t, expected.Hex(), addr.Hex(), "Checksum address should match expected")
	})

	t.Run("Handles invalid hex strings", func(t *testing.T) {
		shortHex := "0x123"
		shortAddr := radius.NewAddressFromHex(shortHex)
		shortExpected := common.HexToAddress(shortHex)
		assert.Equal(t, shortExpected.Hex(), shortAddr.Hex(), "Checksum address should match expected")

		longHex := "0x" + "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
		longAddr := radius.NewAddressFromHex(longHex)
		longExpected := common.HexToAddress(longHex)
		assert.Equal(t, longExpected.Hex(), longAddr.Hex(), "Checksum address should match expected")

		invalidHex := "0xXYZ123"
		invalidAddr := radius.NewAddressFromHex(invalidHex)
		invalidExpected := common.HexToAddress(invalidHex)
		assert.Equal(t, invalidExpected.Hex(), invalidAddr.Hex(), "Checksum address should match expected")
	})

	t.Run("Handles empty input", func(t *testing.T) {
		emptyAddr := radius.NewAddressFromHex("")
		emptyExpected := common.HexToAddress("")
		assert.Equal(t, emptyExpected.Hex(), emptyAddr.Hex(), "Checksum address should match expected")
	})

	t.Run("Confirms EIP-55 checksum behavior", func(t *testing.T) {
		checksumAddr := "0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed"
		addr := radius.NewAddressFromHex(checksumAddr)
		assert.Equal(t, checksumAddr, addr.Hex(), "Should preserve EIP-55 checksum capitalization")

		lowercaseAddr := strings.ToLower(checksumAddr)
		addr = radius.NewAddressFromHex(lowercaseAddr)
		assert.Equal(t, checksumAddr, addr.Hex(), "Should convert to EIP-55 checksum format")
	})
}
