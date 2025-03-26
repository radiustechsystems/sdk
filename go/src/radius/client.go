package radius

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	GasEstimateMultiplier = 1.2
	MaxGas                = uint64(1319413953330)
)

type Client struct {
	ChainID *big.Int
	eth     *ethclient.Client
	rpc     *rpc.Client
}

func NewClient(url string) (*Client, error) {
	return NewClientWithHTTPClient(url, &http.Client{
		Transport: http.DefaultTransport,
	})
}

func NewClientWithHTTPClient(url string, httpClient *http.Client) (*Client, error) {
	ctx := context.Background()
	rpcClient, err := rpc.DialOptions(ctx, url, rpc.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}

	ethClient := ethclient.NewClient(rpcClient)
	chainID, err := ethClient.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	return &Client{ChainID: chainID, eth: ethClient, rpc: rpcClient}, nil
}

func NewClientWithLogging(url string, logger Logger) (*Client, error) {
	return NewClientWithHTTPClient(url, &http.Client{
		Transport: InterceptingRoundTripper{
			Proxied: http.DefaultTransport,
			Log:     logger,
		},
	})
}

func (c *Client) AccountFromPrivateKey(key *ecdsa.PrivateKey) (*Account, error) {
	if key == nil {
		return nil, fmt.Errorf("private key is required; use GeneratePrivateKey to create a new key")
	}
	signer := NewPrivateKeySigner(key, c.ChainID)
	return NewAccount(c, signer), nil
}

func (c *Client) API(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	return c.rpc.CallContext(ctx, result, method, args...)
}

func (c *Client) BalanceAt(ctx context.Context, address Address) (*big.Int, error) {
	return c.eth.BalanceAt(ctx, address, nil)
}

func (c *Client) DeployContract(ctx context.Context, signer Signer, abi ABI, bin []byte, args ...interface{}) (*Contract, error) {
	if signer == nil {
		return nil, fmt.Errorf("signer is required for deploying contracts")
	}

	data := bin
	if len(args) > 0 {
		encoded, err := abi.Pack("", args...)
		if err != nil {
			return nil, fmt.Errorf("failed to encode constructor arguments: %w", err)
		}
		data = append(data, encoded...)
	}

	tx, err := c.PrepareTx(ctx, data, signer, nil, big.NewInt(0))
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, fmt.Errorf("failed to prepare transaction")
	}

	receipt, err := c.SendTx(ctx, tx, signer)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy contract: %w", err)
	}
	if receipt == nil {
		return nil, fmt.Errorf("failed to deploy contract: no receipt returned")
	}
	if receipt.Status != 1 {
		return nil, fmt.Errorf("failed to deploy contract: status %d, transaction hash %s", receipt.Status, receipt.TxHash)
	}

	return NewContract(receipt.ContractAddress, abi, c), nil
}

func (c *Client) DeployContractFromStrings(ctx context.Context, signer Signer, abiStr, binStr string, args ...interface{}) (*Contract, error) {
	abi, err := NewABI(abiStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	bin := BytecodeFromHex(binStr)
	if bin == nil {
		return nil, fmt.Errorf("failed to parse bytecode")
	}

	return c.DeployContract(ctx, signer, abi, bin, args...)
}

func (c *Client) CodeAt(ctx context.Context, address Address) ([]byte, error) {
	return c.eth.CodeAt(ctx, address, nil)
}

func (c *Client) EstimateGas(ctx context.Context, tx *Transaction, from Address) (uint64, error) {
	estimate, err := c.eth.EstimateGas(ctx, ethereum.CallMsg{
		From:  from,
		To:    tx.To(),
		Data:  tx.Data(),
		Value: tx.Value(),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to estimate gas: %w", err)
	}

	gas := uint64(float64(estimate) * GasEstimateMultiplier)

	if gas > MaxGas {
		gas = MaxGas
	}

	return gas, nil
}

func (c *Client) Nonce(ctx context.Context, address Address) (uint64, error) {
	return c.eth.PendingNonceAt(ctx, address)
}

func (c *Client) PrepareTx(ctx context.Context, data []byte, signer Signer, to *Address, value *big.Int) (*Transaction, error) {
	var (
		err   error
		gas   uint64
		nonce uint64
	)

	if signer != nil {
		nonce, err = c.Nonce(ctx, signer.Address())
		if err != nil {
			return nil, fmt.Errorf("failed to get nonce: %w", err)
		}
	}

	gas = 0
	gasPrice := big.NewInt(0)
	tx := NewTransaction(data, gas, gasPrice, nonce, to, value)

	if signer == nil {
		return tx, nil
	}

	gas, err = c.EstimateGas(ctx, tx, signer.Address())
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas: %w", err)
	}

	return NewTransaction(data, gas, big.NewInt(0), nonce, to, value), nil
}

func (c *Client) Send(ctx context.Context, signer Signer, to Address, value *big.Int) (*Receipt, error) {
	if signer == nil {
		return nil, fmt.Errorf("signer is required")
	}

	tx, err := c.PrepareTx(ctx, nil, signer, &to, value)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, fmt.Errorf("failed to prepare transaction")
	}

	receipt, err := c.SendTx(ctx, tx, signer)

	return receipt, err
}

func (c *Client) SendTx(ctx context.Context, tx *Transaction, signer Signer) (*Receipt, error) {
	if signer == nil {
		return nil, fmt.Errorf("signer is required")
	}

	stx, err := signer.SignTx(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}
	if stx == nil {
		return nil, fmt.Errorf("failed to sign transaction")
	}

	return c.SendSignedTx(ctx, stx)
}

func (c *Client) SendSignedTx(ctx context.Context, tx *Transaction) (*Receipt, error) {
	if err := c.eth.SendTransaction(ctx, tx); err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	receipt, err := bind.WaitMinedHash(ctx, c.eth, tx.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction receipt: %w", err)
	}
	if receipt == nil {
		return nil, fmt.Errorf("failed to get transaction receipt: no receipt returned")
	}
	if receipt.Status != 1 {
		return receipt, fmt.Errorf("failed to execute transaction: status %d, transaction hash %s", receipt.Status, receipt.TxHash)
	}

	return receipt, nil
}
