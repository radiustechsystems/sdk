import { PrivateKeySigner, Signer } from '../auth';
import { AccountClient } from './types';

/**
 * A function that configures a Radius account.
 * This is used as a functional option pattern for creating new accounts.
 */
export type AccountOption = (options: AccountOptions) => Promise<void>;

/**
 * Options for creating an account.
 * Contains configuration values that can be set using functional options.
 */
export interface AccountOptions {
  /**
   * The signer to use with this account
   */
  signer?: Signer;
}

/**
 * Create an AccountOption that sets the account address and signer using a private key.
 * The private key will be stored in memory, so for production systems with high security
 * requirements, consider using withSigner instead, along with a hardware security module
 * or key management service.
 *
 * @param key Private key as a hex string
 * @param client AccountClient instance for network operations
 * @returns An AccountOption function that configures an Account with the provided private key
 */
export function withPrivateKey(key: string, client: AccountClient): AccountOption {
  return async (options: AccountOptions) => {
    options.signer = new PrivateKeySigner(key, client);
  };
}

/**
 * Create an AccountOption that sets the account address and signer using a custom Signer implementation.
 * This is useful when you want to use a custom signing implementation, such as a hardware
 * security module or key management service.
 *
 * @param signer Signer instance for signing transactions and messages
 * @returns An AccountOption function that configures an Account with the provided signer
 */
export function withSigner(signer: Signer): AccountOption {
  return async (options: AccountOptions) => {
    options.signer = signer;
  };
}
