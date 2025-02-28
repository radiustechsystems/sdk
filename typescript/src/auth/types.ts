/**
 * The auth types module defines interfaces for signing transactions and messages.
 * It provides the foundation for different signer implementations.
 */
import { Address, Hash, HttpClient, SignedTransaction, Transaction } from '../common';
import { BigNumberish, BytesLike } from '../providers/eth';

/**
 * Signer interface for cryptographically signing messages and transactions
 * Different implementations provide different mechanisms for accessing private keys
 */
export interface Signer {
  /**
   * Returns the Radius account address associated with the Signer
   * @returns The address of the signer
   */
  address(): Address;

  /**
   * Returns the Chain ID associated with the Signer
   * @returns The chain ID used for transaction signing
   */
  chainID(): BigNumberish;

  /**
   * Computes the hash of a transaction
   * @param transaction The transaction to hash
   * @returns The transaction hash
   */
  hash(transaction: Transaction): Hash;

  /**
   * Signs a message using the EIP-191 standard
   * @param message The message bytes to sign
   * @returns The signature bytes
   * @throws Error if signing fails
   */
  signMessage(message: BytesLike): Promise<Uint8Array>;

  /**
   * Signs a transaction using the EIP-155 standard
   * @param transaction The transaction to sign
   * @returns The signed transaction ready to be sent to the network
   * @throws Error if signing fails
   */
  signTransaction(transaction: Transaction): Promise<SignedTransaction>;
}

/**
 * A client interface for the Radius Client methods that may be required by the Signer
 * This interface is implemented by the main Radius Client
 */
export interface SignerClient {
  /**
   * Returns the Radius chain ID, which is used to sign transactions
   * @returns The chain ID of the connected network
   * @throws Error if the chain ID cannot be retrieved
   */
  chainID(): Promise<BigNumberish>;

  /**
   * Returns the HTTP client used by the client to make requests
   * @returns The HTTP client used for API requests
   */
  httpClient(): HttpClient;
}
