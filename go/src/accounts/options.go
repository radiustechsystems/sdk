// Package accounts provides functionality for managing Radius accounts.
// It includes support for creating accounts with different key management strategies,
// querying account balances and nonces, and signing transactions and messages.
package accounts

import (
	"crypto/ecdsa"

	"github.com/radiustechsystems/sdk/go/src/auth"
	"github.com/radiustechsystems/sdk/go/src/auth/privatekey"
	"github.com/radiustechsystems/sdk/go/src/crypto"
)

// Option is a functional option for configuring a new Account.
// Options allow for a flexible API to construct accounts with various configurations.
type Option func(*Account)

// WithPrivateKey creates an Account using a private key.
//
// @param key ECDSA private key to use for signing
// @param client AccountClient used for network operations
// @return An Option function that configures an Account with the provided private key
func WithPrivateKey(key *ecdsa.PrivateKey, client AccountClient) Option {
	return func(a *Account) {
		a.Signer = privatekey.New(key, client)
	}
}

// WithPrivateKeyHex creates an Account using a private key hex string.
// The private key will be stored in memory, so for production systems with high security
// requirements, consider using WithSigner instead, along with a hardware security module
// or key management service.
//
// @param key Private key as a hex string
// @param client AccountClient used for network operations
// @return An Option function that configures an Account with the provided private key
func WithPrivateKeyHex(key string, client AccountClient) Option {
	privateKey, _ := crypto.HexToECDSA(key)
	return WithPrivateKey(privateKey, client)
}

// WithSigner creates an Account using a custom Signer.
// This is useful when you want to use a custom signing implementation, such as a hardware
// security module or key management service.
//
// @param signer Signer used to sign transactions and messages
// @return An Option function that configures an Account with the provided signer
func WithSigner(signer auth.Signer) Option {
	return func(a *Account) {
		a.Signer = signer
	}
}
