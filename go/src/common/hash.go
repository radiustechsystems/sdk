package common

import (
	"encoding/hex"
)

// Hash represents a 32-byte Keccak-256 hash used for transactions, blocks, and states.
// This struct provides methods to access the hash in different formats.
type Hash struct {
	// bytes is the internal byte representation of the hash
	bytes []byte
}

// NewHash creates a new Hash with the given bytes
// @param bytes The byte representation of the hash
// @return A new Hash instance
func NewHash(bytes []byte) Hash {
	return Hash{bytes: bytes}
}

// Bytes returns the bytes of the Hash
// @return The byte representation of the hash
func (h *Hash) Bytes() []byte {
	return h.bytes
}

// Hex returns the hexadecimal string of the Hash with 0x prefix
// @return The hexadecimal string representation of the hash with 0x prefix
func (h *Hash) Hex() string {
	return "0x" + hex.EncodeToString(h.bytes)
}

// HexWithoutPrefix returns the hexadecimal string of the Hash without 0x prefix
// @return The hexadecimal string representation of the hash without 0x prefix
func (h *Hash) HexWithoutPrefix() string {
	return hex.EncodeToString(h.bytes)
}
