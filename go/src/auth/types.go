// Package auth provides interfaces and implementations for signing transactions and messages.
// It includes multiple signer implementations for different security requirements and key management strategies.
package auth

import (
	"context"
	"math/big"
	"net/http"

	"github.com/radiustechsystems/sdk/go/src/common"
)

// Signer is an interface for cryptographically signing messages and transactions.
// Different implementations provide different mechanisms for accessing private keys.
type Signer interface {
	// Address returns the Radius account address associated with the Signer
	// @return The address of the signer
	Address() common.Address

	// ChainID returns the Chain ID associated with the Signer
	// @return The chain ID
	ChainID() *big.Int

	// Hash returns the hash of the given transaction
	// @param tx The transaction to hash
	// @return The transaction hash
	Hash(tx *common.Transaction) common.Hash

	// SignMessage signs the given message using the EIP-191 standard
	// @param msg The message bytes to sign
	// @return The signature bytes, or an error if signing fails
	SignMessage(msg []byte) ([]byte, error)

	// SignTransaction signs the given transaction using the EIP-155 standard
	// @param tx The transaction to sign
	// @return The signed transaction, or an error if signing fails
	SignTransaction(tx *common.Transaction) (*common.SignedTransaction, error)
}

// SignerClient is an interface for the Radius Client methods that may be required by the Signer.
// This interface is implemented by the main Radius Client.
type SignerClient interface {
	// ChainID returns the Radius chain ID, which is used to sign transactions
	// @return The chain ID, or an error if it cannot be retrieved
	ChainID(ctx context.Context) (*big.Int, error)

	// HTTPClient returns the HTTP client used by the client to make requests
	// @return The HTTP client
	HTTPClient() *http.Client
}
