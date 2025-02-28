import { Signer } from '../auth';
import { Address, Receipt, SignedTransaction, Transaction, zeroAddress } from '../common';
import { BigNumberish, BytesLike } from '../providers/eth';
import { AccountOption, AccountOptions } from './options';
import { AccountClient } from './types';

/**
 * Account represents a Radius account that can be used to sign transactions.
 * This class provides methods for checking balance, retrieving nonce, and
 * signing messages and transactions.
 */
export class Account {
  /**
   * The signer used to cryptographically sign messages and transactions
   */
  signer?: Signer;

  /**
   * Creates a new Account instance
   * @param signer Optional signer to use with this account
   */
  constructor(signer?: Signer) {
    this.signer = signer;
  }

  /**
   * Creates a new Account with the given options
   * @param opts Functional options to configure the account (e.g., WithSigner)
   * @returns A new Account instance configured with the provided options
   */
  static async New(...opts: AccountOption[]): Promise<Account> {
    const options: AccountOptions = {};
    for (const opt of opts) {
      await opt(options);
    }
    return new Account(options.signer);
  }

  /**
   * Returns the address of the account
   * @returns The account address, or zero address if no signer is available
   */
  address(): Address {
    return this.signer?.address() ?? zeroAddress();
  }

  /**
   * Returns the balance of the account in wei
   * @param client Radius client instance used to query the balance
   * @returns The account balance in wei
   * @throws Error if the balance cannot be retrieved from the network
   */
  async balance(client: AccountClient): Promise<bigint> {
    return client.balanceAt(this.address());
  }

  /**
   * Returns the next nonce (transaction count) of the account
   * @param client Radius client instance used to query the nonce
   * @returns The next nonce to use for transactions
   * @throws Error if the nonce cannot be retrieved from the network
   */
  async nonce(client: AccountClient): Promise<number> {
    return client.pendingNonceAt(this.address());
  }

  /**
   * Sends native currency to a recipient address
   * @param client Radius client instance used to send the transaction
   * @param recipient Destination address to receive the funds
   * @param value Amount of native currency to send in wei
   * @returns Receipt of the completed transaction
   * @throws Error if no signer is available
   * @throws Error if the transaction fails
   */
  async send(client: AccountClient, recipient: Address, value: BigNumberish): Promise<Receipt> {
    if (!this.signer) {
      throw new Error('Signer is required for sending transactions');
    }
    return client.send(this.signer, recipient, value);
  }

  /**
   * Signs a message using the EIP-191 standard
   * @param message Message bytes to sign
   * @returns The signature bytes
   * @throws Error if no signer is available
   * @throws Error if signing fails
   */
  async signMessage(message: BytesLike): Promise<Uint8Array> {
    if (!this.signer) {
      throw new Error('Signer is required for signing messages');
    }
    return this.signer.signMessage(message);
  }

  /**
   * Signs a transaction using the EIP-155 standard
   * @param transaction Transaction to sign
   * @returns The signed transaction ready to be sent to the network
   * @throws Error if no signer is available
   * @throws Error if signing fails
   */
  async signTransaction(transaction: Transaction): Promise<SignedTransaction> {
    if (!this.signer) {
      throw new Error('Signer is required for sending transactions');
    }
    return this.signer.signTransaction(transaction);
  }
}
