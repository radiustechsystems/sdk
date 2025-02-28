# Contributing to the Radius TypeScript SDK

This guide provides specific guidelines for contributing to the TypeScript SDK.

## Development Setup

1. Install Node.js 20.12.2 or later
2. Install a package manager (npm, pnpm 9.14.2+, or yarn)
3. Clone the repository
4. Install dependencies:
   ```bash
   # Using npm
   npm install
   
   # Using pnpm
   pnpm install
   
   # Using yarn
   yarn
   ```
5. Make sure tests pass: 
   ```bash
   # Using npm
   npm test
   
   # Using pnpm
   pnpm test
   
   # Using yarn
   yarn test
   ```

## TypeScript-Specific Patterns

### Provider Abstraction

The SDK uses an abstraction layer for Ethereum-compatible libraries that provides data structure mappings. This is managed through the `providers/eth` module:

```typescript
// ✅ Correct
import { eth } from '../providers/eth';
import type { BigNumberish } from '../providers/eth';

// ❌ Incorrect
import { ethers } from 'ethers';
import type { BigNumberish } from 'ethers';
```

While the SDK has its own concrete implementations of core data structures (like addresses and transactions), the provider module maps these to Ethereum types when needed to leverage the functionality provided by established libraries.

This abstraction:

- Maintains independent implementations optimized for Radius
- Leverages battle-tested Ethereum libraries when beneficial
- Presents a consistent interface regardless of the underlying implementation
- Handles errors internally with descriptive messages
- Simplifies testing and mocking
- Provides a single point for managing external dependencies

#### Error Handling in Provider Abstraction

The `providers/eth` module handles errors internally to maintain loose coupling with the underlying library:

```typescript
// ✅ Correct: Use the wrapped utility functions that handle errors
function someOperation() {
  try {
    // This function already has internal error handling
    const bytes = eth.getBytes(someValue);
    return bytes;
  } catch (error) {
    // The error from eth.getBytes will have a descriptive message
    throw new Error(`Operation failed: ${error.message}`);
  }
}

// ❌ Incorrect: Don't use ethers directly, missing error handling
function someOperation() {
  try {
    // This doesn't have internal error handling
    const bytes = ethers.getBytes(someValue);
    return bytes;
  } catch (error) {
    // The error from ethers may have a cryptic message
    throw error;
  }
}
```

Key principles for error handling in the abstraction layer:

- Functions in the `providers/eth` module have internal error handling
- Errors are transformed with descriptive messages
- External library implementation details are hidden
- Consumer code doesn't need to handle low-level library errors

### Async Patterns

We use async/await consistently and expose top-level async functions for creating classes:

```typescript
// In class definition
export class MyClass {
    private constructor(/* parameters */) {}
    
    static async New(url: string, ...opts: Option[]): Promise<MyClass> {
        // implementation
        return new MyClass(/* parameters */);
    }
}

// In SDK main export
export async function NewMyClass(url: string, ...opts: Option[]): Promise<MyClass> {
    return MyClass.New(url, ...opts);
}
```

### Options Pattern

We use the functional options pattern for optional configuration, but not for required parameters:

```typescript
// Good use of options pattern - for optional configuration
export async function NewClient(url: string, ...opts: Option[]): Promise<Client> {
    const options = defaultOptions();
    for (const opt of opts) {
        opt(options);
    }
    // Use options to configure the client
    return new Client(url, options);
}

// Not using options for required parameters
export function NewContract(address: Address, abi: ABI): Contract {
    return new Contract(address, abi);
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
- Self-documenting parameter names with the `with*` prefix
- Extensibility without breaking changes
- Graceful defaults for optional parameters

### Type Safety

- Use strict TypeScript configuration
- Avoid `any` type
- Use branded types when appropriate
- Leverage union types for better type safety

### Class Structure Pattern

We use two distinct patterns for classes in our SDK, depending on class complexity and lifecycle requirements:

#### 1. Simple Value Objects

For simple data structures that don't require async initialization or complex configuration:

```typescript
// Simple value object with public constructor
export class Address {
  private readonly data: Uint8Array;

