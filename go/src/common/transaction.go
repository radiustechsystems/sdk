package common

import (
	"fmt"
	"math/big"

	"github.com/radiustechsystems/sdk/go/src/providers/eth"
)

// Transaction is a Radius EVM transaction.
// Contains all the data needed to execute a Radius transaction.
type Transaction struct {
	// Data is the calldata for the transaction (bytecode for contract creation, or method call data)
	Data []byte

	// Gas is the maximum amount of gas units the transaction can consume
	Gas uint64

	// GasPrice is the price per gas unit in wei
	GasPrice *big.Int

	// Nonce is the sequential transaction number for the sending account
	Nonce uint64

	// To is the destination address (nil for contract creation)
	To *Address

	// Value is the amount of native currency to send in wei
	Value *big.Int
}

// EthTransaction converts the Radius Transaction to an eth.Transaction.
//
// @return The transaction converted to an eth.Transaction
func (t *Transaction) EthTransaction() *eth.Transaction {
	return eth.NewTx(&eth.LegacyTx{
		Data:     t.Data,
		Gas:      t.Gas,
		GasPrice: t.GasPrice,
		Nonce:    t.Nonce,
		To:       EthAddressFromRadiusAddress(t.To),
		Value:    t.Value,
	})
}

// ToEthTransaction returns the Transaction as an eth.Transaction.
// @deprecated Use EthTransaction instead
func (t *Transaction) ToEthTransaction() *eth.Transaction {
	return t.EthTransaction()
}

// ToMap returns the Transaction as a map.
func (t *Transaction) ToMap() map[string]interface{} {
	m := map[string]interface{}{
		"nonce": fmt.Sprintf("0x%x", t.Nonce),
		"gas":   fmt.Sprintf("0x%x", t.Gas),
		"data":  fmt.Sprintf("0x%x", t.Data),
	}

	if t.GasPrice == nil {
		m["gasPrice"] = "0x0"
	} else {
		m["gasPrice"] = fmt.Sprintf("0x%x", t.GasPrice)
	}

	if t.Value == nil {
		m["value"] = "0x0"
	} else {
		m["value"] = fmt.Sprintf("0x%x", t.Value)
	}

	if t.To != nil {
		m["to"] = t.To.Hex()
	}

	return m
}

// SignedTransaction is a cryptographically signed Radius EVM transaction
// ready to be sent to Radius. The R, S, and V fields are the raw ECDSA signature values.
type SignedTransaction struct {
	// Transaction is the underlying unsigned transaction
	*Transaction

	// R is the ECDSA signature r value
	R *big.Int

	// S is the ECDSA signature s value
	S *big.Int

	// V is the ECDSA signature v value (recovery id)
	V *big.Int

	// Serialized is the RLP-encoded signed transaction bytes
	Serialized []byte
}

// EthSignedTransaction converts the SignedTransaction to an eth.Transaction.
//
// @return The signed transaction converted to an eth.Transaction
func (s *SignedTransaction) EthSignedTransaction() *eth.Transaction {
	ltx := eth.LegacyTx{
		Data:     s.Data,
		Gas:      s.Gas,
		GasPrice: s.GasPrice,
		Nonce:    s.Nonce,
		To:       EthAddressFromRadiusAddress(s.To),
		Value:    s.Value,
		R:        s.R,
		S:        s.S,
		V:        s.V,
	}
	return eth.NewTx(&ltx)
}
