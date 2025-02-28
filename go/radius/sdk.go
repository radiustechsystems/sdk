package radius

import (
	"crypto/ecdsa"
	"net/http"

	"github.com/radiustechsystems/sdk/go/src/accounts"
	"github.com/radiustechsystems/sdk/go/src/auth"
	"github.com/radiustechsystems/sdk/go/src/auth/clef"
	"github.com/radiustechsystems/sdk/go/src/auth/privatekey"
	"github.com/radiustechsystems/sdk/go/src/client"
	"github.com/radiustechsystems/sdk/go/src/common"
	"github.com/radiustechsystems/sdk/go/src/contracts"
	"github.com/radiustechsystems/sdk/go/src/transport"
)

const MaxGas = common.MaxGas

type (
	ABI               = common.ABI
	Account           = accounts.Account
	AccountClient     = accounts.AccountClient
	AccountOption     = accounts.Option
	Address           = common.Address
	AuthClient        = auth.SignerClient
	ClefSigner        = clef.Signer
	Client            = client.Client
	ClientOption      = client.Option
	Contract          = contracts.Contract
	Event             = common.Event
	Hash              = common.Hash
	Interceptor       = transport.Interceptor
	KeySigner         = privatekey.Signer
	Logf              = transport.Logf
	Receipt           = common.Receipt
	Signer            = auth.Signer
	SignedTransaction = common.SignedTransaction
	Transaction       = common.Transaction
)

// ABIFromJSON creates a new ABI with the given JSON string. If the JSON is invalid, it returns nil.
func ABIFromJSON(json string) *ABI {
	return common.ABIFromJSON(json)
}

// AddressFromHex creates an Address from a hex string. If the hex string is invalid, it returns an error.
func AddressFromHex(h string) (Address, error) {
	return common.AddressFromHex(h)
}

// BytecodeFromHex converts a hex string to a byte slice. If the string is not a valid hex, it returns nil.
func BytecodeFromHex(s string) []byte {
	return common.BytecodeFromHex(s)
}

// NewABI creates a new ABI with the given JSON string.
func NewABI(abiJSON string) (*ABI, error) {
	return common.NewABI(abiJSON)
}

// NewAccount creates a new Radius Account with the given options.
func NewAccount(opts ...AccountOption) *Account {
	return accounts.New(opts...)
}

// NewAddress creates a new Radius Address with the given byte slice.
func NewAddress(b []byte) common.Address {
	return common.NewAddress(b)
}

// NewClefSigner creates a new ClefSigner with the given Address, Radius Client, and Clef URL.
func NewClefSigner(address common.Address, client AuthClient, clefURL string) (*ClefSigner, error) {
	return clef.New(address, client, clefURL)
}

// NewClient creates a new Radius Client with the given URL and options.
func NewClient(url string, opts ...ClientOption) (*Client, error) {
	return client.New(url, opts...)
}

// NewContract creates a new Radius Contract with the given options.
func NewContract(address Address, abi *ABI) *Contract {
	return contracts.New(address, abi)
}

// NewKeySigner creates a new KeySigner with the given private key and Radius Client.
func NewKeySigner(key *ecdsa.PrivateKey, client AuthClient) Signer {
	return privatekey.New(key, client)
}

// WithHTTPClient returns a ContractOption that sets the Radius chain ID for the contract.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return client.WithHTTPClient(httpClient)
}

// WithInterceptor returns a ClientOption that adds a request/response Interceptor to a Radius Client.
func WithInterceptor(interceptor Interceptor) ClientOption {
	return client.WithInterceptor(interceptor)
}

// WithLogger returns a ClientOption that adds request/response logging to a Radius Client.
func WithLogger(logger Logf) ClientOption {
	return client.WithLogger(logger)
}

// WithPrivateKey returns an AccountOption that adds a KeySigner and Address to an Account using a private key.
func WithPrivateKey(key *ecdsa.PrivateKey, client AccountClient) AccountOption {
	return accounts.WithPrivateKey(key, client)
}

// WithPrivateKeyHex returns an AccountOption that adds a KeySigner and Address to an Account using a private key
// hex string.
func WithPrivateKeyHex(key string, client AccountClient) AccountOption {
	return accounts.WithPrivateKeyHex(key, client)
}

// WithSigner returns an AccountOption that adds a Signer to an Account. The Signer is used to derive the Address of
// the Account and sign transactions. The Signer must implement the Signer interface (e.g. ClefSigner, KeySigner).
func WithSigner(signer Signer) AccountOption {
	return accounts.WithSigner(signer)
}

// ZeroAddress returns the zero address.
func ZeroAddress() Address {
	return common.ZeroAddress()
}
