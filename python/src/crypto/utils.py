"""Cryptography utilities for the Radius SDK.

This module provides utilities for cryptographic operations.
"""

from __future__ import annotations

from src.common.hash import Hash
from src.providers.eth.types import BytesLike


def keccak256(data: BytesLike) -> Hash:
    """Compute the Keccak-256 hash of the input data.
    
    Args:
        data: The data to hash
        
    Returns:
        The Keccak-256 hash as a Hash object

    """
    from eth_hash.auto import keccak

    # Normalize the input to bytes
    if isinstance(data, str):
        if data.startswith("0x"):
            bytes_data = bytes.fromhex(data[2:])
        else:
            bytes_data = data.encode()
    else:
        bytes_data = bytes(data)

    # Compute the hash
    hash_bytes = keccak(bytes_data)

    return Hash(hash_bytes)


def ecdsa_recover(message: bytes, signature: bytes) -> bytes:
    """Recover the public key from a message and its signature.
    
    Args:
        message: The message that was signed
        signature: The signature (65 bytes: r, s, v)
        
    Returns:
        The recovered public key
        
    Raises:
        ValueError: If the signature is invalid or the public key cannot be recovered

    """
    from eth_account._utils.signing import to_eth_v
    from eth_keys import keys
    from eth_utils import keccak

    if len(signature) != 65:
        raise ValueError("Signature must be 65 bytes")

    r = int.from_bytes(signature[:32], byteorder="big")
    s = int.from_bytes(signature[32:64], byteorder="big")
    v = signature[64]

    # Convert to the format expected by eth_keys
    v_standard = to_eth_v(v)

    # Create the signature object
    sig = keys.Signature(vrs=(v_standard, r, s))

    # Hash the message
    message_hash = keccak(message)

    # Recover the public key
    try:
        public_key = keys.ecdsa_recover(message_hash, sig)
        return bytes(public_key)
    except Exception as e:
        raise ValueError(f"Failed to recover public key: {e}") from e
