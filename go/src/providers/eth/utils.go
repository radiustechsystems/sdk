// Package eth provides an abstraction layer for Ethereum-compatible libraries in the Radius SDK.
// While the Radius SDK has its own concrete implementations of core data structures like addresses
// and transactions, this package maps those structures to Ethereum types to leverage the functionality
// provided by Ethereum libraries. This approach allows the SDK to benefit from well-tested Ethereum
// libraries while maintaining its own independent implementations that are optimized for Radius's
// custom database architecture for parallel transaction processing.
package eth

import (
	"context"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// BytesToAddress converts a byte slice to an Ethereum address.
//
// @param b Byte slice representing the address
// @return Address instance created from bytes
func BytesToAddress(b []byte) Address {
	return common.BytesToAddress(b)
}

// CreateAddress deterministically computes a contract address from a deployer address and nonce.
//
// @param from Address of the contract deployer
// @param nonce Nonce of the transaction creating the contract
// @return The computed contract address
func CreateAddress(from Address, nonce uint64) Address {
	return crypto.CreateAddress(from, nonce)
}

// NewAddress creates an address from a hex string.
//
// @param s Hex string representation of the address (with or without 0x prefix)
// @return Address instance created from the hex string
func NewAddress(s string) Address {
	return common.HexToAddress(s)
}

// NewClient creates a new Ethereum client connected to the specified URL.
//
// @param url URL of the Ethereum node (e.g. "http://localhost:8545")
// @param httpClient HTTP client to use for the connection
// @return Client instance and nil error on success
// @return nil and error if connection fails
func NewClient(url string, httpClient *http.Client) (*Client, error) {
	rpcClient, err := rpc.DialOptions(context.Background(), url, rpc.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}
	return ethclient.NewClient(rpcClient), nil
}

// NewRPCClient creates a new JSON-RPC client connected to the specified URL.
//
// @param url URL of the Ethereum node (e.g. "http://localhost:8545")
// @param httpClient HTTP client to use for the connection
// @return RPCClient instance and nil error on success
// @return nil and error if connection fails
func NewRPCClient(url string, httpClient *http.Client) (*rpc.Client, error) {
	return rpc.DialOptions(context.Background(), url, rpc.WithHTTPClient(httpClient))
}

// NewEIP155Signer creates a new signer for a specific chain ID.
//
// @param chainID Chain ID to use for the signer
// @return A new signer instance
func NewEIP155Signer(chainID *big.Int) EIP155Signer {
	return types.NewEIP155Signer(chainID)
}

// NewTx creates a new transaction with the given transaction data.
//
// @param inner Transaction data containing fields like recipient, value, etc.
// @return A new transaction instance
func NewTx(inner TxData) *Transaction {
	return types.NewTx(inner)
}

// NewABI creates a new ABI instance from a JSON string.
//
// @param abiStr JSON string representation of the ABI
// @return ABI instance and nil error on success
// @return Empty ABI and error if parsing fails
func NewABI(abiStr string) (ABI, error) {
	return abi.JSON(strings.NewReader(abiStr))
}

// Sender extracts the sender address from a signed transaction.
//
// @param signer Signer to use for extracting the address
// @param tx Signed transaction to extract sender from
// @return Sender address and nil error on success
// @return Zero address and error if extraction fails
func Sender(signer Signer, tx *Transaction) (Address, error) {
	from, err := types.Sender(signer, tx)
	if err != nil {
		return Address{}, err
	}
	return BytesToAddress(from.Bytes()), nil
}

// WaitMined waits for a transaction to be mined on Ethereum.
//
// @param ctx Context for the request (can be used for timeout)
// @param b Backend to use for checking transaction status
// @param tx Transaction to wait for
// @return Transaction receipt and nil error on success
// @return nil and error if waiting fails or times out
func WaitMined(ctx context.Context, b DeployBackend, tx *Transaction) (*Receipt, error) {
	return bind.WaitMined(ctx, b, tx)
}
