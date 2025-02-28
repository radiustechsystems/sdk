import { BytesLike, eth } from '../providers/eth';

/**
 * Hash represents a 32-byte Keccak-256 hash used for transactions, blocks, and states
 * This class provides methods to access the hash in different formats
 */
export class Hash {
  /**
   * The internal byte representation of the hash
   * @private
   */
  private readonly data: Uint8Array;

  /**
   * Creates a new Hash with the given data
   * @param data The hash data as a BytesLike
   */
  constructor(data: BytesLike) {
    this.data = eth.getBytes(data);
  }

  /**
   * Returns the bytes of the Hash
   * @returns The byte representation of the hash
   */
  bytes(): Uint8Array {
    return this.data;
  }

  /**
   * Returns the hexadecimal string of the Hash with 0x prefix
   * @returns The hexadecimal string representation of the hash with 0x prefix
   */
  hex(): string {
    return eth.hexlify(this.data);
  }

  /**
   * Returns the hexadecimal string of the Hash without 0x prefix
   * @returns The hexadecimal string representation of the hash without 0x prefix
   */
  hexWithoutPrefix(): string {
    return this.hex().substring(2);
  }
}
