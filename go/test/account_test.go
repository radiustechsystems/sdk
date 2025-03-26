package test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/radiustechsystems/sdk/go/src/radius"
)

func TestNewAccount(t *testing.T) {
	client := CreateTestClient(t, nil)
	signer := CreateTestSigner()

	account := radius.NewAccount(client, signer)
	assert.Equal(t, client, account.Client, "Client should match")
	assert.Equal(t, signer, account.Signer, "Signer should match")
}

func TestAccount_Address(t *testing.T) {
	signer := CreateTestSigner()

	tests := []struct {
		name         string
		signer       radius.Signer
		expectedAddr radius.Address
		expected     radius.Address
	}{
		{
			name:     "With nil signer",
			signer:   nil,
			expected: radius.Address{},
		},
		{
			name:         "With valid signer",
			signer:       signer,
			expectedAddr: signer.Address(),
		},
		{
			name:     "With failing signer",
			signer:   &radius.PrivateKeySigner{},
			expected: radius.Address{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			account := radius.NewAccount(nil, tc.signer)
			addr := account.Address()
			assert.Equal(t, tc.expectedAddr, addr, "Address should match")
		})
	}
}

func TestAccount_Balance(t *testing.T) {
	ctx := context.Background()
	signer := CreateTestSigner()

	tests := []struct {
		name          string
		client        *radius.Client
		signer        radius.Signer
		expectedError bool
		expectedValue *big.Int
	}{
		{
			name:          "With nil client",
			client:        nil,
			signer:        signer,
			expectedError: true,
		},
		{
			name:          "Client returns error",
			client:        CreateTestClient(t, nil),
			signer:        signer,
			expectedError: true,
		},
		{
			name: "Success case",
			client: CreateTestClient(t, map[string]func(params []interface{}) interface{}{
				"eth_getBalance": func(_ []interface{}) interface{} {
					return fmt.Sprintf("0x%x", OneETH)
				},
			}),
			signer:        signer,
			expectedValue: OneETH,
		},
		{
			name: "Zero balance",
			client: CreateTestClient(t, map[string]func(params []interface{}) interface{}{
				"eth_getBalance": func(_ []interface{}) interface{} {
					return fmt.Sprintf("0x%x", big.NewInt(0))
				},
			}),
			signer:        signer,
			expectedValue: big.NewInt(0),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			account := radius.NewAccount(tc.client, tc.signer)
			balance, err := account.Balance(ctx)

			if tc.expectedError {
				assert.Error(t, err, "Balance should return an error")
			}

			if !tc.expectedError {
				assert.NoError(t, err, "Balance should not return an error")
				assert.Equal(t, 0, tc.expectedValue.Cmp(balance), "Balance should match")
			}
		})
	}
}
