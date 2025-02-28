package common

import (
	"bytes"

	"github.com/radiustechsystems/sdk/go/src/providers/eth"
)

// Address represents a 20-byte Radius account or contract address.
// This struct provides methods to convert between different address representations
// and compare addresses. It serves as the core data structure for identifying
// accounts and smart contracts in the Radius system.
type Address struct {
	// data is the underlying 20-byte address data
	data [20]byte
}

// NewAddress creates an Address instance from a byte slice.
//
// @param b Byte slice representing the address (should be 20 bytes)
// @return Address instance initialized with the provided bytes
func NewAddress(b []byte) Address {
	var a Address
	copy(a.data[:], b)
	return a
}

// Bytes returns the address as a byte slice.
//
// @return Byte slice representation of the 20-byte address
func (a *Address) Bytes() []byte {
	return a.data[:]
}

// EthAddress converts a Radius Address to an eth.Address.
// This method is used when Ethereum library functionality is needed.
//
// @return eth.Address representation of the Radius address
func (a *Address) EthAddress() eth.Address {
	return eth.BytesToAddress(a.Bytes())
}

// ToEthAddress converts a Radius Address to an eth.Address.
//
// @deprecated Use EthAddress instead
// @return eth.Address representation of the address
func (a *Address) ToEthAddress() eth.Address {
	return a.EthAddress()
}

// Hex returns the hexadecimal string representation of the address.
// The returned address string is properly checksummed according to EIP-55.
//
// @return Hex string representation of the address with 0x prefix
func (a *Address) Hex() string {
	// Use the go-ethereum checksum function instead of just formatting as hex
	return eth.BytesToAddress(a.Bytes()).Hex()
}

// Equals compares this address with another address for equality.
//
// @param other Address to compare with this address
// @return true if addresses contain identical bytes, false otherwise
func (a *Address) Equals(other Address) bool {
	return bytes.Equal(a.data[:], other.data[:])
}