  constructor(data: Address | BytesLike | string) {
    // Immediate initialization
    if (data instanceof Address) {
      this.data = data.bytes();
    } else if (typeof data === 'string') {
      this.data = eth.getBytes(data.startsWith('0x') ? data : `0x${data}`);
    } else {
      this.data = eth.getBytes(data);
    }
  }
  
  // Methods
  bytes(): Uint8Array { return this.data; }
}
```

When to use this pattern:
- Simple value objects with synchronous initialization
- Core data structures (Address, ABI, etc.)
- Classes with minimal or no dependencies on external systems
- Objects that don't require complex configuration

#### 2. Service Objects with Async Initialization

For classes that require async initialization or complex configuration:

```typescript
// Service object with both public constructor and static factory method
export class Client {
  private readonly ethClient: Provider;

  // Public constructor for direct instantiation
  constructor(provider: Provider, httpClient?: HttpClient) {
    this.ethClient = provider;
    this._httpClient = httpClient ?? globalThis.fetch;
  }

  // Static factory for configuration and async initialization
  static async New(url: string, ...opts: ClientOption[]): Promise<Client> {
    const options: ClientOptions = {};
    
    for (const opt of opts) {
      opt(options);
    }
    
    // Async initialization work (like connecting to network)
    const provider = new eth.JsonRpcProvider(url);
    await provider.getNetwork();
    
    return new Client(provider, options.httpClient);
  }
}

// In sdk.ts (main export)
export async function NewClient(url: string, ...opts: ClientOption[]): Promise<Client> {
   return Client.New(url, ...opts);
}
```

When to use this pattern:
- Classes that require async initialization or validation
- Services with dependencies on external systems
- Objects that support the functional options pattern
- Classes with complex configuration requirements

Both patterns should:
- Use explicit interface implementation
- Mark properties as readonly when appropriate
- Hide implementation details with private modifiers
- Have corresponding top-level factory functions in the main SDK export

### Method Arguments

- Use variadic arguments (...args) instead of arrays where possible
- When interfacing with underlying libraries that require arrays, convert variadic arguments internally
- Document when methods accept variable arguments

Example:
```typescript
// ✅ Use variadic arguments for optional parameters
async someFunc(foo: string, bar: Number, ...args: unknown[]): Promise<void>

// ❌ Avoid using arrays for optional parameters
async someFunc(foo: string, bar: Number, args: unknown[] = []): Promise<void>
```

## Testing

### Unit Tests

We use Vitest for testing:
```typescript
describe('MyClass_someFunc', () => {
   const tests = [
      {
         name: 'simple case',
         foo: 'test',
         bar: 42,
         want: 'expected result',
      },
      {
         name: 'error case',
         foo: 'invalid',
         wantErr: true,
      },
   ];

   tests.forEach(({ name, foo, bar, want, wantErr }) => {
      test(name, async () => {
         if (wantErr) {
            await expect(MyClass.New(foo, bar)).rejects.toThrow();
         } else {
            const mc = await MyClass.New(foo, bar);
            const result = await mc.someFunc();
            expect(result).toBe(want);
         }
      });
   });
});
```

### Integration Tests

- Use separate test files
- Handle cleanup properly
- Set appropriate timeouts
- Use environment variables for configuration

## Style Guidelines

### Code Organization

- Keep files focused and small
- Clear module exports
- Consistent import ordering
- Follow TypeScript project references
- Maintain our repository structure:
   - `radius/`: Public API package that users import
   - `src/`: Implementation details
   - `test/`: Test utilities and integration tests

### Naming Conventions

- Use PascalCase for classes and interfaces
- Use camelCase for methods and properties
- Use UPPER_CASE for constants

## Quality Checks

Before submitting:

1. Run tests:
   ```bash
   pnpm test
   ```

2. Run linter:
   ```bash
   pnpm lint
   ```

3. Format code:
   ```bash
   pnpm format
   ```

For convenience, you can fix all lint and formatting errors by running the cleanup script:

- On Windows: `cleanup.bat`
- On MacOS/Linux: `./cleanup.sh`

## Build System

We use tsup for building:
- Generates both ESM and CJS outputs
- Handles type definitions
- Manages source maps
- Optimizes production builds

## Dependencies

- Keep dependencies minimal
- Use peer dependencies appropriately
- Lock versions in package.json
- Document why each dependency is needed

## Common Patterns

### Unified Public API

All public SDK functionality should be exposed through the main package:

```typescript
// External usage should only import from the main package
import { NewClient, Address, ABI } from '@radiustechsystems/sdk';

