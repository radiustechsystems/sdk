import { ABI, Account, Client, withLogger, withPrivateKey } from '../radius';
import { bytecodeFromHex } from '../src/common';
import { eth } from '../src/providers/eth';

// Constants for testing
export const MIN_TEST_ACCOUNT_BALANCE = BigInt('1000000000000000000');
export const SIMPLE_STORAGE_ABI = `[{"inputs":[],"name":"get","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"x","type":"uint256"}],"name":"set","outputs":[],"stateMutability":"nonpayable","type":"function"}]`;
export const SIMPLE_STORAGE_BIN =
  '6080604052348015600e575f5ffd5b5060a580601a5f395ff3fe6080604052348015600e575f5ffd5b50600436106030575f3560e01c806360fe47b11460345780636d4ce63c146045575b5f5ffd5b6043603f3660046059565b5f55565b005b5f5460405190815260200160405180910390f35b5f602082840312156068575f5ffd5b503591905056fea26469706673582212207655d86666fa8aa75666db8416e0f5db680914358a57e84aa369d9250218247f64736f6c634300081c0033';

/**
 * Creates a client for testing
 * @param endpoint The endpoint URL
 * @returns A new client instance or null if endpoint is not set
 */
export async function createTestClient(endpoint: string): Promise<Client | null> {
  try {
    return await Client.New(endpoint, withLogger(console.log));
  } catch (error) {
    console.error('Failed to create test client:', error);
    return null;
  }
}

/**
 * Creates a test account using a random private key
 * @param client The client instance
 * @returns A new account
 */
export async function createTestAccount(client: Client): Promise<Account> {
  const wallet = eth.Wallet.createRandom();
  const account = await Account.New(withPrivateKey(wallet.privateKey, client));

  if (!account) {
    throw new Error('Failed to create test account');
  }

  return account;
}

/**
 * Gets a funded account using the private key from environment variables
 * @param client The client instance
 * @param privateKey The private key to use
 * @returns A funded account or null if private key is not set
 */
export async function getFundedAccount(
  client: Client,
  privateKey: string
): Promise<Account | null> {
  const fundedAccount = await Account.New(withPrivateKey(privateKey, client));

  if (!fundedAccount) {
    console.error('Failed to create funded account');
    return null;
  }

  const balance = await fundedAccount.balance(client);
  if (balance < MIN_TEST_ACCOUNT_BALANCE) {
    console.error('Insufficient account balance');
    return null;
  }

  return fundedAccount;
}

/**
 * Helper to get bytecode for testing
 * @param hexString The hex string to convert to bytecode
 * @returns Uint8Array of bytecode
 */
export function getTestBytecode(hexString: string): Uint8Array {
  const bytes = bytecodeFromHex(hexString);
  if (!bytes) {
    throw new Error('Failed to convert hex string to bytecode');
  }
  return bytes;
}

/**
 * Helper to create an ABI object for testing
 * @param abiString The ABI JSON string
 * @returns ABI object
 */
export function getTestABI(abiString: string): ABI {
  return new ABI(abiString);
}
