package common

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/radiustechsystems/sdk/go/src/providers/eth"
)

// ABIFromJSON creates a new ABI (Application Binary Interface) from a JSON string.
//
// @param json ABI definition in JSON string format
// @return A new ABI instance, or nil if the JSON is invalid
func ABIFromJSON(json string) *ABI {
	abi, err := NewABI(json)
	if err != nil {
		return nil
	}
	return abi
}

// AddressFromHex creates an Address from a hex string
// @param hexStr string representing the address in hexadecimal format (with or without 0x prefix)
// @return Address instance, or an error if the hex string is invalid
func AddressFromHex(h string) (Address, error) {
	// Remove 0x prefix if present
	cleanHex := strings.TrimPrefix(h, "0x")

	// Decode hex string to bytes
	addrBytes, err := hex.DecodeString(cleanHex)
	if err != nil {
		return Address{}, fmt.Errorf("invalid hex address: %w", err)
	}

	// Check length
	if len(addrBytes) != 20 {
		return Address{}, fmt.Errorf("address must be 20 bytes, got %d", len(addrBytes))
	}

	return NewAddress(addrBytes), nil
}

// BytecodeFromHex converts a hex string to a byte slice
// @param s Hex string (with or without 0x prefix)
// @return Byte slice representation of the hex string, or nil if the string is not valid hex
func BytecodeFromHex(s string) []byte {
	// Remove 0x prefix if present
	if len(s) >= 2 && s[0:2] == "0x" {
		s = s[2:]
	}

	bytecode, err := hex.DecodeString(s)
	if err != nil {
		return nil
	}
	return bytecode
}

// EthAddressFromRadiusAddress converts a Radius Address pointer to an Ethereum Address pointer
// @param address Radius Address pointer
// @return Ethereum Address pointer, or nil if the input is nil
func EthAddressFromRadiusAddress(address *Address) *eth.Address {
	if address == nil {
		return nil
	}
	ethAddress := address.EthAddress()
	return &ethAddress
}

// EventsFromEthLogs converts Ethereum logs to Radius events
// @param logs Ethereum logs
// @return Slice of Radius events
func EventsFromEthLogs(logs []*eth.Log) []Event {
	events := make([]Event, len(logs))
	for i, log := range logs {
		events[i] = Event{
			Name: log.Topics[0].Hex(),
			Data: make(map[string]interface{}),
			Raw:  log.Data,
		}
	}
	return events
}

// HashFromHex creates a new Hash from a hexadecimal string
// @param h The hexadecimal string representation of the hash (with or without 0x prefix)
// @return A pointer to the new Hash instance, or an error if the hex string is invalid
func HashFromHex(h string) (Hash, error) {
	// Remove 0x prefix if present
	cleanHex := strings.TrimPrefix(h, "0x")

	// Decode hex string to bytes
	hashBytes, err := hex.DecodeString(cleanHex)
	if err != nil {
		return NewHash([]byte{}), err
	}

	return NewHash(hashBytes), nil
}

// ReceiptFromEthReceipt creates a new Radius receipt from an Ethereum receipt
// @param r Ethereum receipt
// @param from Sender address
// @param to Recipient address
// @param value Transaction value
// @return Radius receipt
func ReceiptFromEthReceipt(r *eth.Receipt, from, to Address, value *big.Int) *Receipt {
	return &Receipt{
		From:            from,
		To:              to,
		ContractAddress: NewAddress(r.ContractAddress.Bytes()),
		TxHash:          NewHash(r.TxHash.Bytes()),
		GasUsed:         r.GasUsed,
		Logs:            EventsFromEthLogs(r.Logs),
		Status:          r.Status,
		Value:           value,
	}
}

// ZeroAddress returns the zero address (0x0000000000000000000000000000000000000000).
// Used as a default value or to represent the zero address in the Ethereum ecosystem.
//
// @return An Address instance representing the zero address
func ZeroAddress() Address {
	return NewAddress(make([]byte, 20))
}