// Inside the SDK, import from specific modules
import { Client } from '../src/client';
import { Address } from '../src/common';
```

The `sdk.ts` and `index.ts` files in the `radius/` directory re-export all public types and functions, providing a clean, single-import API for users while maintaining internal organization.

This pattern aligns with our folder structure:
- Implementation details live in the `src/` directory
- Public API is defined in the `radius/` directory
- Test utilities are in the `test/` directory
- Users should never import directly from implementation modules

### Factory Method Naming

Our SDK uses a specific naming convention for factory methods:

```typescript
// Class-level static factory methods
static async New(url: string, ...opts: Option[]): Promise<Client>

// Top-level SDK factory functions (in sdk.ts)
export async function NewClient(url: string, ...opts: Option[]): Promise<Client>
```

Key requirements:
- Class-level factory methods use PascalCase with `New` prefix
- Top-level functions use PascalCase with class name prefix (`NewClass`)
- Factory functions that create objects or perform async initialization should return a Promise
- Top-level functions should generally be thin wrappers around class-level static methods

This pattern provides:
- A clean, consistent API for end users
- Clear distinction between constructors and factory methods
- Better discoverability of SDK capabilities through the main import
- The ability to perform validation and async initialization before object creation

### Error Handling

Follow these principles:
- Use typed errors
- Provide meaningful error messages
- Handle async errors properly
- Use error subclasses when appropriate
- Encapsulate error handling in abstraction layers
- Prevent leaking implementation details in error messages

Example of proper error handling pattern:
```typescript
try {
   const result = await this.provider.send(tx);
   if (!result) {
      throw new Error('Transaction failed: no result returned');
   }
   return result;
} catch (error) {
   throw new Error(`Failed to send transaction: ${error instanceof Error ? error.message : String(error)}`);
}
```

For abstraction layers (like `providers/eth`), include error handling within the utility functions:
```typescript
// In providers/eth/utils.ts
export const eth = {
  // This function includes internal error handling
  getBytes: (value: string | Uint8Array) => {
    try {
      return ethGetBytes(value);
    } catch (error) {
      throw new Error(`Failed to convert to bytes: ${error instanceof Error ? error.message : String(error)}`);
    }
  },
  // ...other functions
};
```

This approach:
- Centralizes error handling in abstraction layers
- Provides consistent, descriptive error messages
- Prevents implementation details from leaking
- Reduces duplicate try/catch blocks across the codebase
- Maintains loose coupling with underlying libraries

### Interface Implementation

Implement interfaces explicitly:
```typescript
interface SomeInterface {
   someFunc(): Promise<string>;
}

export class MyImplementation implements SomeInterface {
   async someFunc(): Promise<string> {
      return 'result';
   }
}
```

## Documentation

- Update README.md
- Include TSDoc comments
- Maintain CHANGELOG.md
- Document breaking changes
- Include usage examples

## Type Definitions

- Export public types
- Use precise types when possible
- Document complex types
- Use utility types appropriately

## Biome Configuration

We use Biome for linting and formatting:
- Follow configured rules
- Use provided formatter
- Address all linting errors
- Maintain configuration in biome.json
- The `dist` and `node_modules` directories are excluded from linting

When running the linter:
```bash
# Lint entire project
pnpm lint

# Lint specific directories
pnpm lint:fix -- ./src/providers/eth/
```

Our Biome configuration includes directory exclusions to avoid linting build artifacts:
```json
{
  "files": {
    "ignore": ["dist/**/*", "node_modules/**/*"]
  },
  // Other Biome configuration...
}
```
