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

func TestAccessTokenSystemIntegration_Deployment(t *testing.T) {
	var (
		provider *radius.Account
		consumer *radius.Account
		client   *radius.Client
		err      error
	)

	url := SkipIfNoRPCEndpoint(t)
	key := SkipIfNoPrivateKey(t)
	consumerKey := radius.GeneratePrivateKey()

	client, err = radius.NewClientWithLogging(url, t.Logf)
	require.NoError(t, err, "Failed to create integration test client")

	provider, err = client.AccountFromPrivateKey(key)
	require.NoError(t, err, "Failed to create provider account")

	consumer, err = client.AccountFromPrivateKey(consumerKey)
	require.NoError(t, err, "Failed to create consumer account")

	balance := SkipIfInsufficientFunds(t, provider)
	t.Log("Provider account balance:", balance.String())

	var (
		tokenSystem *radius.Contract
		receipt     *radius.Receipt
		tierId      uint64 = 1
		price       *big.Int
		ttl         *big.Int
		active      bool
		result      []interface{}
		tierPrice   *big.Int
		tierTTL     *big.Int
		tierActive  bool
	)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	defer cancel()

	price = big.NewInt(1000000000) // 1 Gwei
	ttl = big.NewInt(86400)        // 1 day in seconds
	active = true

	// Deploy AccessTokenSystem
	tokenSystem, err = client.DeployContractFromStrings(
		ctx,
		provider.Signer,
		AccessTokenSystemABI,
		AccessTokenSystemBin,
		"https://example.com/api/token/{id}.json",
	)
	require.NoError(t, err, "Failed to deploy AccessTokenSystem")
	assert.NotNil(t, tokenSystem, "TokenSystem should not be nil")
	assert.NotNil(t, tokenSystem.Address(), "TokenSystem address should not be nil")

	// Create a tier
	receipt, err = tokenSystem.Exec(
		ctx,
		provider.Signer,
		"createTier",
		big.NewInt(int64(tierId)),
		price,
		ttl,
		active,
	)
	require.NoError(t, err, "Failed to create tier")
	assert.NotNil(t, receipt, "Receipt should not be nil")
	assert.Equal(t, uint64(1), receipt.Status, "Receipt status should be 1")

	// Check tier info
	result, err = tokenSystem.Call(ctx, "tiers", big.NewInt(int64(tierId)))
	require.NoError(t, err, "Failed to get tier info")
	require.Len(t, result, 3, "tiers should return 3 values")

	tierPrice = result[0].(*big.Int)
	tierTTL = result[1].(*big.Int)
	tierActive = result[2].(bool)

	assert.Equal(t, price, tierPrice, "Unexpected price")
	assert.Equal(t, ttl, tierTTL, "Unexpected TTL")
	assert.Equal(t, active, tierActive, "Unexpected active status")

	// Test setTierStatus as consumer (value should not change)
	receipt, err = tokenSystem.Exec(
		ctx,
		consumer.Signer,
		"setTierStatus",
		big.NewInt(int64(tierId)),
		false,
	)
	require.Error(t, err, "Expected error when calling setTierStatus as non-owner")
	require.Nil(t, receipt, "Receipt should be nil")

	// Test setTierStatus as provider (value should change)
	receipt, err = tokenSystem.Exec(
		ctx,
		provider.Signer,
		"setTierStatus",
		big.NewInt(int64(tierId)),
		false,
	)
	require.NoError(t, err, "Failed to set tier status")
	require.NotNil(t, receipt, "Receipt should not be nil")
	require.Equal(t, uint64(1), receipt.Status, "Receipt status should be 1")

	// Check tier is inactive
	result, err = tokenSystem.Call(ctx, "tiers", big.NewInt(int64(tierId)))
	require.NoError(t, err, "Failed to get tier info")
	tierActive = result[2].(bool)
	assert.False(t, tierActive, "Tier should be inactive")

	// Set tier active again for next tests
	receipt, err = tokenSystem.Exec(
		ctx,
		provider.Signer,
		"setTierStatus",
		big.NewInt(int64(tierId)),
		true,
	)
	require.NoError(t, err, "Failed to set tier status")
}

