package test

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/radiustechsystems/sdk/go/src/radius"
)

func TestNewClient(t *testing.T) {
	t.Run("Creates client with valid URL", func(t *testing.T) {
		server := MockJSONRPCServer(t, nil)
		defer server.Close()

		client, err := radius.NewClient(server.URL)
		require.NoError(t, err, "NewClient should not return an error with valid URL")
		require.NotNil(t, client, "Client should not be nil")

		assert.Equal(t, TestnetChainID, client.ChainID, "ChainID should match the expected value")
	})

	t.Run("Returns error with invalid URL", func(t *testing.T) {
		client, err := radius.NewClient("invalid://url")
		assert.Error(t, err, "NewClient should return an error with invalid URL")
		assert.Nil(t, client, "Client should be nil when error occurs")
	})

	t.Run("Returns error when server is unavailable", func(t *testing.T) {
		server := MockJSONRPCServer(t, nil)
		server.Close() // Close server to force connection error

		client, err := radius.NewClient(server.URL)
		assert.Error(t, err, "NewClient should return an error when server is unavailable")
		assert.Nil(t, client, "Client should be nil when error occurs")
	})
}

func TestClient_API(t *testing.T) {
	t.Run("Successfully calls RPC method", func(t *testing.T) {
		blockNumber := time.Now().UnixMilli()
		expectedResult := fmt.Sprintf("0x%x", blockNumber)
		server := MockJSONRPCServer(t, map[string]func(params []interface{}) interface{}{
			"eth_blockNumber": func(_ []interface{}) interface{} {
				return expectedResult
			},
		})
		defer server.Close()

		client, err := radius.NewClient(server.URL)
		require.NoError(t, err, "Failed to create client")

		var result string
		ctx := context.Background()
		err = client.API(ctx, &result, "eth_blockNumber")
		assert.NoError(t, err, "API call should not return an error")
		assert.Equal(t, expectedResult, result, "Result should match expected value")

		intValue, err := strconv.ParseInt(result, 0, 64)
		require.NoError(t, err, "Failed to parse result as integer")

		// Radius uses Unix timestamps (in milliseconds) as block numbers, because Radius is not a blockchain
		expectedTimestamp := time.UnixMilli(blockNumber)
		actualTimestamp := time.UnixMilli(intValue)
		assert.Equal(t, expectedTimestamp, actualTimestamp, "Result should match expected value")
	})
}

func TestClient_BalanceAt(t *testing.T) {
	t.Run("Returns balance for address", func(t *testing.T) {
		expectedResult := OneETH
		server := MockJSONRPCServer(t, map[string]func(params []interface{}) interface{}{
			"eth_getBalance": func(_ []interface{}) interface{} {
				return fmt.Sprintf("0x%x", expectedResult)
			},
		})
		defer server.Close()

		client, err := radius.NewClient(server.URL)
		require.NoError(t, err, "Failed to create client")

		address := radius.NewAddressFromHex("0x8ba1f109551bD432803012645Ac136ddd64DBA72")
		ctx := context.Background()

		balance, err := client.BalanceAt(ctx, address)
		assert.NoError(t, err, "BalanceAt should not return an error")
		assert.Equal(t, expectedResult, balance, "Balance should match expected value (from mock)")
	})
}

func TestClient_CodeAt(t *testing.T) {
	t.Run("Retrieves code at address", func(t *testing.T) {
		expectedResult := []byte(SimpleStorageBin)
		server := MockJSONRPCServer(t, map[string]func(params []interface{}) interface{}{
			"eth_getCode": func(_ []interface{}) interface{} {
				return fmt.Sprintf("0x%x", []byte(SimpleStorageBin))
			},
		})
		defer server.Close()

		client, err := radius.NewClient(server.URL)
		require.NoError(t, err, "Failed to create client")

		address := radius.NewAddressFromHex("0x8ba1f109551bD432803012645Ac136ddd64DBA72")
		ctx := context.Background()

		code, err := client.CodeAt(ctx, address)
		require.NoError(t, err, "CodeAt should not return an error")

		assert.Equal(t, expectedResult, code, "Code should match expected value")
	})
}

func TestClient_EstimateGas(t *testing.T) {
	t.Run("Estimates gas with safety margin", func(t *testing.T) {
		gasEstimate := uint64(21000)
		expectedResult := uint64(float64(gasEstimate) * radius.GasEstimateMultiplier)
		server := MockJSONRPCServer(t, map[string]func(params []interface{}) interface{}{
			"eth_estimateGas": func(_ []interface{}) interface{} {
				return fmt.Sprintf("0x%x", gasEstimate)
			},
		})

		client, err := radius.NewClient(server.URL)
		require.NoError(t, err, "Failed to create client")

		fromAddr := radius.Address{}
		toAddr := radius.NewAddressFromHex("0x8ba1f109551bD432803012645Ac136ddd64DBA72")
		tx := CreateTestTransaction(toAddr)
		ctx := context.Background()

		estimate, err := client.EstimateGas(ctx, tx, fromAddr)
		require.NoError(t, err, "EstimateGas should not return an error")

		assert.Equal(t, expectedResult, estimate, "Gas estimate should match expected value")
	})
}

