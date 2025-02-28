// Package accounts provides functionality for managing Radius accounts.
// It includes support for creating accounts with different key management strategies,
// querying account balances and nonces, and signing transactions and messages.
package accounts

import (
	"context"
	"fmt"
	"math/big"

	"github.com/radiustechsystems/sdk/go/src/auth"
	"github.com/radiustechsystems/sdk/go/src/common"
)

// Account represents a Radius account that can be used to sign transactions.
// This struct provides methods for checking balance, retrieving nonce, and
// signing messages and transactions.
type Account struct {
	// Signer used to cryptographically sign messages and transactions
	Signer auth.Signer
}

// New creates a new Account with the given Option(s).
//
// @param opts Functional options to configure the account (e.g., WithSigner)
// @return A new Account instance configured with the provided options
func New(opts ...Option) *Account {
	a := &Account{}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// Address returns the address of the account.
//
// @return The account address, or zero address if no signer is available
func (a *Account) Address() common.Address {
	if a.Signer != nil {
		return a.Signer.Address()
	}
	return common.ZeroAddress()
}

// Balance returns the balance of the account in wei.
//
// @param ctx Context for the request
// @param client Radius client instance used to query the balance
// @return The account balance in wei and nil error on success
// @return nil and error if the balance cannot be retrieved from the network
func (a *Account) Balance(ctx context.Context, client AccountClient) (*big.Int, error) {
	return client.BalanceAt(ctx, a.Address())
}

// Nonce returns the next nonce (transaction count) of the account.
//
// @param ctx Context for the request
// @param client Radius client instance used to query the nonce
// @return The next nonce to use for transactions and nil error on success
// @return 0 and error if the nonce cannot be retrieved from the network
func (a *Account) Nonce(ctx context.Context, client AccountClient) (uint64, error) {
	return client.PendingNonceAt(ctx, a.Address())
}

// Send sends native currency to a recipient address.
//
// @param ctx Context for the request
// @param client Radius client instance used to send the transaction
// @param recipient Destination address to receive the funds
// @param amount Amount of native currency to send in wei
// @return Receipt of the completed transaction and nil error on success
// @return nil and error if no signer is available
// @return nil and error if the transaction fails
func (a *Account) Send(ctx context.Context, client AccountClient, recipient common.Address, amount *big.Int) (*common.Receipt, error) {
	if a.Signer == nil {
		return nil, fmt.Errorf("signer is required for sending transactions")
	}
	return client.Send(ctx, a.Signer, recipient, amount)
}

// SignMessage signs a message using the EIP-191 standard.
//
// @param msg Message bytes to sign
// @return The signature bytes and nil error on success
// @return nil and error if no signer is available
// @return nil and error if signing fails
func (a *Account) SignMessage(msg []byte) ([]byte, error) {
	if a.Signer == nil {
		return nil, fmt.Errorf("signer is required for signing messages")
	}

	signature, err := a.Signer.SignMessage(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to sign message: %w", err)
	}

	return signature, nil
}

// SignTransaction signs a transaction using the EIP-155 standard.
//
// @param tx Transaction to sign
// @return The signed transaction ready to be sent to the network and nil error on success
// @return nil and error if no signer is available
// @return nil and error if signing fails
func (a *Account) SignTransaction(tx *common.Transaction) (*common.SignedTransaction, error) {
	if a.Signer == nil {
		return nil, fmt.Errorf("signer is required for signing transactions")
	}

	signedTx, err := a.Signer.SignTransaction(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to create signed transaction: %w", err)
	}

	return signedTx, nil
}
