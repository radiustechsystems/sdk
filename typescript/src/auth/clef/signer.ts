/**
 * Implementation of the Signer interface using the Clef external signing tool.
 * This module provides a way to sign transactions and messages with Clef,
 * which manages keys securely outside the application.
 */
import { Address, Hash, SignedTransaction, Transaction } from '../../common';
import { BigNumberish, BytesLike, JsonRpcProvider, eth } from '../../providers/eth';
import { Signer, SignerClient } from '../types';

/**
 * A Signer implementation that uses the Clef JSON-RPC API.
 * Clef is a secure key management service that can be used to sign transactions without exposing
 * the private key to the application. This is useful for securing private keys in production systems.
 * Clef must be running and accessible to the application in order to use this signer.
 * Learn more about Clef here: https://geth.ethereum.org/docs/tools/clef/introduction
 */
export class ClefSigner implements Signer {
  /**
   * The address associated with this signer
   * @private
   */
  private readonly _address: Address;

  /**
   * The chain ID used for EIP-155 transaction signing
   * @private
   */
  private _chainID: BigNumberish;

  /**
   * The JSON-RPC client used to communicate with the Clef server
   * @private
   */
  private readonly client: JsonRpcProvider;

  /**
   * Create a new ClefSigner instance
   * @param address The address to use for signing
   * @param client The Radius client
   * @param clefURL The URL of the Clef server (e.g. "http://localhost:8550")
   * @throws Error if the connection to Clef fails
   */
  constructor(address: Address, client: SignerClient, clefURL: string) {
    this._address = address;
    this._chainID = 0n;
    this.client = new eth.JsonRpcProvider(clefURL, undefined, {
      staticNetwork: true,
      cacheTimeout: -1,
    });

    // Fetch chain ID from client
    client
      .chainID()
      .then((id) => {
        this._chainID = id;
      })
      .catch(() => {
        // Default to 0 if we can't get chain ID
      });

    // Verify Clef connection
    this.verifyClefConnection().catch((error) => {
      console.error('Failed to verify Clef connection:', error);
    });
  }

  /**
   * Get the account address of the signer
   * @returns The account address
   */
  address(): Address {
    return this._address;
  }

  /**
   * Get the chain ID used by the signer
   * @returns The chain ID
   */
  chainID(): BigNumberish {
    return this._chainID;
  }

  /**
   * Compute the hash of a transaction
   * @param transaction The transaction to hash
   * @returns The transaction hash
   */
  hash(transaction: Transaction): Hash {
    const tx = this.prepareTransaction(transaction);
    const unsignedTx = eth.Transaction.from(tx);
    return new Hash(unsignedTx.hash || '0x');
  }

  /**
   * Sign a message using the EIP-191 standard
   * @param message The message to sign
   * @returns The signature
   * @throws Error if signing fails
   */
  async signMessage(message: BytesLike): Promise<Uint8Array> {
    const messageHex = eth.hexlify(message).substring(2); // Remove 0x prefix

    try {
      const result = await this.client.send('account_signData', [
        'text/plain',
        this._address.hex(),
        messageHex,
      ]);

      // Clef returns hex string with 0x prefix
      const signature = result.startsWith('0x') ? result.substring(2) : result;
      return eth.getBytes(`0x${signature}`);
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : String(error);
      throw new Error(`Clef signing failed: ${errorMessage}`);
    }
  }

  /**
   * Sign a transaction using the EIP-155 standard
   * @param transaction The transaction to sign
   * @returns The signed transaction
   * @throws Error if signing fails
   */
  async signTransaction(transaction: Transaction): Promise<SignedTransaction> {
    // Prepare transaction data for Clef
    const args: Record<string, string> = {
      from: this._address.hex(),
      chainId: `0x${this._chainID.toString(16)}`,
    };

    // Add transaction properties
    if (transaction.to) {
      args.to = transaction.to.hex();
    }

    if (transaction.data) {
      args.data = eth.hexlify(transaction.data);
    }

    if (transaction.value !== undefined && transaction.value !== null) {
      args.value = `0x${transaction.value.toString(16)}`;
    }

    if (transaction.gas !== undefined && transaction.gas !== null) {
      args.gas = `0x${transaction.gas.toString(16)}`;
    }

    if (transaction.gasPrice !== undefined && transaction.gasPrice !== null) {
      args.gasPrice = `0x${transaction.gasPrice.toString(16)}`;
    }

    if (transaction.nonce !== undefined) {
      args.nonce = `0x${transaction.nonce.toString(16)}`;
    }

    try {
      // Call Clef to sign the transaction
      const result = await this.client.send('account_signTransaction', [args]);

      // Parse signature components
      const signedTx = result as SignedTransactionResponse;

      // Parse R, S, V values from hex strings
      const r = BigInt(signedTx.tx.r);
      const s = BigInt(signedTx.tx.s);
      const v = parseInt(
        signedTx.tx.v.startsWith('0x') ? signedTx.tx.v.substring(2) : signedTx.tx.v,
        16
      );

      // Parse serialized transaction
      const rawHex = signedTx.raw.startsWith('0x') ? signedTx.raw : `0x${signedTx.raw}`;

      return {
        ...transaction,
        r,
        s,
        v,
        serialized: rawHex,
      };
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : String(error);
      throw new Error(`Clef signing failed: ${errorMessage}`);
    }
  }

  /**
   * Verify that we can connect to the Clef server
   * @returns A promise that resolves when the connection is verified
   * @throws Error if verification fails
   * @private
   */
  private async verifyClefConnection(): Promise<void> {
    try {
      const version = await this.client.send('account_version', []);

      if (!version) {
        throw new Error('Failed to get Clef version');
      }
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : String(error);
      throw new Error(`Failed to verify Clef connection: ${errorMessage}`);
    }
  }

  /**
   * Prepare a transaction for processing
   * @param transaction The transaction to prepare
   * @returns The prepared transaction
   * @private
   */
  private prepareTransaction(transaction: Transaction): Record<string, unknown> {
    return {
      to: transaction.to ? transaction.to.ethAddress() : undefined,
      data: transaction.data ? eth.hexlify(transaction.data) : undefined,
      value: transaction.value,
      nonce: transaction.nonce || 0,
      gasLimit: transaction.gas,
      gasPrice: transaction.gasPrice,
      chainId: this._chainID,
    };
  }
}

/**
 * Interface for the response from Clef when signing a transaction.
 * Contains the raw signed transaction and signature components.
 */
interface SignedTransactionResponse {
  /** The raw signed transaction data as a hex string */
  raw: string;
  tx: {
    /** The transaction hash as a hex string */
    hash: string;
    /** The v component of the signature as a hex string */
    v: string;
    /** The r component of the signature as a hex string */
    r: string;
    /** The s component of the signature as a hex string */
    s: string;
  };
}
