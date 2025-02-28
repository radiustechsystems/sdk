//go:build integration

package test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/radiustechsystems/sdk/go/radius"
)

func TestIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	endpoint := SkipIfNoEndpoint(t)
	client, err := radius.NewClient(endpoint, radius.WithLogger(t.Logf))
	require.NoError(t, err, "Failed to create client")

	account := GetFundedTestAccount(t, client)
	require.NotNil(t, account, "Failed to create account")
	require.NotNil(t, account.Address(), "Account address should not be nil")

	t.Run("Send", func(t *testing.T) {
		initialBalance := SkipIfInsufficientTestAccountBalance(ctx, t, account, client)
		recipient := CreateTestAccount(t, client)
		amount := big.NewInt(100)

		// Send ETH from test account to recipient
		var receipt *radius.Receipt
		receipt, err = account.Send(ctx, client, recipient.Address(), amount)
		assert.NoError(t, err, "Failed to send value to recipient")
		require.NotNil(t, receipt, "Receipt should not be nil")
		assert.Equal(t, account.Address(), receipt.From, "Unexpected sender address")
		assert.Equal(t, recipient.Address(), receipt.To, "Unexpected recipient address")
		assert.Equal(t, amount, receipt.Value, "Unexpected value")

		// Check sender balance
		var senderBalance *big.Int
		senderBalance, err = account.Balance(ctx, client)
		assert.NoError(t, err, "Failed to get sender balance")
		assert.Equal(t, initialBalance.Sub(initialBalance, amount), senderBalance, "Unexpected sender balance")

		// Check recipient balance
		var recipientBalance *big.Int
		recipientBalance, err = recipient.Balance(ctx, client)
		assert.NoError(t, err, "Failed to get recipient balance")
		assert.Equal(t, amount, recipientBalance, "Unexpected recipient balance")
	})

	t.Run("SimpleStorage", func(t *testing.T) {
		var (
			contract *radius.Contract
			receipt  *radius.Receipt
			result   []interface{}
		)

		abi := radius.ABIFromJSON(SimpleStorageABI)
		require.NotNil(t, abi, "Failed to parse ABI")

		bytecode := radius.BytecodeFromHex(SimpleStorageBin)
		require.NotNil(t, bytecode, "Failed to parse bytecode")

		contract, err = client.DeployContract(ctx, account.Signer, bytecode, abi)
		require.NoError(t, err, "Failed to deploy contract")
		require.NotNil(t, contract.Address(), "Contract address should not be nil")
		require.NotEqual(t, radius.Address{}, contract.Address(), "Contract address should not be empty")

		value := big.NewInt(42)
		receipt, err = contract.Execute(ctx, client, account.Signer, "set", value)
		assert.NoError(t, err, "Failed to call contract method")
		assert.NotNil(t, receipt, "Receipt should not be nil")
		assert.Equal(t, account.Address(), receipt.From, "Unexpected from address")
		assert.Equal(t, contract.Address(), receipt.To, "Unexpected to address")

		result, err = contract.Call(ctx, client, "get")
		assert.NoError(t, err, "Failed to call contract method")
		assert.Len(t, result, 1, "Unexpected result length")
		assert.Equal(t, value, result[0].(*big.Int), "Unexpected result value")
	})
}