func TestAccessTokenSystemIntegration_PurchaseAndUseAccessToken(t *testing.T) {
	var (
		provider *radius.Account
		consumer *radius.Account
		client   *radius.Client
		err      error
	)

	url := SkipIfNoRPCEndpoint(t)
	key := SkipIfNoPrivateKey(t)
	consumerKey := radius.GeneratePrivateKey()

	client, err = radius.NewClientWithLogging(url, t.Logf)
	require.NoError(t, err, "Failed to create integration test client")

	provider, err = client.AccountFromPrivateKey(key)
	require.NoError(t, err, "Failed to create provider account")

	consumer, err = client.AccountFromPrivateKey(consumerKey)
	require.NoError(t, err, "Failed to create consumer account")

	balance := SkipIfInsufficientFunds(t, provider)
	t.Log("Provider account balance:", balance.String())

	var (
		tokenSystem *radius.Contract
		receipt     *radius.Receipt
		tierId      uint64 = 1
		price       *big.Int
		ttl         *big.Int
		active      bool
		isValid     bool
		expiryTime  *big.Int
	)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
	defer cancel()

	_, err = client.Send(ctx, provider.Signer, consumer.Address(), OneGwei)
	require.NoError(t, err, "Failed to pre-fund consumer account from provider account")

	price = big.NewInt(1000000000) // 1 Gwei
	ttl = big.NewInt(86400)        // 1 day in seconds
	active = true

	// Deploy AccessTokenSystem
	tokenSystem, err = client.DeployContractFromStrings(
		ctx,
		provider.Signer,
		AccessTokenSystemABI,
		AccessTokenSystemBin,
		"https://example.com/api/token/{id}.json",
	)
	require.NoError(t, err, "Failed to deploy AccessTokenSystem")

	// Create a tier
	receipt, err = tokenSystem.Exec(
		ctx,
		provider.Signer,
		"createTier",
		big.NewInt(int64(tierId)),
		price,
		ttl,
		active,
	)
	require.NoError(t, err, "Failed to create tier")

	// Purchase access with consumer provider
	initialBalance, err := consumer.Balance(ctx)
	require.NoError(t, err, "Failed to get consumer account balance")

	// Call purchaseAccess on the contract
	receipt, err = tokenSystem.ExecWithValue(
		ctx,
		consumer.Signer,
		price,
		"purchaseAccess",
		big.NewInt(int64(tierId)),
	)
	require.NoError(t, err, "Failed to purchase access")
	assert.NotNil(t, receipt, "Receipt should not be nil")
	assert.Equal(t, uint64(1), receipt.Status, "Receipt status should be 1")

	// Check consumer balance decreased by approximately price (accounting for gas fees)
	newBalance, err := consumer.Balance(ctx)
	require.NoError(t, err, "Failed to get consumer account balance")
	assert.True(t, initialBalance.Cmp(newBalance) > 0, "Consumer balance should have decreased")

	// Check token balance
	result, err := tokenSystem.Call(ctx, "balanceOf", consumer.Address(), big.NewInt(int64(tierId)))
	require.NoError(t, err, "Failed to get token balance")
	balance = result[0].(*big.Int)
	assert.Equal(t, big.NewInt(1), balance, "Consumer should have 1 token")

	// Check expiration time
	result, err = tokenSystem.Call(ctx, "expiresAt", consumer.Address(), big.NewInt(int64(tierId)))
	require.NoError(t, err, "Failed to get expiration time")
	expiryTime = result[0].(*big.Int)
	assert.True(t, expiryTime.Cmp(big.NewInt(0)) > 0, "Expiry time should be set")

	// Check isValid
	result, err = tokenSystem.Call(ctx, "isValid", consumer.Address(), big.NewInt(int64(tierId)))
	require.NoError(t, err, "Failed to check isValid")
	isValid = result[0].(bool)
	assert.True(t, isValid, "Token should be valid")

	// Try to revoke access as consumer (should fail)
	receipt, err = tokenSystem.Exec(
		ctx,
		consumer.Signer,
		"revokeAccess",
		consumer.Address(),
		big.NewInt(int64(tierId)),
	)
	require.Error(t, err, "Expected error when calling revokeAccess as consumer")

	// Revoke access as provider
	receipt, err = tokenSystem.Exec(
		ctx,
		provider.Signer,
		"revokeAccess",
		consumer.Address(),
		big.NewInt(int64(tierId)),
	)
	require.NoError(t, err, "Failed to revoke access")
	assert.NotNil(t, receipt, "Receipt should not be nil")
	assert.Equal(t, uint64(1), receipt.Status, "Receipt status should be 1")

	// Check token is no longer valid
	result, err = tokenSystem.Call(ctx, "isValid", consumer.Address(), big.NewInt(int64(tierId)))
	require.NoError(t, err, "Failed to check isValid")
	isValid = result[0].(bool)
	assert.False(t, isValid, "Token should not be valid after revocation")

	// Check revocation status using bitmapping
	result, err = tokenSystem.Call(ctx, "revocations", consumer.Address())
	require.NoError(t, err, "Failed to get revocation status")
	revocationBits := result[0].(*big.Int)

	// Check if the bit for tierId is set (assuming tierId = 1)
	expectedBit := big.NewInt(1).Lsh(big.NewInt(1), uint(tierId%256))
	expectedRevocation := big.NewInt(0).And(revocationBits, expectedBit)
	assert.NotEqual(t, big.NewInt(0), expectedRevocation, "Revocation bit should be set for tier")
}

