# Radius TypeScript SDK

A TypeScript client library for interacting with [Radius](https://radiustech.xyz/), providing a simple and type-safe
way to interact with the Radius platform.

## Features

- Account management and transaction signing
- Smart contract deployment and interaction
- Optional request logging and interceptors
- EVM compatibility with high performance
- Type-safe contract interactions

## Requirements

- Node.js >= 20.12.2

## Installation

```bash
# Using npm
npm install @radiustechsystems/sdk

# Using pnpm
pnpm add @radiustechsystems/sdk

# Using yarn
yarn add @radiustechsystems/sdk
```

## Quickstart Examples

### Connect to Radius and Create an Account

```typescript
import { NewClient, NewAccount, withPrivateKey } from '@radiustechsystems/sdk';

// Connect to Radius
const client = await NewClient('https://your-radius-endpoint');

// Create an account using a private key
const account = await NewAccount(withPrivateKey('your-private-key', client));
```

### Transfer Value Between Accounts

```typescript
import { AddressFromHex } from '@radiustechsystems/sdk';

// Send 100 tokens to another account
const recipient = AddressFromHex('0x...');
const amount = BigInt(100);
const receipt = await account.send(client, recipient, amount);

console.log('Transaction hash:', receipt.txHash.hex());
```

### Deploy a Smart Contract

```typescript
import { ABIFromJSON, BytecodeFromHex } from '@radiustechsystems/sdk';

// Parse ABI and bytecode
const abi = ABIFromJSON(`[{"inputs":[],"name":"get","outputs":[{"type":"uint256"}],"type":"function"},{"inputs":[{"type":"uint256"}],"name":"set","type":"function"}]`);
const bytecode = BytecodeFromHex('608060405234801561001057600080fd5b50610150806100...');

// Deploy the contract
const contract = await client.deployContract(account.signer, bytecode, abi);
```

### Interact with a Smart Contract

```typescript
import { NewContract, AddressFromHex, ABIFromJSON } from '@radiustechsystems/sdk';

// Reference an existing contract
const address = AddressFromHex('0x...');
const abi = ABIFromJSON(`[{"inputs":[],"name":"get","outputs":[{"type":"uint256"}],"type":"function"},{"inputs":[{"type":"uint256"}],"name":"set","type":"function"}]`);
const contract = NewContract(address, abi);

// Write to the contract
const value = BigInt(42);
const receipt = await contract.execute(client, account.signer, 'set', value);

// Read from the contract
const result = await contract.call(client, 'get');
console.log('Stored value:', result[0]);
```

## Advanced Features

### Custom Transaction Signing

```typescript
import { Address, BigNumberish, BytesLike, Hash, SignedTransaction, Transaction } from '@radiustechsystems/sdk';

class MyCustomSigner {
    address(): Address { /* ... */ }
    chainID(): BigNumberish { /* ... */ }
    hash(transaction: Transaction): Hash { /* ... */ }
    signMessage(message: BytesLike): Promise<Uint8Array> { /* ... */ }
    signTransaction(transaction: Transaction): Promise<SignedTransaction> { /* ... */ }
    constructor(...args) { /* ... */ }
}
const signer = new MyCustomSigner(...args);
const account = NewAccount(withSigner(signer));
```

### Logging and Request Interceptors

```typescript
import { NewClient, withLogger, withInterceptor } from '@radiustechsystems/sdk';

const client = await NewClient('https://your-radius-endpoint',
    withLogger((message, data) => {
        console.log(message, data);
    }),
    withInterceptor(async (reqBody, response) => {
        // Examine request body, modify response, etc.
        return response;
    })
);
```

### Custom HTTP Client

```typescript
import { NewClient, withHttpClient } from '@radiustechsystems/sdk';

const client = await NewClient('https://your-radius-endpoint',
    withHttpClient(async (url: string | URL | Request, init?: RequestInit | undefined): Promise<Response> => {
        // Make a custom HTTP request, or use a library like axios
    })
);
```

## Resources

- [Website](https://radiustech.xyz/)
- [Testnet Access](https://docs.radiustech.xyz/radius-testnet-access) 
- [GitHub Issues](https://github.com/radiustechsystems/sdks/issues)
- [Changelog](CHANGELOG.md)

## Contributing

Please see the [TypeScript SDK Contributing Guide](CONTRIBUTING.md) for detailed information about contributing to this
SDK. For repository-wide guidelines, see the [General Contributing Guide](../CONTRIBUTING.md).

## License

This project is licensed under the [MIT License](../LICENSE).
