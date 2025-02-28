// Package accounts provides functionality for managing Radius accounts.
// It includes support for creating accounts with different key management strategies,
// querying account balances and nonces, and signing transactions and messages.
package accounts

import (
	"context"
	"math/big"
	"net/http"

	"github.com/radiustechsystems/sdk/go/src/auth"
	"github.com/radiustechsystems/sdk/go/src/common"
)

// AccountClient is an interface for account operations.
// This interface is implemented by the main Radius Client.
type AccountClient interface {
	// BalanceAt returns the balance of an account in wei.
	//
	// @param ctx Context for the request
	// @param address Address to check the balance for
	// @return The account balance in wei and nil error on success
	// @return nil and error if the balance cannot be retrieved from the network
	BalanceAt(ctx context.Context, address common.Address) (*big.Int, error)

	// ChainID returns the Radius chain ID, which is used to sign transactions.
	//
	// @param ctx Context for the request
	// @return The chain ID of the connected network and nil error on success
	// @return nil and error if the chain ID cannot be retrieved
	ChainID(ctx context.Context) (*big.Int, error)

	// EstimateGas estimates the gas cost of a transaction with a safety margin.
	//
	// @param ctx Context for the request
	// @param tx Transaction to estimate gas for
	// @return The estimated gas cost in gas units and nil error on success
	// @return 0 and error if the gas estimation fails
	EstimateGas(ctx context.Context, tx *common.Transaction) (uint64, error)

	// HTTPClient returns the HTTP client used by the client to make requests.
	//
	// @return The HTTP client used for API requests
	HTTPClient() *http.Client

	// PendingNonceAt returns the next nonce (transaction count) for an account.
	//
	// @param ctx Context for the request
	// @param address Address to check the nonce for
	// @return The next nonce to use for transactions and nil error on success
	// @return 0 and error if the nonce cannot be retrieved from the network
	PendingNonceAt(ctx context.Context, address common.Address) (uint64, error)

	// Send sends native currency to a recipient address.
	//
	// @param ctx Context for the request
	// @param signer The signer used to sign the transaction
	// @param recipient Destination address to receive the funds
	// @param amount Amount of native currency to send in wei
	// @return Receipt of the completed transaction and nil error on success
	// @return nil and error if the transaction fails
	// @return nil and error if the transaction receipt is not returned
	Send(ctx context.Context, signer auth.Signer, recipient common.Address, amount *big.Int) (*common.Receipt, error)
}
