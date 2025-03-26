//go:build integration

package test

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/radiustechsystems/sdk/go/src/radius"
)

func TestIntegration(t *testing.T) {
	var (
		account *radius.Account
		client  *radius.Client
		err     error
	)

	url := SkipIfNoRPCEndpoint(t)
	key := SkipIfNoPrivateKey(t)

	client, err = radius.NewClientWithLogging(url, t.Logf)
	require.NoError(t, err, "Failed to create integration test client")

	account, err = client.AccountFromPrivateKey(key)
	require.NoError(t, err, "Failed to create integration test account")

	balance := SkipIfInsufficientFunds(t, account)
	t.Log(fmt.Sprintf("Integration test account balance: %s", balance.String()))

	t.Run("Send value to another account", func(t *testing.T) {
		var (
			newBalance   *big.Int
			receipt      *radius.Receipt
			toBalance    *big.Int
			toNewBalance *big.Int
		)

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
		defer cancel()

		// Create a new account to send value to
		toAddr := radius.NewAddressFromPrivateKey(radius.GeneratePrivateKey())
		toBalance, err = client.BalanceAt(ctx, toAddr)

		// Send value from account to recipient
		receipt, err = client.Send(ctx, account.Signer, toAddr, OneGwei)
		require.NoError(t, err, "Failed to send value to recipient")
		assert.NotNil(t, receipt, "Receipt should not be nil")
		assert.Equal(t, uint64(1), receipt.Status, "Receipt status should be 1")

		// Confirm sender balance decreased
		newBalance, err = account.Balance(ctx)
		require.NoError(t, err, "Failed to get sender balance")
		assert.Equal(t, balance.Sub(balance, OneGwei), newBalance, "Unexpected sender balance")

		// Confirm recipient balance increased
		toNewBalance, err = client.BalanceAt(ctx, toAddr)
		require.NoError(t, err, "Failed to get recipient balance")
		assert.Equal(t, toBalance.Add(toBalance, OneGwei), toNewBalance, "Unexpected recipient balance")
	})

	t.Run("Deploy and interact with a contract", func(t *testing.T) {
		var (
			contract *radius.Contract
			receipt  *radius.Receipt
			result   []interface{}
		)

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
		defer cancel()

		contract, err = client.DeployContractFromStrings(ctx, account.Signer, SimpleStorageABI, SimpleStorageBin)
		require.NoError(t, err, "Failed to deploy contract")
		assert.NotNil(t, contract, "Contract should not be nil")
		assert.NotNil(t, contract.Address(), "Contract address should not be nil")
		assert.NotEqual(t, radius.Address{}, contract.Address(), "Contract address should not be empty")

		// Set value in contract
		expectedValue := big.NewInt(42)
		receipt, err = contract.Exec(ctx, account.Signer, "set", expectedValue)
		require.NoError(t, err, "Failed to execute contract method")
		assert.NotNil(t, receipt, "Receipt should not be nil")
		assert.Equal(t, uint64(1), receipt.Status, "Receipt status should be 1")

		// Get value from contract
		result, err = contract.Call(ctx, "get")
		require.NoError(t, err, "Failed to call contract method")
		assert.Equal(t, expectedValue, result[0], "Unexpected contract value")
	})

	t.Run("Interact with a previously deployed contract", func(t *testing.T) {
		var (
			c        *radius.Contract
			contract *radius.Contract
			receipt  *radius.Receipt
			result   []interface{}
		)

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
		defer cancel()

		c, err = client.DeployContractFromStrings(ctx, account.Signer, SimpleStorageABI, SimpleStorageBin)
		require.NoError(t, err, "Failed to deploy contract")
		require.NotNil(t, c, "Contract should not be nil")

		contractAddr := *c.Address()
		require.NotNil(t, contractAddr, "Contract address should not be nil")
		require.NotEqual(t, radius.Address{}, contractAddr, "Contract address should not be empty")

		// Create a new contract instance using the address of the previously deployed contract
		contract = radius.NewContract(contractAddr, c.ABI, client)
		assert.Equal(t, contractAddr, *contract.Address(), "Unexpected contract address")

		// Set value in contract
		expectedValue := big.NewInt(42)
		receipt, err = contract.Exec(ctx, account.Signer, "set", expectedValue)
		require.NoError(t, err, "Failed to execute contract method")
		assert.NotNil(t, receipt, "Receipt should not be nil")
		assert.Equal(t, uint64(1), receipt.Status, "Receipt status should be 1")

		// Get value from contract
		result, err = contract.Call(ctx, "get")
		require.NoError(t, err, "Failed to call contract method")
		assert.Equal(t, expectedValue, result[0], "Unexpected contract value")
	})
}
