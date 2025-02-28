// Package crypto provides cryptographic utilities for use with the Radius platform.
// It includes functions for working with ECDSA keys, Keccak256 hashing,
// address derivation, and transaction signing.
package crypto

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/radiustechsystems/sdk/go/src/common"
)

// HexToECDSA converts a hexadecimal string to an ECDSA private key.
// The input string should be a hex-encoded string of the private key (with or without 0x prefix).
//
// @param key The hex string representation of the private key
// @return The ECDSA private key and nil error on success
// @return nil and error if the key is invalid or cannot be parsed
func HexToECDSA(key string) (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(key)
}

// Keccak256 calculates the Keccak256 hash of the input data.
// This is the hashing algorithm used by Ethereum for various cryptographic operations.
// Multiple byte slices will be concatenated before hashing.
//
// @param data One or more byte slices to hash
// @return The 32-byte Keccak256 hash of the input data
func Keccak256(data ...[]byte) []byte {
	return crypto.Keccak256(data...)
}

// PubkeyToAddress converts an ECDSA public key to a Radius address.
// The address is derived by taking the Keccak256 hash of the public key
// and keeping the last 20 bytes (same algorithm as Ethereum).
//
// @param p The ECDSA public key to convert
// @return The corresponding Radius address
func PubkeyToAddress(p ecdsa.PublicKey) common.Address {
	return common.NewAddress(crypto.PubkeyToAddress(p).Bytes())
}

// Sign creates a cryptographic signature of a digest hash using an ECDSA private key.
// The signature is in the Ethereum format: [R || S || V] where V is 0 or 1.
//
// @param digestHash The 32-byte hash to sign (typically a Keccak256 hash)
// @param prv The ECDSA private key to sign with
// @return The signature bytes and nil error on success
// @return nil and error if signing fails
func Sign(digestHash []byte, prv *ecdsa.PrivateKey) (sig []byte, err error) {
	return crypto.Sign(digestHash, prv)
}
