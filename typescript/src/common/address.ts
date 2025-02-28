import { BytesLike, eth } from '../providers/eth';

/**
 * Represents a 20-byte Radius account or contract address.
 *
 * This class provides methods to convert between different address representations
 * and compare addresses. It serves as the core data structure for identifying
 * accounts and smart contracts in the Radius system.
 */
export class Address {
  /**
   * The address data as a byte array
   * @private
   */
  private readonly data: Uint8Array;

  /**
   * Creates a new Address instance from various input formats.
   *
   * @param data Address data as Uint8Array, BytesLike, hex string, or another Address instance
   * @throws Error if the address is not exactly 20 bytes long
   */
  constructor(data: Address | BytesLike | string) {
    if (data instanceof Address) {
      this.data = data.bytes();
    } else if (typeof data === 'string') {
      const cleanHex = data.startsWith('0x') ? data : `0x${data}`;
      this.data = eth.getBytes(cleanHex);
    } else {
      const bytes = eth.getBytes(data);
      if (bytes.length !== 20) {
        throw new Error('Address must be 20 bytes');
      }
      this.data = bytes;
    }
  }

  /**
   * Returns the address as a byte array.
   *
   * @returns Byte array representation of the 20-byte address
   */
  bytes(): Uint8Array {
    return this.data;
  }

  /**
   * Converts a Radius Address to an Ethereum address format.
   * This method is used when Ethereum library functionality is needed.
   *
   * @returns Checksummed Ethereum address string
   */
  ethAddress(): string {
    return eth.getAddress(this.hex());
  }

  /**
   * Returns the hexadecimal string representation of the address.
   *
   * @returns Hex string representation of the address with 0x prefix
   */
  hex(): string {
    return eth.hexlify(this.data);
  }

  /**
   * Compares this address with another address for equality.
   *
   * @param other Address to compare with this address
   * @returns True if addresses are equal (case-insensitive comparison), false otherwise
   */
  equals(other: Address): boolean {
    return this.hex().toLowerCase() === other.hex().toLowerCase();
  }
}
