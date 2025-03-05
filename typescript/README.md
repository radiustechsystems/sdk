# Radius TypeScript SDK

The official TypeScript client library for interacting with the [Radius platform](https://radiustech.xyz/), providing
a simple and idiomatic way to interact with Radius services.

## Features

- Account management and transaction signing
- Smart contract deployment and interaction
- Optional request logging and interceptors
- EVM compatibility with high performance & low latency

## Requirements

- Node.js >= 20.12
- Radius JSON-RPC endpoint: https://docs.radiustech.xyz/radius-testnet-access
- Ethereum private key: https://ethereum.org/en/developers/docs/accounts/#account-creation

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

Be sure to use your own `RADIUS_ENDPOINT` and `PRIVATE_KEY` values, as mentioned in the [Requirements](#requirements).

```typescript
import { Account, Client, NewClient, NewAccount, withPrivateKey } from '@radiustechsystems/sdk';

const RADIUS_ENDPOINT = "https://rpc.testnet.tryradi.us/xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx";
const PRIVATE_KEY = "fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd036415f";

const client: Client = await NewClient(RADIUS_ENDPOINT);
const account: Account = await NewAccount(withPrivateKey(PRIVATE_KEY, client));
```

Alternatively, using plain JavaScript and CommonJS `require` syntax:

```javascript
const { Account, Client, NewClient, NewAccount, withPrivateKey } = require('@radiustechsystems/sdk');
```

### Transfer Value Between Accounts

Here, we send 100 tokens to another account. Be sure to replace the recipient's address with one of your own.

```typescript
import { Address, AddressFromHex, Receipt } from '@radiustechsystems/sdk';

const recipient: Address = AddressFromHex('0x5e97870f263700f46aa00d967821199b9bc5a120'); // Recipient's address
const amount: bigint = BigInt(100);
const receipt: Receipt = await account.send(client, recipient, amount);

console.log('Transaction hash:', receipt.txHash.hex());
```

### Deploy a Smart Contract

Here, we deploy the [SimpleStorage.sol](https://github.com/radiustechsystems/sdk/tree/main/contracts/solidity)
example contract included in this SDK, with the application binary interface (ABI) and bytecode that were generated
using the Solidity compiler [solcjs](https://docs.soliditylang.org/en/latest/installing-solidity.html#npm-node-js).

```typescript
import { ABI, ABIFromJSON, BytecodeFromHex } from '@radiustechsystems/sdk';

// Parse ABI and bytecode of the SimpleStorage contract
const abi: ABI = ABIFromJSON(`[{"inputs":[],"name":"get","outputs":[{"type":"uint256"}],"type":"function"},{"inputs":[{"type":"uint256"}],"name":"set","type":"function"}]`);
const bytecode: Uint8Array = BytecodeFromHex('6080604052348015600e575f5ffd5b5060a580601a5f395ff3fe6080604052348015600e575f5ffd5b50600436106030575f3560e01c806360fe47b11460345780636d4ce63c146045575b5f5ffd5b6043603f3660046059565b5f55565b005b5f5460405190815260200160405180910390f35b5f602082840312156068575f5ffd5b503591905056fea26469706673582212207655d86666fa8aa75666db8416e0f5db680914358a57e84aa369d9250218247f64736f6c634300081c0033');

// Deploy the contract
const contract = await client.deployContract(account.signer, bytecode, abi);
```

### Interact with a Smart Contract

Assuming the contract was previously deployed (which is typically the case), we can interact with it using the contract
address and ABI. Be sure to replace the contract address with that of your own deployed contract.

```typescript
import { ABI, Address, AddressFromHex, ABIFromJSON, Contract, NewContract, Receipt } from '@radiustechsystems/sdk';

// Reference a previously deployed contract
const address: Address = AddressFromHex('0x5e97870f263700f46aa00d967821199b9bc5a120'); // Contract address
const abi: ABI = ABIFromJSON(`[{"inputs":[],"name":"get","outputs":[{"type":"uint256"}],"type":"function"},{"inputs":[{"type":"uint256"}],"name":"set","type":"function"}]`);
const contract: Contract = NewContract(address, abi);

// Write to the contract
const value: bigint = BigInt(42);
const receipt: Receipt = await contract.execute(client, account.signer, 'set', value);

// Read from the contract
const result: unknown[] = await contract.call(client, 'get');
console.log('Stored value:', result[0]);
```

## Advanced Features

### Custom Transaction Signing

```typescript
import { Address, BigNumberish, BytesLike, Hash, SignedTransaction, Signer, Transaction } from '@radiustechsystems/sdk';

class MyCustomSigner implements Signer {
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

All Radius SDKs are released under the [MIT License](../LICENSE).
