# Contributing to the Radius Go SDK

This guide provides specific guidelines for contributing to the Go SDK.

## Development Setup

1. Install Go 1.23.3 or later
2. Clone the repository
3. Run `go mod download`
4. Make sure tests pass: `go test ./...`

## Go-Specific Patterns

### Provider Abstraction

The SDK uses an abstraction layer for Ethereum-compatible libraries that provides data structure mappings. This is managed through the `providers/eth` package:

```go
// ✅ Correct
import "github.com/radiustechsystems/sdk/go/radius/providers/eth"

// ❌ Incorrect
import "github.com/ethereum/go-ethereum"
import "github.com/ethereum/go-ethereum/common"
```

While the SDK has its own concrete implementations of core data structures (like addresses and transactions), the provider package maps these to Ethereum types when needed to leverage the functionality provided by established libraries.

This abstraction:

- Maintains independent implementations optimized for Radius
- Leverages battle-tested Ethereum libraries when beneficial
- Presents a consistent interface regardless of the underlying implementation
- Simplifies testing and mocking
- Provides a single point for managing external dependencies

### Options Pattern

We use the functional options pattern for optional configuration, but not for required parameters:

```go
// Good use of options pattern - for optional configuration
func NewClient(url string, opts ...Option) (*Client, error) {
    options := defaultOptions()
    for _, opt := range opts {
        opt(options)
    }
    // Use options to configure the client
}

// Not using options for required parameters
func NewContract(address Address, abi *ABI) *Contract {
    return &Contract{
        address: address,
        abi:     abi,
    }
}
```

When to use the options pattern:
- For truly optional configuration (logging, interceptors, timeouts)
- When extending functionality without changing method signatures
- When configurations may grow over time

When NOT to use the options pattern:
- For required parameters (addresses, ABIs, bytecode)
- When parameters are essential to the object's function
- When clarity of required inputs is important

Key benefits:
- Clear distinction between required and optional parameters
- Self-documenting parameter names with the `With*` prefix
- Extensibility without breaking changes
- Graceful defaults for optional parameters

### Interface Design

Interfaces are defined where they are used, not implemented:

```go
type SomeInterface interface {
    SomeFunc(ctx context.Context) error
}
```

Our SDK follows these interface design principles:

- We use the `Signer` interface for authentication abstraction, allowing multiple signing methods (private key, Clef, etc.)
- We use client interfaces (like `AccountClient`) to define specific capabilities needed by each package
- We prefer concrete types (`Account`, `Address`, `Contract`) for core data structures that have properties users need to access directly
- Each package-specific client interface contains only the methods needed for that package's functionality
- Public interfaces are well-documented with clear purpose statements

### Error Handling

- Use explicit error types when beneficial
- Return errors rather than panic
- Include context in error messages
- Wrap errors with `fmt.Errorf`: `return fmt.Errorf("failed to deploy contract: %w", err)`

### Type Conversions

Use explicit conversion methods between types without the `To` prefix:
```go
// Preferred conversion method naming
func (a *Address) EthAddress() eth.Address {
    return eth.BytesToAddress(a.Bytes())
}

// Avoid using "To" prefix in method names
// func (a *Address) ToEthAddress() eth.Address { ... } // ❌ Not preferred
```

This naming convention:
- Is more concise and reads more naturally
- Follows common Go conventions (e.g., `String()` not `ToString()`)
- Makes generated code more readable with less redundancy

## Testing

### Test Organization

- Place test helpers and shared test utilities in a dedicated `test` directory
- We maintain this structure for consistency across all SDK implementations (Go, TypeScript, etc.)
- This approach helps to avoid import cycles and keeps the test code well-organized

### Unit Tests

- Place tests in `_test.go` files
- Use table-driven tests when appropriate
- Use testify for assertions
- Mock interfaces for isolation

Example:
```go
func TestMyClass_SomeFunc(t *testing.T) {
    tests := []struct {
        name    string
        foo     string
        bar     int
        want    string
        wantErr bool
    }{
        {
            name: "simple case",
            foo: "test",
            bar: 42,
            want: "expected result",
        },
        {
            name: "error case",
            foo: "invalid",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mc := NewMyClass(tt.foo, tt.bar)
            got, err := mc.SomeFunc()
            if (err != nil) != tt.wantErr {
                t.Errorf("MyClass.SomeFunc() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("MyClass.SomeFunc() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Integration Tests

- Use build tags: `//go:build integration`
- Require environment variables
- Clean up resources
- Handle timeouts
- Skip tests gracefully when required environment variables are not set

## Style Guidelines

### Code Organization

- Keep files focused and small
- Group related types and functions
- Use clear package names
- Follow standard Go layout
- Maintain our repository structure:
  - `radius/`: Public API package that users import
  - `src/`: Implementation details
  - `test/`: Test utilities and integration tests

### Naming Conventions

- Use descriptive names
- Follow Go conventions:
    - PascalCase for exported names
    - camelCase for internal names
    - Short names for short scopes

### Comments

- Begin comments with the name:
  ```go
  // Account represents a Radius account.
  type Account struct {
  ```
- Document all exported items
- Include examples for complex functions

## Quality Checks

Before submitting:

1. Run tests:
   ```bash
   go test -v -race -coverprofile=coverage.txt ./...
   ```
Note: On macOS, you may see a linker warning about `malformed LC_DYSYMTAB`. This is a known toolchain issue that
doesn't affect test execution or results and can be safely ignored.

2. Run linter:
   ```bash
   golangci-lint run
   ```

3. Format code:
   ```bash
   go fmt ./...
   ```

For convenience, you can fix all lint and formatting errors by running the cleanup script:

- On Windows: `cleanup.bat`
- On MacOS/Linux: `./cleanup.sh`

## Dependencies

- Keep dependencies minimal
- Use standard library when possible
- Pin versions in go.mod
- Document why each dependency is needed

## Common Patterns

### Unified Public API

All public SDK functionality should be exposed through the main `radius` package:

```go
// External usage should only import the main radius package
import "github.com/radiustechsystems/sdk/go/radius"

// Inside the SDK, use the appropriate package from src
import "github.com/radiustechsystems/sdk/go/src/common"
```

The `sdk.go` file in the `radius` package re-exports all public types and functions, providing a clean, single-import API for users while maintaining internal organization.

This pattern aligns with our folder structure:
- Implementation details live in the `src/` directory
- Public API is defined in the `radius/` directory
- Test utilities are in the `test/` directory
- Users should never import directly from implementation modules

### Context Usage

- Accept context.Context as first parameter
- Pass through to underlying calls
- Respect cancellation

### Builders and Factories

Use New* functions for construction:
```go
type MyClass struct {
    foo string
    bar int
}

func NewMyClass(foo string, bar int, opts ...Option) (*MyClass, error) {
    options := &Options{}
    for _, opt := range opts {
        opt(options)
    }
    
    return &MyClass{
        foo: foo,
        bar: bar,
    }, nil
}
```

### Error Handling Patterns

Wrap errors with context:
```go
func (m *MyClass) SomeFunc() error {
    result, err := someOperation()
    if err != nil {
        return fmt.Errorf("failed to perform operation: %w", err)
    }
    return nil
}
```

## Documentation

- Update README.md with new features
- Include godoc examples
- Keep CHANGELOG.md current
- Document breaking changes
