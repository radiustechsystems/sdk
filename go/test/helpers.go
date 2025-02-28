package test

import (
	"context"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/radiustechsystems/sdk/go/radius"
)

// Constants for testing
const (
	SimpleStorageABI = `[{"inputs":[],"name":"get","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"x","type":"uint256"}],"name":"set","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
	SimpleStorageBin = "6080604052348015600e575f5ffd5b5060a580601a5f395ff3fe6080604052348015600e575f5ffd5b50600436106030575f3560e01c806360fe47b11460345780636d4ce63c146045575b5f5ffd5b6043603f3660046059565b5f55565b005b5f5460405190815260200160405180910390f35b5f602082840312156068575f5ffd5b503591905056fea26469706673582212207655d86666fa8aa75666db8416e0f5db680914358a57e84aa369d9250218247f64736f6c634300081c0033"
)

// MinTestAccountBalance is the minimum balance required for a test account
var MinTestAccountBalance = big.NewInt(1000000000000000000) // 1 ETH

// SkipIfNoPrivateKey skips the test if the RADIUS_PRIVATE_KEY environment variable is not set
func SkipIfNoPrivateKey(t *testing.T) string {
	privateKeyHex := os.Getenv("RADIUS_PRIVATE_KEY")
	if privateKeyHex == "" {
		t.Skip("RADIUS_PRIVATE_KEY environment variable not set")
	}
	return privateKeyHex
}

// GetFundedTestAccount returns a test account with funds from the private key environment variable
func GetFundedTestAccount(t *testing.T, client *radius.Client) *radius.Account {
	privateKeyHex := SkipIfNoPrivateKey(t)
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	require.NoError(t, err, "Failed to parse private key")

	account := radius.NewAccount(radius.WithPrivateKey(privateKey, client))
	require.NotNil(t, account, "Failed to create account")

	return account
}

// CreateTestAccount creates a new test account with a random private key
func CreateTestAccount(t *testing.T, client *radius.Client) *radius.Account {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err, "Failed to generate private key")

	account := radius.NewAccount(radius.WithPrivateKey(privateKey, client))
	require.NotNil(t, account, "Failed to create account")
	require.NotNil(t, account.Address(), "Account address should not be nil")

	return account
}

// SkipIfNoEndpoint skips the test if the RADIUS_ENDPOINT environment variable is not set
func SkipIfNoEndpoint(t *testing.T) string {
	endpoint := os.Getenv("RADIUS_ENDPOINT")
	if endpoint == "" {
		t.Skip("RADIUS_ENDPOINT environment variable not set")
	}
	return endpoint
}

// SkipIfInsufficientTestAccountBalance skips the test if the account doesn't have enough funds
func SkipIfInsufficientTestAccountBalance(ctx context.Context, t *testing.T, account *radius.Account, client *radius.Client) *big.Int {
	balance, err := account.Balance(ctx, client)
	require.NoError(t, err, "Failed to get balance")
	if balance.Cmp(MinTestAccountBalance) == -1 {
		t.Skip("Test account has insufficient balance")
	}
	return balance
}
