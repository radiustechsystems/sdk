# Radius Go SDK

The official Go client library for interacting with the [Radius platform](https://radiustech.xyz/), providing a simple
and idiomatic way to interact with Radius services.

## Features

- Account management and transaction signing
- Smart contract deployment and interaction
- Optional request logging and interceptors
- EVM compatibility with high performance & low latency

## Requirements

- Go 1.13 or higher
- Radius JSON-RPC endpoint: https://docs.radiustech.xyz/radius-testnet-access
- Ethereum private key: https://ethereum.org/en/developers/docs/accounts/#account-creation

## Installation

```bash
go get github.com/radiustechsystems/sdk/go
```

## Quickstart Examples

### Connect to Radius

Be sure to use your own `RADIUS_ENDPOINT` and `PRIVATE_KEY` values, as mentioned in the [Requirements](#requirements).

```go
const (
	RADIUS_ENDPOINT = "https://rpc.testnet.tryradi.us/xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx";
	PRIVATE_KEY = "fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd036415f";
)

client, err := radius.NewClient(RADIUS_ENDPOINT)
account := radius.NewAccount(radius.WithPrivateKeyHex(PRIVATE_KEY, client))
```

### Transfer Value Between Accounts

Here, we send 100 tokens to another account. Be sure to replace the recipient's address with one of your own.

```go
// Send 100 tokens to another account
recipient, err := radius.AddressFromHex("0x5e97870f263700f46aa00d967821199b9bc5a120") // Recipient's address
amount := big.NewInt(100)
receipt, err := account.Send(context.Background(), client, recipient, amount)
if err != nil {
	log.Fatal(err)
}

log.Printf("Transaction hash: %s", receipt.TxHash.Hex())
```

### Deploy a Smart Contract

Here, we deploy the [SimpleStorage.sol](https://github.com/radiustechsystems/sdk/tree/main/contracts/solidity)
example contract included in this SDK, with the application binary interface (ABI) and bytecode that were generated
using the Solidity compiler [solcjs](https://docs.soliditylang.org/en/latest/installing-solidity.html#npm-node-js).

```go
// Parse ABI and bytecode of the SimpleStorage contract
abi := radius.ABIFromJSON(`[{"inputs":[],"name":"get","outputs":[{"type":"uint256"}],"type":"function"},{"inputs":[{"type":"uint256"}],"name":"set","type":"function"}]`)
bytecode := radius.BytecodeFromHex("6080604052348015600e575f5ffd5b5060a580601a5f395ff3fe6080604052348015600e575f5ffd5b50600436106030575f3560e01c806360fe47b11460345780636d4ce63c146045575b5f5ffd5b6043603f3660046059565b5f55565b005b5f5460405190815260200160405180910390f35b5f602082840312156068575f5ffd5b503591905056fea26469706673582212207655d86666fa8aa75666db8416e0f5db680914358a57e84aa369d9250218247f64736f6c634300081c0033")

// Deploy the contract
contract, err := client.DeployContract(
	context.Background(),
	account.Signer(),
	bytecode,
	abi,
)
```

### Interact with a Smart Contract

Assuming the contract was previously deployed (which is typically the case), we can interact with it using the contract
address and ABI. Be sure to replace the contract address with that of your own deployed contract.

```go
// Reference a previously deployed contract
address, err := radius.AddressFromHex("0x5e97870f263700f46aa00d967821199b9bc5a120") // Contract address
abi := radius.ABIFromJSON(`[{"inputs":[],"name":"get","outputs":[{"type":"uint256"}],"type":"function"},{"inputs":[{"type":"uint256"}],"name":"set","type":"function"}]`)
contract := radius.NewContract(address, abi)

// Write to the contract
value := big.NewInt(42)
receipt, err := contract.Execute(context.Background(), client, account.Signer(), "set", value)

// Read from the contract
result, err := contract.Call(context.Background(), client, "get")
log.Printf("Stored value: %v", result[0])
```

## Advanced Features

### Custom Transaction Signing

```go
type MyCustomSigner struct{
	// ...
}
func (s *MyCustomSigner) Address() radius.Address { /* ... */ }
func (s *MyCustomSigner) ChainID() *big.Int { /* ... */ }
func (s *MyCustomSigner) Hash(tx *radius.Transaction) radius.Hash { /* ... */ }
func (s *MyCustomSigner) SignMessage(message []byte) ([]byte, error) { /* ... */ }
func (s *MyCustomSigner) SignTransaction(tx *radius.Transaction) (*radius.SignedTransaction, error) { /* ... */ }

customSigner := &MyCustomSigner{}
customSignerAccount := radius.NewAccount(radius.WithSigner(customSigner))
```

### Logging and Request Interceptors

```go
client, err := radius.NewClient("https://your-radius-endpoint",
	radius.WithLogger(func(format string, args ...any) {
		log.Printf(format, args...)
	}),
	radius.WithInterceptor(func(reqBody string, resp *http.Response) (*http.Response, error) {
		// Examine request body, modify response, etc.
		return resp, nil
	}),
)
```

### Custom HTTP Client

```go
httpClient := &http.Client{
	Timeout: time.Second * 30,
}

client, err := radius.NewClient("https://your-radius-endpoint",
	radius.WithHttpClient(httpClient),
)
```

## Resources

- [Website](https://radiustech.xyz/)
- [Testnet Access](https://docs.radiustech.xyz/radius-testnet-access)
- [GitHub Issues](https://github.com/radiustechsystems/sdks/issues)
- [Changelog](CHANGELOG.md)

## Contributing

Please see the [Go SDK Contributing Guide](CONTRIBUTING.md) for detailed information about contributing to this SDK.
For repository-wide guidelines, see the [General Contributing Guide](../CONTRIBUTING.md).

## License

All Radius SDKs are released under the [MIT License](../LICENSE).
