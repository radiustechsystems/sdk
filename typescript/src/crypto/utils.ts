import { Address } from '../common';
import { BytesLike, eth } from '../providers/eth';
import { SigningKey } from './types';

/**
 * Convert a hex string private key to a SigningKey.
 * Creates an Ethereum wallet from the private key first, then extracts the signing key.
 * @param key Hex string of the private key (with or without 0x prefix)
 * @returns SigningKey containing both public and private keys
 */
export function hexToSigningKey(key: string): SigningKey {
  const formattedKey = key.startsWith('0x') ? key : `0x${key}`;
  const wallet = new eth.Wallet(formattedKey);

  return {
    publicKey: eth.getBytes(wallet.signingKey.publicKey),
    privateKey: eth.getBytes(wallet.signingKey.privateKey),
  };
}

/**
 * Calculate the Keccak256 hash of the input data.
 * This is the hashing algorithm used by Ethereum for various cryptographic operations.
 * @param data Input data as a single value or an array of values to be concatenated before hashing
 * @returns Keccak256 hash as a Uint8Array
 */
export function keccak256(data: BytesLike | BytesLike[]): Uint8Array {
  if (Array.isArray(data)) {
    return eth.getBytes(eth.keccak256(eth.concat(data.map((d) => eth.getBytes(d)))));
  }
  return eth.getBytes(eth.keccak256(eth.getBytes(data)));
}

/**
 * Convert a public key to an account address.
 * The address is derived by taking the Keccak256 hash of the public key
 * (without the prefix byte) and keeping the last 20 bytes.
 * @param publicKey Public key as BytesLike
 * @returns Account address as an Address object
 */
export function pubkeyToAddress(publicKey: BytesLike): Address {
  const bytes = eth.getBytes(publicKey);
  const hash = eth.keccak256(bytes.slice(1)); // Remove the prefix byte
  return new Address(eth.getAddress(`0x${hash.slice(-40)}`));
}

/**
 * Sign a digest hash with a signing key.
 * The signature is in the Ethereum format: [R || S || V] where V is 0 or 1.
 * @param digestHash Digest hash to sign (typically a Keccak256 hash)
 * @param key Signing key containing the private key
 * @returns The signature as a Uint8Array
 */
export function sign(digestHash: BytesLike, key: SigningKey): Uint8Array {
  const signingKey = new eth.SigningKey(eth.hexlify(key.privateKey));
  const signature = signingKey.sign(eth.getBytes(digestHash));
  return eth.getBytes(signature.serialized);
}
