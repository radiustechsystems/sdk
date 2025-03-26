package radius

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
)

type Contract struct {
	address Address
	ABI     ABI
	Client  *Client
	code    []byte
}

func NewContract(address Address, abi ABI, client *Client) *Contract {
	return &Contract{ABI: abi, address: address, Client: client}
}

func (c *Contract) Address() *Address {
	return &c.address
}

func (c *Contract) Call(ctx context.Context, method string, args ...interface{}) ([]interface{}, error) {
	if c.Client == nil {
		return nil, fmt.Errorf("radius client is required for contract calls")
	}

	params, err := c.ABI.Pack(method, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to encode method call: %w", err)
	}

	tx, err := c.Client.PrepareTx(ctx, params, nil, c.Address(), big.NewInt(0))
	if err != nil {
		return nil, fmt.Errorf("failed to prepare transaction: %w", err)
	}

	data, err := c.Client.eth.CallContract(ctx, ethereum.CallMsg{
		To:    tx.To(),
		Data:  tx.Data(),
		Value: tx.Value(),
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("contract call failed: %w", err)
	}

	result, err := c.ABI.Unpack(method, data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode result: %w", err)
	}

	return result, nil
}

func (c *Contract) Code(ctx context.Context) ([]byte, error) {
	var err error

	if c.code != nil {
		return c.code, nil
	}
	if c.Client == nil {
		return nil, fmt.Errorf("radius client is required to fetch contract code")
	}

	c.code, err = c.Client.eth.CodeAt(ctx, *c.Address(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract code: %w", err)
	}

	return c.code, nil
}

func (c *Contract) Exec(ctx context.Context, signer Signer, method string, args ...interface{}) (*Receipt, error) {
	return c.ExecWithValue(ctx, signer, big.NewInt(0), method, args...)
}

func (c *Contract) ExecWithValue(ctx context.Context, signer Signer, value *big.Int, method string, args ...interface{}) (*Receipt, error) {
	if c.Client == nil {
		return nil, fmt.Errorf("radius client is required for sending transactions")
	}

	if signer == nil {
		return nil, fmt.Errorf("signer is required for sending transactions")
	}

	data, err := c.ABI.Pack(method, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to encode method call: %w", err)
	}

	tx, err := c.Client.PrepareTx(ctx, data, signer, c.Address(), value)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare transaction: %w", err)
	}

	return c.Client.SendTx(ctx, tx, signer)
}