func TestAccessTokenSystemIntegration_VerifyAccessToken(t *testing.T) {
	var (
		provider *radius.Account
		consumer *radius.Account
		client   *radius.Client
		err      error
	)

	url := SkipIfNoRPCEndpoint(t)
	key := SkipIfNoPrivateKey(t)
	consumerKey := radius.GeneratePrivateKey()

	client, err = radius.NewClientWithLogging(url, t.Logf)
	require.NoError(t, err, "Failed to create integration test client")

	provider, err = client.AccountFromPrivateKey(key)
	require.NoError(t, err, "Failed to create provider account")

	consumer, err = client.AccountFromPrivateKey(consumerKey)
	require.NoError(t, err, "Failed to create consumer account")

	balance := SkipIfInsufficientFunds(t, provider)
	t.Log("Provider account balance:", balance.String())

	var (
		tokenSystem   *radius.Contract
		receipt       *radius.Receipt
		tierId        uint64 = 1
		price         *big.Int
		ttl           *big.Int
		active        bool
		verifyResult  bool
		invalidResult bool
	)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
	defer cancel()

	_, err = client.Send(ctx, provider.Signer, consumer.Address(), OneGwei)
	require.NoError(t, err, "Failed to pre-fund consumer account from provider account")

	price = big.NewInt(1000000000) // 1 Gwei
	ttl = big.NewInt(86400)        // 1 day in seconds
	active = true

	// Deploy AccessTokenSystem
	tokenSystem, err = client.DeployContractFromStrings(
		ctx,
		provider.Signer,
		AccessTokenSystemABI,
		AccessTokenSystemBin,
		"https://example.com/api/token/{id}.json",
	)
	require.NoError(t, err, "Failed to deploy AccessTokenSystem")

	// Create a tier
	receipt, err = tokenSystem.Exec(
		ctx,
		provider.Signer,
		"createTier",
		big.NewInt(int64(tierId)),
		price,
		ttl,
		active,
	)
	require.NoError(t, err, "Failed to create tier")
	require.NotNil(t, receipt, "Receipt should not be nil")
	require.Equal(t, uint64(1), receipt.Status, "Receipt status should be 1")

	// Purchase access
	receipt, err = tokenSystem.ExecWithValue(
		ctx,
		consumer.Signer,
		price,
		"purchaseAccess",
		big.NewInt(int64(tierId)),
	)
	require.NoError(t, err, "Failed to purchase access")
	require.NotNil(t, receipt, "Receipt should not be nil")
	require.Equal(t, uint64(1), receipt.Status, "Receipt status should be 1")

	// Generate a random challenge as bytes32
	challenge := fmt.Sprintf("auth-challenge-%d", time.Now().UnixNano())
	t.Logf("Generated challenge: %s", challenge)

	// Sign the challenge with consumer's account
	consumerSignature, err := consumer.Signer.Sign([]byte(challenge))
	require.NoError(t, err, "Failed to sign challenge with consumer account")

	// Call verifyAccess with consumer's signature (should return true)
	result, err := tokenSystem.Call(
		ctx,
		"verifyAccess",
		consumer.Address(),
		big.NewInt(int64(tierId)),
		challenge,
		consumerSignature,
	)
	require.NoError(t, err, "Failed to call verifyAccess with consumer signature")
	require.Len(t, result, 1, "verifyAccess should return 1 value")
	verifyResult = result[0].(bool)
	assert.True(t, verifyResult, "verifyAccess should return true for consumer signature")

	// Sign with provider's account (should return false)
	providerSignature, err := provider.Signer.Sign([]byte(challenge))
	require.NoError(t, err, "Failed to sign challenge with provider account")

	result, err = tokenSystem.Call(
		ctx,
		"verifyAccess",
		consumer.Address(),
		big.NewInt(int64(tierId)),
		challenge,
		providerSignature,
	)
	require.NoError(t, err, "Failed to call verifyAccess with provider signature")
	require.Len(t, result, 1, "verifyAccess should return 1 value")
	invalidResult = result[0].(bool)
	assert.False(t, invalidResult, "verifyAccess should return false for provider signature")

	// Modify challenge and verify (should return false)
	modifiedChallenge := fmt.Sprintf("auth-modified-challenge-%d", time.Now().UnixNano())
	result, err = tokenSystem.Call(
		ctx,
		"verifyAccess",
		consumer.Address(),
		big.NewInt(int64(tierId)),
		modifiedChallenge,
		consumerSignature,
	)
	require.NoError(t, err, "Failed to call verifyAccess with modified challenge")
	require.Len(t, result, 1, "verifyAccess should return 1 value")
	invalidResult = result[0].(bool)
	assert.False(t, invalidResult, "verifyAccess should return false for modified challenge")

	// Revoke token and verify signature (should return false even with valid signature)
	receipt, err = tokenSystem.Exec(
		ctx,
		provider.Signer,
		"revokeAccess",
		consumer.Address(),
		big.NewInt(int64(tierId)),
	)
	require.NoError(t, err, "Failed to revoke access")

	result, err = tokenSystem.Call(
		ctx,
		"verifyAccess",
		consumer.Address(),
		big.NewInt(int64(tierId)),
		challenge,
		consumerSignature,
	)
	require.NoError(t, err, "Failed to call verifyAccess after revocation")
	require.Len(t, result, 1, "verifyAccess should return 1 value")
	invalidResult = result[0].(bool)
	assert.False(t, invalidResult, "verifyAccess should return false after revocation")
}

