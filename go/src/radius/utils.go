package radius

import (
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
)

func BytecodeFromHex(s string) []byte {
	// Remove 0x prefix if present
	if len(s) >= 2 && s[0:2] == "0x" {
		s = s[2:]
	}

	bytecode, err := hex.DecodeString(s)
	if err != nil {
		return nil
	}
	return bytecode
}

func GeneratePrivateKey() *ecdsa.PrivateKey {
	key, err := crypto.GenerateKey()
	if err != nil {
		return nil
	}
	return key
}

func GeneratePrivateKeyFromHex(h string) *ecdsa.PrivateKey {
	key, err := crypto.HexToECDSA(h)
	if err != nil {
		return nil
	}
	return key
}

func PadBytes(data []byte, size int) []byte {
	padded := make([]byte, size)

	if len(data) > size {
		copy(padded, data[len(data)-size:])
	} else {
		copy(padded[size-len(data):], data)
	}

	return padded
}
