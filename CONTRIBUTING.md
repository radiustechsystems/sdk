# Contributing to Radius SDKs

Thank you for your interest in contributing to the Radius SDKs! This document provides guidelines and instructions for
contributing.

## Table of Contents
- [Code of Conduct](#code-of-conduct)
- [Repository Structure](#repository-structure)
- [Core Design Principles](#core-design-principles)
- [Development Workflow](#development-workflow)
- [Commit Messages](#commit-messages)
- [Pull Requests](#pull-requests)
- [Documentation](#documentation)
- [Language-Specific Guidelines](#language-specific-guidelines)
- [Questions](#questions)

## Code of Conduct

This project and everyone participating in it is governed by our [Code of Conduct](CODE_OF_CONDUCT.md). By
participating, you are expected to uphold this code. Please report unacceptable behavior to: 
[opensource@radiustech.xyz](mailto:opensource@radiustech.xyz).

## Repository Structure

The Radius SDKs are organized in a monorepo structure with a consistent layout across all language implementations:

```
radius-sdks/
├── CONTRIBUTING.md          # This file
├── README.md                # Repository overview
├── go/                      # Go SDK
│   ├── CONTRIBUTING.md      # Go-specific guidelines
│   ├── README.md            # Go SDK documentation
│   ├── radius/              # Public SDK package
│   ├── src/                 # Source code implementation
│   └── test/                # Test utilities and integration tests
├── typescript/              # TypeScript SDK
│   ├── CONTRIBUTING.md      # TypeScript-specific guidelines
│   ├── README.md            # TypeScript SDK documentation
│   ├── src/                 # Source code implementation
│   └── test/                # Test utilities and integration tests
├── python/                  # Python SDK (coming soon)
│   └── README.md            # Python SDK information
└── rust/                    # Rust SDK (coming soon)
    └── README.md            # Rust SDK information
```

Each SDK maintains its own:
- Build configuration
- Tests
- Documentation
- Contributing guidelines

The consistent folder structure enables:
- Clear separation of public API and implementation details
- Isolation of test code from production code
- Consistent development experience across languages
- Easier integration and maintenance

## Core Design Principles

The Radius SDKs follow these core design principles:

### 1. Consistent Interface, Idiomatic Implementation

While the SDKs provide a consistent interface across languages, they implement these interfaces in ways that feel
natural to each language:

- **Go SDK**: Uses interfaces, struct embedding, and options pattern
- **TypeScript SDK**: Leverages classes, async/await, and strong typing

Example of consistent interface with idiomatic implementation:

```go
// Go
client, err := radius.NewClient(url, radius.WithLogger(log.Printf))
```

```typescript
// TypeScript
const client = await NewClient(url, withLogger(console.log))
```

### 2. Selective Use of the Functional Options Pattern

Both SDKs use the functional options pattern selectively for configuration. We apply this pattern when:

- Configuring behavior with truly optional parameters
- Extending functionality without changing method signatures
- Providing configuration that may grow over time

We do NOT use the options pattern for:
- Required parameters
- Factory methods that need specific inputs to create an object
- Core functionality that should always be explicitly provided

This selective approach allows:
- Clarity about what's required vs. optional
- Type safety for all parameters
- Self-documenting method signatures
- Future extensibility for configuration options

### 3. Clear Error Handling

Each SDK follows language-specific best practices for error handling:
- Go: Multiple return values with explicit error types
- TypeScript: Promise rejections and Error objects

### 4. Interface-Based Design with Concrete Implementations

The SDKs use interfaces to define contracts between components while providing concrete implementations that are
optimized for each language:

- **Go**: Uses interfaces primarily for client capabilities and uses concrete types for core data structures
- **TypeScript**: Uses TypeScript interfaces for type safety while providing concrete class implementations

This approach enables:
- Easier testing through mocks
- Flexible implementations
- Clear API boundaries
- Type safety and intellisense support
- Direct access to commonly used properties

### 5. Consistent Directory Structure

Both SDKs follow similar package/module organization:
- `accounts/`: Account management and value transfer operations
- `auth/`: Authentication and transaction signing
- `client/`: Main client implementation with core Radius interactions
- `common/`: Core data structures (Address, Transaction, ABI, etc.) and utilities
- `contracts/`: Smart contract deployment and interaction
- `providers/eth/`: Abstraction layer for Ethereum library functionality
- `transport/`: HTTP communication and request/response handling

### 6. Unified Public Interface

While we maintain a well-organized internal structure, all public API elements should be accessible through a single
import path:

```go
// Go - all public functionality accessible through single import
import "github.com/radiustechsystems/sdk/go/radius"
```

```typescript
// TypeScript - all public functionality accessible through single import
import { NewClient, NewAccount, Address } from '@radiustechsystems/sdk';
```

This approach:
- Simplifies the developer experience
- Reduces import statements in user code
- Provides a clear distinction between public and internal APIs 
- Allows internal refactoring without breaking user code

Our folder structure enforces this pattern:
- `radius/` directory contains the public API
- `src/` directory contains the implementation details
- `test/` directory contains test utilities and integration tests
- Source code should never be directly imported from `src/` in user code

## Development Workflow

1. Fork the repository
2. Clone your fork: `git clone https://github.com/your-username/sdk.git`
3. Create a new branch: `git checkout -b my-feature`
4. Make your changes
5. Run tests and linting (using language-specific scripts)
6. Push to your fork and submit a pull request

## Commit Messages

We use [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

Types:
- `feat:` New features
- `fix:` Bug fixes
- `docs:` Documentation only
- `style:` Code style changes
- `refactor:` Non-bug-fixing code changes
- `test:` Test updates
- `chore:` Build process updates

## Pull Requests

1. Clear title following conventional commits
2. Detailed description of changes
3. Reference related issues
4. Update documentation
5. Add tests
6. Update CHANGELOG.md
7. Ensure CI checks pass

## Documentation

- Update README.md files
- Document public APIs
- Update CHANGELOG.md
- Include examples
- Follow language-specific documentation styles

## Language-Specific Guidelines

Please refer to the language-specific CONTRIBUTING.md files:
- [Go Contributing Guide](go/CONTRIBUTING.md)
- [TypeScript Contributing Guide](typescript/CONTRIBUTING.md)

These guides contain:
- Setup instructions
- Testing requirements
- Style guidelines
- Tool configurations
- Language-specific patterns

## Questions?

If you have questions:
1. Check existing issues
2. Create a new issue with `question` label
3. Ask in your PR if you're working on code

Thank you for contributing to Radius SDKs!
