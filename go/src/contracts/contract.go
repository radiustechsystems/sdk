// Package contracts provides functionality for interacting with smart contracts on the Radius platform.
// It includes tools for creating contract instances, encoding/decoding method calls, and executing
// transactions on Radius contracts.
package contracts

import (
	"context"

	"github.com/radiustechsystems/sdk/go/src/auth"
	"github.com/radiustechsystems/sdk/go/src/common"
)

// Contract represents an EVM smart contract on the Radius platform.
// It provides methods to call read-only methods and execute state-changing methods,
// handling the ABI encoding/decoding automatically.
type Contract struct {
	// ABI is the contract's Application Binary Interface
	// Used for encoding and decoding method calls and return values
	ABI *common.ABI

	// address is the contract's address on Radius
	address common.Address
}

// New creates a new Contract with the given ABI and address.
//
// @param address The contract's address on Radius
// @param abi The contract's ABI (Application Binary Interface)
// @return A new Contract instance
func New(address common.Address, abi *common.ABI) *Contract {
	return &Contract{
		ABI:     abi,
		address: address,
	}
}

// Address returns the address of the contract.
//
// @return The contract's address on Radius
func (c *Contract) Address() common.Address {
	return c.address
}

// Call executes a contract method call and returns the decoded result. This is used for read-only contract methods,
// and does not require a transaction to be sent to Radius.
//
// @param ctx Context for the request
// @param client Radius client instance used to make the call
// @param method Name of the method to call on the contract
// @param args Arguments to pass to the contract method
// @return Array of decoded return values from the contract method and nil error on success
// @return nil and error if the contract ABI is missing
// @return nil and error if the contract address is missing or zero
// @return nil and error if the contract method call fails
func (c *Contract) Call(ctx context.Context, client ContractClient, method string, args ...interface{}) ([]interface{}, error) {
	return client.Call(ctx, c, method, args...)
}

// Execute executes a contract method call and returns the transaction receipt. This is used for state-changing contract
// methods, and requires a transaction to be sent to Radius.
//
// @param ctx Context for the request
// @param client Radius client instance used to execute the transaction
// @param signer The signer used to sign the transaction
// @param method Name of the method to execute on the contract
// @param args Arguments to pass to the contract method
// @return Transaction receipt after the method execution and nil error on success
// @return nil and error if the contract ABI is missing
// @return nil and error if the contract address is missing or zero
// @return nil and error if the transaction fails or is reverted
// @return nil and error if the transaction receipt is not returned
func (c *Contract) Execute(ctx context.Context, client ContractClient, signer auth.Signer, method string, args ...interface{}) (*common.Receipt, error) {
	return client.Execute(ctx, c, signer, method, args...)
}
