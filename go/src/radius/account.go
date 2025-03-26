package radius

import (
	"context"
	"fmt"
	"math/big"
)

type Account struct {
	Client *Client
	Signer Signer
}

func NewAccount(client *Client, signer Signer) *Account {
	return &Account{Client: client, Signer: signer}
}

func (a *Account) Address() Address {
	if a.Signer == nil {
		return Address{}
	}

	return a.Signer.Address()
}

func (a *Account) Balance(ctx context.Context) (*big.Int, error) {
	if a.Client == nil {
		return nil, fmt.Errorf("radius client is required for account calls")
	}

	return a.Client.BalanceAt(ctx, a.Address())
}
