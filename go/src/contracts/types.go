// Package contracts provides functionality for interacting with smart contracts on the Radius platform.
// It includes tools for creating contract instances, encoding/decoding method calls, and executing
// transactions on Radius contracts.
package contracts

import (
	"context"

	"github.com/radiustechsystems/sdk/go/src/auth"
	"github.com/radiustechsystems/sdk/go/src/common"
)

// ContractClient is an interface for interacting with EVM smart contracts via a Radius Client.
// This interface is implemented by the main Radius Client.
type ContractClient interface {
	// Call executes a contract method call and returns the decoded result. This is used for read-only contract methods,
	// and does not require a transaction to be sent to Radius.
	//
	// @param ctx Context for the request
	// @param contract Contract instance to interact with
	// @param method Name of the method to call on the contract
	// @param args Arguments to pass to the contract method
	// @return Array of decoded return values from the contract method and nil error on success
	// @return nil and error if the contract ABI is missing
	// @return nil and error if the contract address is missing or zero
	// @return nil and error if the contract method call fails
	Call(ctx context.Context, contract *Contract, method string, args ...interface{}) ([]interface{}, error)

	// Execute executes a contract method that modifies Radius state. This is used for write operations, and
	// requires a transaction to be sent to Radius.
	//
	// @param ctx Context for the request
	// @param contract Contract instance to interact with
	// @param signer The signer used to sign the transaction
	// @param method Name of the method to execute on the contract
	// @param args Arguments to pass to the contract method
	// @return Transaction receipt after the method execution and nil error on success
	// @return nil and error if the contract ABI is missing
	// @return nil and error if the contract address is missing or zero
	// @return nil and error if the transaction fails or is reverted
	// @return nil and error if the transaction receipt is not returned
	Execute(ctx context.Context, contract *Contract, signer auth.Signer, method string, args ...interface{}) (*common.Receipt, error)
}
