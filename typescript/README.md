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
- npm, yarn, or pnpm

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
const abi = new ABIFromJSON(`[{"inputs":[],"name":"get","outputs":[{"type":"uint256"}],"type":"function"},{"inputs":[{"type":"uint256"}],"name":"set","type":"function"}]`);
const bytecode = BytecodeFromHex('608060405234801561001057600080fd5b50610150806100...');

// Deploy the contract
const contract = await client.deployContract(account.signer, bytecode, abi);
```

### Interact with a Smart Contract

```typescript
import { NewContract, AddressFromHex, ABIFromJSON } from '@radiustechsystems/sdk';

// Reference an existing contract
const address = AddressFromHex('0x...');
const abi = new ABIFromJSON(`[{"inputs":[],"name":"get","outputs":[{"type":"uint256"}],"type":"function"},{"inputs":[{"type":"uint256"}],"name":"set","type":"function"}]`);
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
import { NewClefSigner, NewAccount, withSigner, AddressFromHex } from '@radiustechsystems/sdk';

// Use an external signer (e.g., Clef)
const address = AddressFromHex('0x...');
const clefSigner = NewClefSigner(address, client, 'http://localhost:8550');
const account = await NewAccount(withSigner(clefSigner));

// Or create a custom signer
const customSigner = new MyCustomSigner();
const customSignerAccount = NewAccount(withSigner(customSigner));
```

### Logging and Request Interceptors

```typescript
import { NewClient, withLogger, withInterceptor } from '@radiustechsystems/sdk';

const client = await NewClient('https://your-radius-endpoint',
    withLogger((message, data) => {
        console.log(message, data);
    }),
    withInterceptor(async (reqBody, response) => {
        // Process or log responses
        return response;
    })
);
```

### Custom HTTP Client

```typescript
import { NewClient, withHttpClient } from '@radiustechsystems/sdk';

const client = await NewClient('https://your-radius-endpoint',
    withHttpClient(customHttpClient)
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
