// Package eth provides an abstraction layer for Ethereum-compatible libraries in the Radius SDK.
// While the Radius SDK has its own concrete implementations of core data structures like addresses
// and transactions, this package maps those structures to Ethereum types to leverage the functionality
// provided by Ethereum libraries. This approach allows the SDK to benefit from well-tested Ethereum
// libraries while maintaining its own independent implementations that are optimized for Radius's
// custom database architecture for parallel transaction processing.
package eth

import (
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// This file contains type aliases that map Radius SDK types to Ethereum library types.
// The SDK has its own concrete implementations of these structures, and we use these
// aliases only to leverage Ethereum library functionality when needed.
type (
	// ABI represents a smart contract's Application Binary Interface.
	// Used for encoding and decoding interactions with smart contracts.
	ABI = abi.ABI

	// Address represents a 20-byte account or contract address in Radius.
	// Used to identify accounts and smart contracts in the Radius system.
	Address = common.Address

	// CallMsg contains parameters for contract method calls in Radius.
	// Used when calling read-only contract methods.
	CallMsg = ethereum.CallMsg

	// Client is a client for interacting with a Radius node.
	// Provides methods for querying state and sending transactions.
	Client = ethclient.Client

	// DeployBackend is an interface for deploying contracts to Radius.
	// Abstracts the backend used for contract deployment.
	DeployBackend = bind.DeployBackend

	// EIP155Signer implements standardized transaction signing for Radius.
	// Used to create signatures for transactions with replay protection.
	EIP155Signer = types.EIP155Signer

	// Log represents a smart contract event log in Radius.
	// Contains data emitted by contract events during transaction execution.
	Log = types.Log

	// LegacyTx is a transaction in the original format for Radius.
	// Used for compatibility with EVM transaction formats.
	LegacyTx = types.LegacyTx

	// Signer is an interface for producing signatures for Radius transactions.
	// Provides methods to sign transactions.
	Signer = types.Signer

	// Transaction represents a Radius transaction.
	// Contains all data needed to execute a state change in the Radius system.
	Transaction = types.Transaction

	// TxData is an interface for different transaction types in Radius.
	// Allows supporting multiple transaction formats.
	TxData = types.TxData

	// Receipt is a Radius transaction receipt.
	// Contains information about a completed transaction, including status and logs.
	Receipt = types.Receipt

	// RPCClient is a client for making JSON-RPC calls to Radius.
	// Used for low-level communication with Radius JSON-RPC endpoints.
	RPCClient = rpc.Client
)
