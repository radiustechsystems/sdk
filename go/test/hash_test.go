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

func TestNewHash(t *testing.T) {
	t.Run("Creates hash from bytes", func(t *testing.T) {
		hexStr := "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"
		bytes, err := hex.DecodeString(hexStr)
		require.NoError(t, err, "Failed to decode hex string")

		hash := radius.NewHash(bytes)
		expected := common.BytesToHash(bytes)
		assert.Equal(t, expected, hash, "Hash should match expected value")
		assert.Equal(t, expected.Hex(), hash.Hex(), "Hex representation should match expected")
	})

	t.Run("Handles different length byte arrays", func(t *testing.T) {
		shortBytes, err := hex.DecodeString("010203")
		require.NoError(t, err)

		shortHash := radius.NewHash(shortBytes)
		shortExpected := common.BytesToHash(shortBytes)
		assert.Equal(t, shortExpected.Hex(), shortHash.Hex(), "Hex representation should match expected")

		longStr := "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f2021222324252627"
		longBytes, err := hex.DecodeString(longStr)
		require.NoError(t, err)

		longHash := radius.NewHash(longBytes)
		longExpected := common.BytesToHash(longBytes)
		assert.Equal(t, longExpected.Hex(), longHash.Hex(), "Hex representation should match expected")
	})
}

func TestNewHashFromHex(t *testing.T) {
	t.Run("Creates hash from hex string with 0x prefix", func(t *testing.T) {
		hexStr := "0xd4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3"
		hash := radius.NewHashFromHex(hexStr)
		expected := common.HexToHash(hexStr)
		assert.Equal(t, expected, hash, "Hash should match expected value")
		assert.Equal(t, expected.Hex(), hash.Hex(), "Hex representation should match expected")
	})

	t.Run("Creates hash from hex string without 0x prefix", func(t *testing.T) {
		hexStr := "d4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3"
		hash := radius.NewHashFromHex(hexStr)
		expected := common.HexToHash(hexStr)
		assert.Equal(t, expected.Hex(), hash.Hex(), "Hex representation should match expected")
		assert.Equal(t, "0x"+hexStr, hash.Hex(), "Should add 0x prefix if missing")
	})

	t.Run("Handles mixed case hex strings", func(t *testing.T) {
		mixedCaseHex := "0xD4e56740F876aEf8c010b86A40d5f56745A118d0906a34E69aEc8C0db1CB8fa3"
		hash := radius.NewHashFromHex(mixedCaseHex)
		expected := common.HexToHash(mixedCaseHex)
		assert.Equal(t, expected.Hex(), hash.Hex(), "Hex representation should match expected")
		assert.Equal(t, strings.ToLower(mixedCaseHex), hash.Hex(), "Hash should be lowercase")
	})

	t.Run("Handles invalid hex strings", func(t *testing.T) {
		shortHex := "0x123"
		shortHash := radius.NewHashFromHex(shortHex)
		shortExpected := common.HexToHash(shortHex)
		assert.Equal(t, shortExpected.Hex(), shortHash.Hex(), "Hex representation should match expected")

		longHex := "0x" + strings.Repeat("1234567890abcdef", 4) + "extra"
		longHash := radius.NewHashFromHex(longHex)
		longExpected := common.HexToHash(longHex)
		assert.Equal(t, longExpected.Hex(), longHash.Hex(), "Hex representation should match expected")

		invalidHex := "0xXYZ123"
		invalidHash := radius.NewHashFromHex(invalidHex)
		invalidExpected := common.HexToHash(invalidHex)
		assert.Equal(t, invalidExpected.Hex(), invalidHash.Hex(), "Hex representation should match expected")
	})

	t.Run("Handles empty input", func(t *testing.T) {
		emptyHash := radius.NewHashFromHex("")
		emptyExpected := common.HexToHash("")
		assert.Equal(t, emptyExpected.Hex(), emptyHash.Hex(), "Hex representation should match expected")
		assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000000",
			emptyHash.Hex(), "Empty hash should be all zeros")
	})
}
