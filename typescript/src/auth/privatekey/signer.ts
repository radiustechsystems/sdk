/**
 * The privatekey package provides a Signer implementation using ECDSA private keys.
 * This is the simplest approach for signing but requires careful key management.
 */
import { Address, Hash, SignedTransaction, Transaction } from '../../common';
import { keccak256 } from '../../crypto';
import { BigNumberish, BytesLike, Wallet, eth } from '../../providers/eth';
import { Signer, SignerClient } from '../types';

/**
 * A Signer implementation that uses a private key to sign messages and transactions
 * This is the simplest way to sign messages and transactions, but it requires keeping the private key in memory.
 * For production systems with high security requirements, consider using a hardware security module or key management service.
 * @implements {Signer}
 */
export class PrivateKeySigner implements Signer {
  /**
   * The ethers.js wallet used for signing operations
   * @private
   */
  private readonly wallet: Wallet;

  /**
   * The address associated with this signer
   * @private
   */
  private readonly _address: Address;

  /**
   * The chain ID used for EIP-155 transaction signing
   * @private
   */
  private _chainID: BigNumberish = 0;

  /**
   * Creates a new PrivateKeySigner instance
   * @param key The private key as a hex string (with or without 0x prefix)
   * @param client The Radius client used to retrieve the chain ID
   * @throws Error if the private key is invalid
   */
  constructor(key: string, client: SignerClient) {
    const formattedKey = key.startsWith('0x') ? key : `0x${key}`;

    this.wallet = new eth.Wallet(formattedKey);
    this._address = new Address(this.wallet.address);

    client
      .chainID()
      .then((id) => {
        this._chainID = id;
      })
      .catch(() => {
        // Default to 0 if we can't get chain ID
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
   */
  async signMessage(message: BytesLike): Promise<Uint8Array> {
    const messageHash = keccak256([
      eth.toUtf8Bytes('\x19Ethereum Signed Message:\n'),
      eth.toUtf8Bytes(String(eth.getBytes(message).length)),
      message,
    ]);

    const signature = await this.wallet.signMessage(eth.getBytes(messageHash));
    return eth.getBytes(signature);
  }

  /**
   * Sign a transaction using the EIP-155 standard
   * @param transaction The transaction to sign
   * @returns The signed transaction
   */
  async signTransaction(transaction: Transaction): Promise<SignedTransaction> {
    const tx = {
      to: transaction.to?.ethAddress(),
      data: transaction.data ? eth.hexlify(transaction.data) : undefined,
      value: transaction.value,
      nonce: transaction.nonce,
      gasLimit: transaction.gas,
      gasPrice: transaction.gasPrice,
      chainId: this._chainID,
    };

    try {
      const unsignedTx = eth.Transaction.from(tx);
      const signature = await this.wallet.signTransaction(unsignedTx);
      const signedTx = eth.Transaction.from(signature);
      return {
        ...transaction,
        r: BigInt(signedTx.signature?.r || 0),
        s: BigInt(signedTx.signature?.s || 0),
        v: signedTx.signature?.v || 0,
        serialized: signedTx.serialized,
      };
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : String(error);
      throw new Error(`Failed to sign transaction: ${errorMessage}`);
    }
  }

  /**
   * Prepare a transaction for processing by converting it to the format expected by ethers.js
   * This handles the conversion between Radius SDK transaction format and the ethers.js format
   *
   * @param transaction The Radius transaction to prepare
   * @returns The prepared transaction in ethers.js format
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
