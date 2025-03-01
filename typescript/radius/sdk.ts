import { Account } from '../src/accounts';
import type { AccountOption } from '../src/accounts';
import { withPrivateKey, withSigner } from '../src/accounts';
import { ClefSigner, PrivateKeySigner, Signer, SignerClient } from '../src/auth';
import { Client, ClientOption } from '../src/client';
import { withHttpClient, withInterceptor, withLogger } from '../src/client';
import { ABI, Address, Event, Hash, HttpClient, Receipt, Transaction } from '../src/common';
import type { SignedTransaction } from '../src/common';
import {
  abiFromJSON as ABIFromJSON,
  addressFromHex as AddressFromHex,
  bytecodeFromHex as BytecodeFromHex,
  hashFromHex as HashFromHex,
  zeroAddress,
} from '../src/common';
import { Contract } from '../src/contracts';
import type { BigNumberish, BytesLike } from '../src/providers/eth';
import type { Interceptor, Logf } from '../src/transport';

// Re-export classes
export {
  ABI,
  Account,
  Address,
  Client,
  ClefSigner,
  Contract,
  Event,
  Hash,
  PrivateKeySigner,
  Receipt,
  Transaction,
};

// Re-export types
export type {
  BigNumberish,
  BytesLike,
  HttpClient,
  Interceptor,
  Logf,
  Signer,
  SignerClient,
  SignedTransaction,
};

// Re-export functions
export {
  ABIFromJSON,
  AddressFromHex,
  BytecodeFromHex,
  HashFromHex,
  withHttpClient,
  withInterceptor,
  withLogger,
  withPrivateKey,
  withSigner,
};

/**
 * Maximum gas allowed for a transaction.
 * This value represents the upper limit that can be used for gas in a transaction (0xffffffff).
 */
export const MaxGas = BigInt('0xffffffff');

/**
 * Creates a new ClefSigner with the given address, client, and Clef URL.
 * ClefSigner provides a way to sign transactions using Clef as an external signing service.
 *
 * @param address The address to use for signing
 * @param client The Radius client to use for transaction-related operations
 * @param clefURL The URL of the Clef server
 * @returns A new ClefSigner instance
 * @throws Error if unable to connect to the Clef server
 */
export function NewClefSigner(address: Address, client: SignerClient, clefURL: string): ClefSigner {
  return new ClefSigner(address, client, clefURL);
}

/**
 * Creates a new Client with the given URL and options.
 * The client is the main entry point for interacting with the Radius platform.
 *
 * @param url The URL of the Ethereum node
 * @param opts Additional options for the client configuration
 * @returns A new Client instance
 * @throws Error if the client cannot be created or cannot connect to the network
 */
export async function NewClient(url: string, ...opts: ClientOption[]): Promise<Client> {
  return Client.New(url, ...opts);
}

/**
 * Creates a new Contract with the given address and ABI.
 * The Contract object allows interaction with smart contracts deployed on Radius.
 *
 * @param address The contract address on Radius
 * @param abi The contract ABI (Application Binary Interface) defining its methods and events
 * @returns A new Contract instance
 */
export function NewContract(address: Address, abi: ABI): Contract {
  return new Contract(address, abi);
}

/**
 * Creates a new Account with the given options.
 * Accounts are used to represent Radius accounts and manage their keys.
 *
 * @param opts Account options for configuring the account (private key, signer, etc.)
 * @returns A new Account instance
 * @throws Error if the account cannot be created with the provided options
 */
export async function NewAccount(...opts: AccountOption[]): Promise<Account> {
  return Account.New(...opts);
}

/**
 * Creates a zero address (0x0000000000000000000000000000000000000000)
 * Used as a default value or to represent the zero address in the Ethereum ecosystem
 * @returns An Address instance representing the zero address
 */
export function ZeroAddress(): Address {
  return zeroAddress();
}