func TestAccessTokenSystemIntegration_BatchOperations(t *testing.T) {
	var (
		provider *radius.Account
		consumer *radius.Account
		client   *radius.Client
		err      error
	)

	url := SkipIfNoRPCEndpoint(t)
	key := SkipIfNoPrivateKey(t)
	consumerKey := radius.GeneratePrivateKey()

	client, err = radius.NewClientWithLogging(url, t.Logf)
	require.NoError(t, err, "Failed to create integration test client")

	provider, err = client.AccountFromPrivateKey(key)
	require.NoError(t, err, "Failed to create provider account")

	consumer, err = client.AccountFromPrivateKey(consumerKey)
	require.NoError(t, err, "Failed to create consumer account")

	balance := SkipIfInsufficientFunds(t, provider)
	t.Log("Provider account balance:", balance.String())

	var (
		tokenSystem *radius.Contract
		receipt     *radius.Receipt
		tierIds     []*big.Int
		prices      []*big.Int
		ttls        []*big.Int
		actives     []bool
		totalPrice  *big.Int
		isValid     bool
	)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
	defer cancel()

	_, err = client.Send(ctx, provider.Signer, consumer.Address(), new(big.Int).Div(OneETH, big.NewInt(100)))
	require.NoError(t, err, "Failed to pre-fund consumer account from provider account")

	// Create data for multiple tiers
	tierIds = []*big.Int{big.NewInt(10), big.NewInt(20), big.NewInt(30)}
	prices = []*big.Int{
		big.NewInt(1000000000), // 1 Gwei
		big.NewInt(2000000000), // 2 Gwei
		big.NewInt(3000000000), // 3 Gwei
	}
	ttls = []*big.Int{
		big.NewInt(86400),  // 1 day
		big.NewInt(172800), // 2 days
		big.NewInt(259200), // 3 days
	}
	actives = []bool{true, true, true}
	totalPrice = big.NewInt(0)

	// Deploy AccessTokenSystem
	tokenSystem, err = client.DeployContractFromStrings(
		ctx,
		provider.Signer,
		AccessTokenSystemABI,
		AccessTokenSystemBin,
		"https://example.com/api/token/{id}.json",
	)
	require.NoError(t, err, "Failed to deploy AccessTokenSystem")

	// Create tiers
	for i := 0; i < len(tierIds); i++ {
		receipt, err = tokenSystem.Exec(
			ctx,
			provider.Signer,
			"createTier",
			tierIds[i],
			prices[i],
			ttls[i],
			actives[i],
		)
		require.NoError(t, err, "Failed to create tier")
		totalPrice = totalPrice.Add(totalPrice, prices[i])
	}

	// Batch purchase
	receipt, err = tokenSystem.ExecWithValue(
		ctx,
		consumer.Signer,
		totalPrice,
		"batchPurchaseAccess",
		tierIds,
	)
	require.NoError(t, err, "Failed to batch purchase access")
	assert.NotNil(t, receipt, "Receipt should not be nil")
	assert.Equal(t, uint64(1), receipt.Status, "Receipt status should be 1")

	// Check balances of all tokens
	for _, tierId := range tierIds {
		result, err := tokenSystem.Call(ctx, "balanceOf", consumer.Address(), tierId)
		require.NoError(t, err, "Failed to get token balance")
		balance := result[0].(*big.Int)
		assert.Equal(t, big.NewInt(1), balance, "Consumer should have 1 token for tier "+tierId.String())

		// Verify token is valid
		result, err = tokenSystem.Call(ctx, "isValid", consumer.Address(), tierId)
		require.NoError(t, err, "Failed to check isValid")
		isValid = result[0].(bool)
		assert.True(t, isValid, "Token should be valid for tier "+tierId.String())
	}
}
