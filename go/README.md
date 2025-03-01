# Radius Go SDK

A Go client library for interacting with the [Radius platform](https://radiustech.xyz/), providing a simple and
idiomatic way to interact with Radius services.

## Features

- Account management and transaction signing
- Smart contract deployment and interaction
- Optional request logging and interceptors
- EVM compatibility with high performance
- Comprehensive error handling

## Requirements

- Go 1.23.3 or higher

## Installation

```bash
go get github.com/radiustechsystems/sdk/go
go mod tidy
```

## Quickstart Examples

### Connect to Radius and Create an Account

```go
// Connect to Radius
client, err := radius.NewClient("https://your-radius-endpoint")
if err != nil {
    log.Fatal(err)
}

// Create an account using a private key
account := radius.NewAccount(
    radius.WithPrivateKeyHex("your-private-key", client),
)
```

### Transfer Value Between Accounts

```go
// Send 100 tokens to another account
recipient, err := radius.AddressFromHex("0x...")
amount := big.NewInt(100)
receipt, err := account.Send(context.Background(), client, recipient, amount)
if err != nil {
    log.Fatal(err)
}

log.Printf("Transaction hash: %s", receipt.TxHash.Hex())
```

### Deploy a Smart Contract

```go
// Parse ABI and bytecode
abi := radius.ABIFromJSON(`[{"inputs":[],"name":"get","outputs":[{"type":"uint256"}],"type":"function"},{"inputs":[{"type":"uint256"}],"name":"set","type":"function"}]`)
bytecode := radius.BytecodeFromHex("608060405234801561001057600080fd5b50610150806100...")

// Deploy the contract
contract, err := client.DeployContract(
    context.Background(),
    account.Signer(),
    bytecode,
    abi,
)
```

### Interact with a Smart Contract

```go
// Reference an existing contract
address, err := radius.AddressFromHex("0x...")
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

This project is licensed under the [MIT License](../LICENSE).
