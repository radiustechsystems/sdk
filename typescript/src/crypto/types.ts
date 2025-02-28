/**
 * A key pair used for signing and verifying messages.
 * @property publicKey Public key
 * @property privateKey Private key
 */
export interface SigningKey {
  publicKey: Uint8Array;
  privateKey: Uint8Array;
}