func TestClient_NewAccount(t *testing.T) {
	t.Run("Creates account with private key", func(t *testing.T) {
		server := MockJSONRPCServer(t, nil)
		defer server.Close()

		client, err := radius.NewClient(server.URL)
		require.NoError(t, err, "Failed to create client")

		privateKey := radius.GeneratePrivateKey()
		account, err := client.AccountFromPrivateKey(privateKey)
		require.NoError(t, err, "AccountFromPrivateKey should not return an error with valid private key")
		require.NotNil(t, account, "Account should not be nil")

		expectedAddress := radius.NewAddressFromPrivateKey(privateKey)
		assert.Equal(t, expectedAddress, account.Address(), "Account address should match derived address")
	})

	t.Run("Returns error with nil private key", func(t *testing.T) {
		server := MockJSONRPCServer(t, nil)
		defer server.Close()

		client, err := radius.NewClient(server.URL)
		require.NoError(t, err, "Failed to create client")

		account, err := client.AccountFromPrivateKey(nil)
		assert.Error(t, err, "AccountFromPrivateKey should return an error with nil private key")
		assert.Nil(t, account, "Account should be nil when error occurs")
	})
}

func TestClient_Nonce(t *testing.T) {
	t.Run("Returns nonce for address", func(t *testing.T) {
		expectedResult := uint64(42)
		server := MockJSONRPCServer(t, map[string]func(params []interface{}) interface{}{
			"eth_getTransactionCount": func(_ []interface{}) interface{} {
				return fmt.Sprintf("0x%x", expectedResult)
			},
		})
		defer server.Close()

		client, err := radius.NewClient(server.URL)
		require.NoError(t, err, "Failed to create client")

		address := radius.NewAddressFromHex("0x8ba1f109551bD432803012645Ac136ddd64DBA72")
		ctx := context.Background()

		nonce, err := client.Nonce(ctx, address)
		assert.NoError(t, err, "Nonce should not return an error")
		assert.Equal(t, expectedResult, nonce, "Nonce should match expected value (from mock)")
	})
}

func TestClient_PrepareTx(t *testing.T) {
	t.Run("Prepares transaction with nil signer", func(t *testing.T) {
		server := MockJSONRPCServer(t, nil)
		defer server.Close()

		client, err := radius.NewClient(server.URL)
		require.NoError(t, err, "Failed to create client")

		data := []byte("test data")
		to := radius.NewAddressFromHex("0x8ba1f109551bD432803012645Ac136ddd64DBA72")
		value := OneETH
		ctx := context.Background()

		tx, err := client.PrepareTx(ctx, data, nil, &to, value)
		require.NoError(t, err, "PrepareTx should not return an error with nil signer")
		require.NotNil(t, tx, "Transaction should not be nil")

		assert.Equal(t, data, tx.Data(), "Transaction data should match")
		assert.Equal(t, uint64(0), tx.Gas(), "Gas should be zero with nil signer")
		assert.Equal(t, big.NewInt(0), tx.GasPrice(), "GasPrice should be zero")
		assert.Equal(t, uint64(0), tx.Nonce(), "Nonce should be zero with nil signer")
		assert.Equal(t, to, *tx.To(), "To address should match")
		assert.Equal(t, value, tx.Value(), "Value should match")
	})

	t.Run("Prepares transaction with signer", func(t *testing.T) {
		gasEstimate := uint64(21000)
		expectedGasWithMargin := uint64(float64(gasEstimate) * radius.GasEstimateMultiplier)
		expectedGasPrice := big.NewInt(0)
		expectedNonce := uint64(42)
		server := MockJSONRPCServer(t, map[string]func(params []interface{}) interface{}{
			"eth_estimateGas": func(_ []interface{}) interface{} {
				return fmt.Sprintf("0x%x", gasEstimate)
			},
			"eth_getTransactionCount": func(_ []interface{}) interface{} {
				return fmt.Sprintf("0x%x", expectedNonce)
			},
		})
		defer server.Close()

		client, err := radius.NewClient(server.URL)
		require.NoError(t, err, "Failed to create client")

		privateKey := radius.GeneratePrivateKey()
		signer := radius.NewPrivateKeySigner(privateKey, TestnetChainID)
		ctx := context.Background()
		data := []byte("test data")
		to := radius.NewAddressFromHex("0x8ba1f109551bD432803012645Ac136ddd64DBA72")
		value := OneETH

		tx, err := client.PrepareTx(ctx, data, signer, &to, value)
		require.NoError(t, err, "PrepareTx should not return an error with signer")
		require.NotNil(t, tx, "Transaction should not be nil")

		assert.Equal(t, data, tx.Data(), "Transaction data should match")
		assert.Equal(t, expectedGasWithMargin, tx.Gas(), "Gas should match expected value")
		assert.Equal(t, expectedGasPrice, tx.GasPrice(), "GasPrice should match expected value")
		assert.Equal(t, expectedNonce, tx.Nonce(), "Nonce should match expected value")
		assert.Equal(t, to, *tx.To(), "To address should match")
		assert.Equal(t, value, tx.Value(), "Value should match")
	})
}
