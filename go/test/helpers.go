package test

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/radiustechsystems/sdk/go/src/radius"
)

func CreateTestClient(t *testing.T, handlers map[string]func(params []interface{}) interface{}) *radius.Client {
	server := MockJSONRPCServer(t, handlers)
	client, err := radius.NewClient(server.URL)
	require.NoError(t, err, "Failed to create test client")

	return client
}

func CreateTestSigner() radius.Signer {
	key := radius.GeneratePrivateKey()

	return radius.NewPrivateKeySigner(key, TestnetChainID)
}

func CreateTestTransaction(to radius.Address) *radius.Transaction {
	return radius.NewTransaction([]byte("data"), 21000, OneGwei, 1, &to, OneETH)
}

func SkipIfInsufficientFunds(t *testing.T, account *radius.Account) *big.Int {
	balance, err := account.Balance(context.Background())
	require.NoError(t, err, "Failed to get balance")
	if balance.Cmp(MinBalance) == -1 {
		t.Skip("Test account has insufficient balance")
	}
	return balance
}

func SkipIfNoPrivateKey(t *testing.T) *ecdsa.PrivateKey {
	if PrivateKeyHex == "" {
		t.Skip("RADIUS_PRIVATE_KEY environment variable not set")
	}
	key, err := crypto.HexToECDSA(PrivateKeyHex)
	if err != nil {
		t.Skip("RADIUS_PRIVATE_KEY environment variable not set")
	}
	return key
}

func SkipIfNoRPCEndpoint(t *testing.T) string {
	if RPCEndpoint == "" {
		t.Skip("RADIUS_RPC_ENDPOINT environment variable not set")
	}
	return RPCEndpoint
}

func ToByte32(s string) [32]byte {
	var accessID [32]byte
	copy(accessID[:], crypto.Keccak256([]byte(s)))
	return accessID
}
