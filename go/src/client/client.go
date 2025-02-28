// Package client provides the primary interface for interacting with the Radius platform.
// It implements methods for account management, contract deployment, transaction handling,
// and querying Radius state.
package client

import (
	"context"
	"fmt"
	"math/big"
	"net/http"

	"github.com/radiustechsystems/sdk/go/src/auth"
	"github.com/radiustechsystems/sdk/go/src/common"
	"github.com/radiustechsystems/sdk/go/src/contracts"
	"github.com/radiustechsystems/sdk/go/src/providers/eth"
	"github.com/radiustechsystems/sdk/go/src/transport"
)

// Client is used to interact with the Radius platform.
// It serves as the main entry point for working with the Radius ecosystem.
// It provides methods for account management, contract deployment and interaction,
// transaction handling, and querying Radius state.
type Client struct {
	// httpClient is the HTTP client used for making API requests
	httpClient *http.Client

	// ethClient is the Ethereum client used to communicate with Radius
	ethClient *eth.Client
}

// New creates a new Radius Client with the given URL and ClientOption(s).
//
// @param url URL of the Radius node
// @param opts Optional client configuration options
// @return New Radius Client instance and nil error on success
// @return nil and error if client creation fails
func New(url string, opts ...Option) (*Client, error) {
	options := &Options{
		httpClient: &http.Client{},
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.httpClient.Transport == nil {
		options.httpClient.Transport = http.DefaultTransport
	}

	if options.logger != nil || options.interceptor != nil {
		irt := transport.InterceptingRoundTripper{
			Proxied:     options.httpClient.Transport,
			Interceptor: options.interceptor,
			Logf:        options.logger,
		}
		options.httpClient.Transport = irt
	}

	ethClient, err := eth.NewClient(url, options.httpClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create Radius client: %w", err)
	}

	return &Client{
		httpClient: options.httpClient,
		ethClient:  ethClient,
	}, nil
}

// BalanceAt returns the balance of the given address in wei.
//
// @param ctx Context for the request
// @param address Address to check the balance for
// @return Balance in wei and nil error on success
// @return nil and error if the balance cannot be retrieved from the network
func (c *Client) BalanceAt(ctx context.Context, address common.Address) (*big.Int, error) {
	balance, err := c.ethClient.BalanceAt(ctx, address.EthAddress(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}
	return balance, nil
}

// Call executes a contract method call and returns the decoded result. This is used for read-only contract methods,
// and does not require a transaction to be sent to Radius. Alternatively, you can use the contracts.Contract method
// Call, which provides a more convenient interface for interacting with smart contracts.
func (c *Client) Call(ctx context.Context, contract *contracts.Contract, method string, args ...interface{}) ([]interface{}, error) {
	if contract.ABI == nil {
		return nil, fmt.Errorf("contract ABI is required")
	}

	address := contract.Address()
	if address.Equals(common.ZeroAddress()) {
		return nil, fmt.Errorf("contract address is required")
	}

	data, err := contract.ABI.Pack(method, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to encode method call: %w", err)
	}

	params := txParams{
		to:    &address,
		data:  data,
		value: big.NewInt(0),
	}

	tx, err := c.prepareTx(ctx, params)
	if err != nil {
		return nil, err
	}

	result, err := c.ethClient.CallContract(ctx, eth.CallMsg{
		To:    common.EthAddressFromRadiusAddress(tx.To),
		Data:  tx.Data,
		Value: tx.Value,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("contract call failed: %w", err)
	}

	decoded, err := contract.ABI.Unpack(method, result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode result: %w", err)
	}

	return decoded, nil
}

// ChainID returns the chain ID of the connected Radius network.
//
// @param ctx Context for the request
// @return Chain ID of the network and nil error on success
// @return nil and error if the chain ID cannot be retrieved from the network
func (c *Client) ChainID(ctx context.Context) (*big.Int, error) {
	chainID, err := c.ethClient.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	return chainID, nil
}

// CodeAt returns the contract code at the given address.
//
// @param ctx Context for the request
// @param address Address of the contract to retrieve code for
// @return Contract bytecode and nil error on success
// @return nil and error if the code cannot be retrieved from the network
func (c *Client) CodeAt(ctx context.Context, address common.Address) ([]byte, error) {
	code, err := c.ethClient.CodeAt(ctx, address.EthAddress(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get code: %w", err)
	}
	return code, nil
}

// DeployContract deploys the given EVM smart contract bytecode to Radius. If the contract has a constructor, the
// ABI and constructor arguments must be provided.
func (c *Client) DeployContract(ctx context.Context, signer auth.Signer, bytecode []byte, abi *common.ABI, args ...interface{}) (*contracts.Contract, error) {
	if signer == nil {
		return nil, fmt.Errorf("signer is required for deploying contracts")
	}

	data := bytecode
	if len(args) > 0 && abi != nil {
		encodedConstructorArgs, err := abi.Pack("", args...)
		if err != nil {
			return nil, fmt.Errorf("failed to encode constructor arguments: %w", err)
		}
		data = append(data, encodedConstructorArgs...)
	}

	receipt, err := c.prepareAndSendTx(ctx, txParams{
		data:   data,
		signer: signer,
		value:  big.NewInt(0),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to deploy contract: %w", err)
	}
	if receipt == nil {
		return nil, fmt.Errorf("failed to deploy contract: no receipt returned")
	}
	if receipt.Status != 1 {
		return nil, fmt.Errorf("failed to deploy contract: status %d, transaction hash %s", receipt.Status, receipt.TxHash)
	}

	return contracts.New(receipt.ContractAddress, abi), nil
}

// EstimateGas estimates the gas cost of the given transaction. This is handled automatically by the Execute, Send,
// and Transact methods, so you only need to call this method if you need to get the gas cost manually.
func (c *Client) EstimateGas(ctx context.Context, tx *common.Transaction) (uint64, error) {
	estimate, err := c.ethClient.EstimateGas(ctx, eth.CallMsg{
		To:    common.EthAddressFromRadiusAddress(tx.To),
		Data:  tx.Data,
		Value: tx.Value,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to estimate gas: %w", err)
	}

	// Apply safety margin of 20% to the estimated gas cost
	margin := estimate / 5
	gas := estimate + margin

	// Limit gas to maxGas
	if gas > common.MaxGas {
		gas = common.MaxGas
	}

	return gas, nil
}

// Execute executes a contract method call and returns the transaction receipt. This is used for state-changing contract
// methods, and requires a transaction to be sent to Radius. A more convenient interface for interacting with smart
// contracts is provided by the contracts.Contract method Execute.
func (c *Client) Execute(ctx context.Context, contract *contracts.Contract, signer auth.Signer, method string, args ...interface{}) (*common.Receipt, error) {
	if contract.ABI == nil {
		return nil, fmt.Errorf("contract ABI is required")
	}

	address := contract.Address()
	if address.Equals(common.ZeroAddress()) {
		return nil, fmt.Errorf("contract address is required")
	}

	data, err := contract.ABI.Pack(method, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to encode method call: %w", err)
	}

	return c.prepareAndSendTx(ctx, txParams{
		to:     &address,
		data:   data,
		signer: signer,
		value:  big.NewInt(0),
	})
}

// HTTPClient returns the underlying HTTP client used by the Radius Client.
func (c *Client) HTTPClient() *http.Client {
	return c.httpClient
}

// PendingNonceAt returns the pending nonce of the given address. In most cases, you should not need to call this
// method directly.
func (c *Client) PendingNonceAt(ctx context.Context, address common.Address) (uint64, error) {
	nonce, err := c.ethClient.PendingNonceAt(ctx, address.EthAddress())
	if err != nil {
		return 0, fmt.Errorf("failed to get nonce: %w", err)
	}
	return nonce, nil
}

// Send sends value to the recipient address, and returns the Radius transaction Receipt.
func (c *Client) Send(
	ctx context.Context,
	signer auth.Signer,
	recipient common.Address,
	value *big.Int,
) (*common.Receipt, error) {
	if signer == nil {
		return nil, fmt.Errorf("signer is required for sending transactions")
	}

	receipt, err := c.prepareAndSendTx(ctx, txParams{
		signer: signer,
		to:     &recipient,
		value:  value,
	})
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}
	if receipt == nil {
		return nil, fmt.Errorf("transaction failed: no receipt returned")
	}
	if receipt.Status != 1 {
		return nil, fmt.Errorf("transaction failed: %v", receipt)
	}

	return receipt, nil
}

// Transact sends a signed transaction to the Radius platform, and returns the Radius transaction Receipt.
func (c *Client) Transact(
	ctx context.Context,
	signer auth.Signer,
	tx *common.SignedTransaction,
) (*common.Receipt, error) {
	if signer == nil {
		return nil, fmt.Errorf("signer is required for sending transactions")
	}

	if tx == nil {
		return nil, fmt.Errorf("no signed transaction provided")
	}

	ethTx := tx.EthSignedTransaction()

	if err := c.ethClient.SendTransaction(ctx, ethTx); err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	receipt, err := eth.WaitMined(ctx, c.ethClient, ethTx)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction receipt: %w", err)
	}
	if receipt == nil {
		return nil, fmt.Errorf("failed to get transaction receipt: no receipt returned")
	}
	if receipt.Status != 1 {
		return nil, fmt.Errorf("transaction failed: status %d, transaction hash %s", receipt.Status, receipt.TxHash)
	}

	from := signer.Address()
	to := common.ZeroAddress()
	if tx.To != nil {
		to = *tx.To
	}
	value := tx.Value

	return common.ReceiptFromEthReceipt(receipt, from, to, value), nil
}

// prepareTx prepares a Radius transaction, ensuring that the nonce is set correctly. In most cases, you should use the
// Execute or Send methods instead, which provide a more convenient interface.
func (c *Client) prepareTx(ctx context.Context, params txParams) (*common.Transaction, error) {
	var (
		err   error
		nonce uint64
	)

	// Get the pending nonce for the signer address, if necessary
	if params.signer != nil {
		nonce, err = c.PendingNonceAt(ctx, params.signer.Address())
		if err != nil {
			return nil, fmt.Errorf("failed to get nonce: %w", err)
		}
	}

	// Must set Transaction.To value to nil if it is the zero address
	to := params.to
	if params.to == nil || params.to.Equals(common.ZeroAddress()) {
		to = nil
	}

	// Create the initial transaction used to estimate gas
	tx := &common.Transaction{
		Data:     params.data,
		Nonce:    nonce,
		Gas:      0,
		GasPrice: big.NewInt(0),
		To:       to,
		Value:    params.value,
	}

	// Estimate gas cost for the transaction
	tx.Gas, err = c.EstimateGas(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas: %w", err)
	}

	return tx, nil
}

// prepareAndSendTx prepares and sends a Radius transaction, ensuring that the transaction is signed correctly. In
// most cases, you should use the Execute or Send methods instead, which provide a more convenient interface.
func (c *Client) prepareAndSendTx(ctx context.Context, params txParams) (*common.Receipt, error) {
	if params.signer == nil {
		return nil, fmt.Errorf("signer is required for sending transactions")
	}

	tx, err := c.prepareTx(ctx, params)
	if err != nil {
		return nil, err
	}

	signedTx, err := params.signer.SignTransaction(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	return c.Transact(ctx, params.signer, signedTx)
}

// txParams contains the parameters required to prepare and send a Radius transaction.
// This is an internal struct used by the Client for transaction preparation.
type txParams struct {
	// data is the transaction data (bytecode for contract creation or method call data)
	data []byte

	// signer is used to sign the transaction
	signer auth.Signer

	// to is the destination address for the transaction (nil for contract creation)
	to *common.Address

	// value is the amount of native currency to send with the transaction
	value *big.Int
}
