import { Signer } from '../auth';
import { Address, HttpClient, Receipt, Transaction } from '../common';
import { BigNumberish } from '../providers/eth';

/**
 * Client interface for account operations.
 * This interface is implemented by the main Radius Client and provides
 * core functionality for interacting with Radius accounts.
 */
export interface AccountClient {
  /**
   * Gets the balance of an account in wei.
   *
   * @param address Address to check the balance for
   * @returns The account balance in wei
   * @throws Error if the balance cannot be retrieved from the network
   */
  balanceAt(address: Address): Promise<bigint>;

  /**
   * Returns the Radius chain ID, which is used to sign transactions.
   *
   * @returns The chain ID of the connected network
   * @throws Error if the chain ID cannot be retrieved
   */
  chainID(): Promise<bigint>;

  /**
   * Estimates the gas cost of a transaction with a safety margin.
   *
   * @param tx Transaction to estimate gas for
   * @returns The estimated gas cost in gas units
   * @throws Error if the gas estimation fails
   */
  estimateGas(tx: Transaction): Promise<bigint>;

  /**
   * Returns the HTTP client used by the client to make requests.
   *
   * @returns The HTTP client used for API requests
   */
  httpClient(): HttpClient;

  /**
   * Returns the next nonce (transaction count) for an account.
   *
   * @param address Address to check the nonce for
   * @returns The next nonce to use for transactions
   * @throws Error if the nonce cannot be retrieved from the network
   */
  pendingNonceAt(address: Address): Promise<number>;

  /**
   * Sends native currency to a recipient address.
   *
   * @param signer The signer used to sign the transaction
   * @param recipient Destination address to receive the funds
   * @param value Amount of native currency to send in wei
   * @returns Receipt of the completed transaction
   * @throws Error if the transaction fails
   * @throws Error if the transaction receipt is not returned
   */
  send(signer: Signer, recipient: Address, value: BigNumberish): Promise<Receipt>;
}
