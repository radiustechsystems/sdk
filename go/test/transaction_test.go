package test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/radiustechsystems/sdk/go/src/radius"
)

func TestNewTransaction(t *testing.T) {
	t.Run("Creates a transaction with all fields set", func(t *testing.T) {
		data := []byte("data payload")
		gas := uint64(21000)
		gasPrice := OneGwei
		nonce := uint64(42)
		toAddr := radius.NewAddressFromHex("0x8ba1f109551bD432803012645Ac136ddd64DBA72")
		value := OneETH

		tx := radius.NewTransaction(data, gas, gasPrice, nonce, &toAddr, value)
		require.NotNil(t, tx, "Transaction should not be nil")
		assert.Equal(t, data, tx.Data(), "Data field should match")
		assert.Equal(t, gas, tx.Gas(), "Gas field should match")
		assert.Equal(t, gasPrice, tx.GasPrice(), "GasPrice field should match")
		assert.Equal(t, nonce, tx.Nonce(), "Nonce field should match")
		assert.Equal(t, toAddr, *tx.To(), "To field should match")
		assert.Equal(t, value, tx.Value(), "Value field should match")
		assert.Equal(t, uint8(0x0), tx.Type(), "Transaction type should be LegacyTxType")
	})

	t.Run("Creates a transaction with nil address (contract deployment)", func(t *testing.T) {
		data := []byte("contract bytecode")
		gas := uint64(21000)
		gasPrice := big.NewInt(0)
		nonce := uint64(1)
		value := big.NewInt(0)

		tx := radius.NewTransaction(data, gas, gasPrice, nonce, nil, value)
		assert.Nil(t, tx.To(), "To field should be nil for contract deployment")
	})

	t.Run("Creates a transaction with zero values", func(t *testing.T) {
		addr := radius.NewAddressFromHex("0x0")

		tx := radius.NewTransaction([]byte{}, 0, big.NewInt(0), 0, &addr, big.NewInt(0))
		assert.Equal(t, 0, len(tx.Data()), "Data should be empty")
		assert.Equal(t, uint64(0), tx.Gas(), "Gas should be zero")
		assert.Equal(t, big.NewInt(0), tx.GasPrice(), "GasPrice should be zero")
		assert.Equal(t, uint64(0), tx.Nonce(), "Nonce should be zero")
		assert.Equal(t, big.NewInt(0), tx.Value(), "Value should be zero")
	})
}

func TestReceipt(t *testing.T) {
	t.Run("Receipt type is properly aliased", func(t *testing.T) {
		blockHash := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
		blockNumber := big.NewInt(12345)
		contractAddress := common.HexToAddress("0x0000000000000000000000000000000000000000")
		cumulativeGasUsed := uint64(21000)
		effectiveGasPrice := big.NewInt(0)
		gasUsed := uint64(21000)
		var logs []*types.Log
		bloom := types.Bloom{}
		status := uint64(1)
		txHash := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
		transactionIndex := uint(0)

		receipt := &radius.Receipt{
			Type:              uint8(0),
			PostState:         []byte{},
			Status:            status,
			CumulativeGasUsed: cumulativeGasUsed,
			Bloom:             bloom,
			Logs:              logs,
			TxHash:            txHash,
			ContractAddress:   contractAddress,
			GasUsed:           gasUsed,
			EffectiveGasPrice: effectiveGasPrice,
			BlockHash:         blockHash,
			BlockNumber:       blockNumber,
			TransactionIndex:  transactionIndex,
		}

		assert.NotNil(t, receipt, "Should be able to assign types.Receipt to radius.Receipt")
		assert.Equal(t, blockHash, receipt.BlockHash, "BlockHash field should match")
		assert.Equal(t, blockNumber, receipt.BlockNumber, "BlockNumber field should match")
		assert.Equal(t, status, receipt.Status, "Status field should match")
		assert.Equal(t, txHash, receipt.TxHash, "TxHash field should match")
		assert.Equal(t, gasUsed, receipt.GasUsed, "GasUsed field should match")
	})
}
