package radius

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// ABI is an Application Binary Interface
//
// The ABI holds information about a contract's context and available invocable methods. It enables you to type check
// function calls and packs data accordingly.
type ABI = abi.ABI

// NewABI creates a new ABI from a JSON string
func NewABI(s string) (ABI, error) {
	return abi.JSON(strings.NewReader(s))
}
