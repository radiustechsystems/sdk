// Package clef provides a Signer implementation using the Clef external signing tool.
// This approach allows for more secure private key management by delegating signing operations
// to an external process that can implement additional security measures.
package clef

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/radiustechsystems/sdk/go/src/auth"
	"github.com/radiustechsystems/sdk/go/src/common"
	"github.com/radiustechsystems/sdk/go/src/providers/eth"
)

// Signer implements the Signer interface using the Clef JSON-RPC API.
// Clef is a secure key management service that can be used to sign transactions without exposing
// the private key to the application. This is useful for securing private keys in production systems.
// Clef must be running and accessible to the application in order to use this signer.
// Learn more about Clef here: https://geth.ethereum.org/docs/tools/clef/introduction
type Signer struct {
	// address is the Radius address associated with this signer
	address common.Address

	// chainID is the network chain ID used for EIP-155 transaction signing
	chainID *big.Int

	// client is the RPC client used to communicate with the Clef server
	client *eth.RPCClient

	// signer is the underlying Ethereum signer implementation
	signer eth.Signer
}

// New creates a new Signer with the given address, Radius Client, and Clef server URL.
// @param address The address to use for signing
// @param client The Radius client
// @param clefURL The URL of the Clef server (e.g. "http://localhost:8550")
// @return A new Signer instance, or an error if the connection fails
func New(address common.Address, client auth.SignerClient, clefURL string) (*Signer, error) {
	clefClient, err := eth.NewRPCClient(clefURL, client.HTTPClient())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Clef: %w", err)
	}

	// Verify we can actually connect by checking version
	var version string
	if err = clefClient.Call(&version, "account_version"); err != nil {
		return nil, fmt.Errorf("failed to verify Clef connection: %w", err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		chainID = new(big.Int)
	}

	return &Signer{
		address: address,
		chainID: chainID,
		client:  clefClient,
		signer:  eth.NewEIP155Signer(chainID),
	}, nil
}

// Address implements the Signer interface
// @return The Radius Address associated with the Signer
func (s *Signer) Address() common.Address {
	return s.address
}

// ChainID implements the Signer interface
// @return The Chain ID associated with the Signer
func (s *Signer) ChainID() *big.Int {
	return s.chainID
}

// Hash implements the Signer interface
// @param tx The transaction to hash
// @return The hash of the given transaction
func (s *Signer) Hash(tx *common.Transaction) common.Hash {
	ethTx := tx.EthTransaction()
	ethHash := s.signer.Hash(ethTx)
	return common.NewHash(ethHash.Bytes())
}

// SignMessage implements the Signer interface
// @param msg The message bytes to sign
// @return The signature bytes, or an error if signing fails
func (s *Signer) SignMessage(msg []byte) ([]byte, error) {
	var result string // Clef returns hex string

	if err := s.client.Call(&result, "account_signData",
		"text/plain", // Use text/plain content type as shown in docs
		s.address.Hex(),
		hex.EncodeToString(msg)); err != nil {
		return nil, fmt.Errorf("clef signing failed: %w", err)
	}

	return hex.DecodeString(strings.TrimPrefix(result, "0x"))
}

// SignTransaction implements the Signer interface
// @param tx The transaction to sign
// @return The signed transaction, or an error if signing fails
func (s *Signer) SignTransaction(tx *common.Transaction) (*common.SignedTransaction, error) {
	var result signedTransaction

	args := tx.ToMap()
	args["from"] = s.address.Hex()
	args["chainId"] = fmt.Sprintf("0x%x", s.chainID)

	if err := s.client.Call(&result, "account_signTransaction", args); err != nil {
		return nil, fmt.Errorf("clef signing failed: %w", err)
	}

	return result.ToRadiusSignedTransaction(tx)
}

// signedTransaction represents a transaction signed by Clef.
// It contains the raw signed transaction data and signature components.
type signedTransaction struct {
	// Raw is the raw signed transaction data as a hex string
	Raw string `json:"raw"`

	// Tx contains the transaction hash and signature components
	Tx struct {
		// Hash is the transaction hash as a hex string
		Hash string `json:"hash"`

		// V is the v component of the signature as a hex string
		V string `json:"v"`

		// R is the r component of the signature as a hex string
		R string `json:"r"`

		// S is the s component of the signature as a hex string
		S string `json:"s"`
	} `json:"tx"`
}

// ToRadiusSignedTransaction converts a signedTransaction to a common.SignedTransaction
// @param origTx The original unsigned transaction
// @return The signed transaction, or an error if conversion fails
func (c *signedTransaction) ToRadiusSignedTransaction(origTx *common.Transaction) (*common.SignedTransaction, error) {
	r, ok := new(big.Int).SetString(strings.TrimPrefix(c.Tx.R, "0x"), 16)
	if !ok {
		return nil, fmt.Errorf("invalid R value: %s", c.Tx.R)
	}

	s, ok := new(big.Int).SetString(strings.TrimPrefix(c.Tx.S, "0x"), 16)
	if !ok {
		return nil, fmt.Errorf("invalid S value: %s", c.Tx.S)
	}

	v, ok := new(big.Int).SetString(strings.TrimPrefix(c.Tx.V, "0x"), 16)
	if !ok {
		return nil, fmt.Errorf("invalid V value: %s", c.Tx.V)
	}

	rawBytes, err := hex.DecodeString(strings.TrimPrefix(c.Raw, "0x"))
	if err != nil {
		return nil, fmt.Errorf("invalid raw transaction: %s", c.Raw)
	}

	return &common.SignedTransaction{
		Transaction: origTx,
		R:           r,
		S:           s,
		V:           v,
		Serialized:  rawBytes,
	}, nil
}
