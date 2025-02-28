// Package common provides core data types and utilities for the Radius SDK.
//
// This package contains the fundamental data structures and utilities needed for
// interacting with Radius smart contracts, including addresses, transactions,
// events, and cryptographic functions. It provides Radius-specific implementations
// that are independent from but compatible with Ethereum libraries.
package common

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// ABI represents an Application Binary Interface for smart contracts.
// It provides methods for encoding and decoding contract method calls and return values,
// which are essential for interacting with smart contracts deployed on Radius.
type ABI struct {
	// abi is the underlying ABI implementation
	abi abi.ABI
}

// NewABI creates a new ABI instance from a JSON string representation.
//
// @param abiJSON String representing the ABI in JSON format
// @return An ABI instance if successful, or an error if the JSON string is empty or invalid
func NewABI(abiJSON string) (*ABI, error) {
	if abiJSON == "" {
		return nil, fmt.Errorf("ABI JSON string is empty")
	}

	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	return &ABI{abi: parsedABI}, nil
}

// Pack encodes contract input data for method calls or constructor invocations.
//
// @param name Name of the method to call, or an empty string for constructor
// @param args Variadic list of arguments for the method
// @return Encoded binary data ready for contract interaction, or an error if the method is not found or encoding fails
func (a *ABI) Pack(name string, args ...interface{}) ([]byte, error) {
	// Special case for constructor
	if name == "" {
		return a.abi.Pack("", args...)
	}

	// Regular method call
	data, err := a.abi.Pack(name, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to pack arguments: %w", err)
	}

	return data, nil
}

// Unpack decodes contract output data returned from a method call.
//
// @param name Name of the method that produced the output, or an empty string for constructor
// @param data Encoded binary data received from the contract
// @return List of decoded values representing the method's return values, or an error if the method is not found or decoding fails
func (a *ABI) Unpack(name string, data []byte) ([]interface{}, error) {
	// Special case for constructor which has no return value
	if name == "" {
		return []interface{}{}, nil
	}

	method, ok := a.abi.Methods[name]
	if !ok {
		return nil, fmt.Errorf("method %s not found in ABI", name)
	}

	result := make(map[string]interface{})
	if err := method.Outputs.UnpackIntoMap(result, data); err != nil {
		return nil, fmt.Errorf("failed to unpack output: %w", err)
	}

	values := make([]interface{}, len(method.Outputs))
	for i, output := range method.Outputs {
		values[i] = result[output.Name]
	}

	return values, nil
}
