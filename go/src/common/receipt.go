package common

import (
	"math/big"
)

// Receipt represents the result of a successfully mined transaction.
// Contains information about the transaction execution including gas usage,
// emitted events, and contract creation if applicable.
type Receipt struct {
	// From is the address of the sender
	From Address

	// To is the address of the recipient
	To Address

	// ContractAddress is the address of the created contract (if any)
	ContractAddress Address

	// Value is the amount of ETH transferred
	Value *big.Int

	// GasUsed is the amount of gas used by the transaction
	GasUsed uint64

	// TxHash is the transaction hash
	TxHash Hash

	// Logs is the list of events emitted by the transaction
	Logs []Event

	// Status is the transaction status (1 for success, 0 for failure)
	Status uint64
}

// NewReceipt creates a new receipt
// @param from The sender address
// @param to The recipient address
// @param contractAddress The created contract address (if any)
// @param hash The transaction hash
// @param gasUsed The amount of gas used
// @param status The transaction status
// @param logs The transaction logs/events
// @param value The amount of ETH transferred
// @return A new Receipt instance
func NewReceipt(
	from Address,
	to Address,
	contractAddress Address,
	hash Hash,
	gasUsed uint64,
	status uint64,
	logs []Event,
	value *big.Int,
) *Receipt {
	return &Receipt{
		From:            from,
		To:              to,
		ContractAddress: contractAddress,
		TxHash:          hash,
		GasUsed:         gasUsed,
		Status:          status,
		Logs:            logs,
		Value:           value,
	}
}
