package radius

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

type Transaction = types.Transaction

func NewTransaction(
	data []byte,
	gas uint64,
	gasPrice *big.Int,
	nonce uint64,
	to *Address,
	value *big.Int,
) *types.Transaction {
	return types.NewTx(&types.LegacyTx{
		Data:     data,
		Gas:      gas,
		GasPrice: gasPrice,
		Nonce:    nonce,
		To:       to,
		Value:    value,
	})
}

type Receipt = types.Receipt
