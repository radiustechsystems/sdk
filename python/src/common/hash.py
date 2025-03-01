"""Hash module for the Radius SDK.

This module provides the Hash class and related utilities for working with hashes.
"""

from __future__ import annotations

import re
from typing import ClassVar


class Hash:
    """Represents a hash in the Radius blockchain.
    
    A hash is a unique identifier generated from data using a cryptographic hashing function.
    """

    # Class constants
    ZERO: ClassVar[bytes] = bytes(32)
    PATTERN: ClassVar[re.Pattern] = re.compile(r"^(0x)?[0-9a-fA-F]{64}$")

    def __init__(self, hash_bytes: bytes) -> None:
        """Initialize a hash with its underlying bytes.
        
        Args:
            hash_bytes: The underlying bytes of the hash (32 bytes)
            
        Raises:
            ValueError: If hash_bytes is not exactly 32 bytes long

        """
        if len(hash_bytes) != 32:
            raise ValueError(f"Hash must be exactly 32 bytes, got {len(hash_bytes)}")
        self._bytes = hash_bytes

    @property
    def bytes(self) -> bytes:
        """Get the raw bytes of the hash.
        
        Returns:
            The raw bytes of the hash

        """
        return self._bytes

    def hex(self) -> str:
        """Get the hex representation of the hash.
        
        Returns:
            The hex representation of the hash with 0x prefix

        """
        return "0x" + self._bytes.hex()

    def __str__(self) -> str:
        """Get the string representation of the hash.
        
        Returns:
            The hex representation of the hash with 0x prefix

        """
        return self.hex()

    def __repr__(self) -> str:
        """Get the representation of the hash.
        
        Returns:
            A string representation of the Hash object

        """
        return f"Hash({self.hex()})"

    def __eq__(self, other: object) -> bool:
        """Check if two hashes are equal.
        
        Args:
            other: The object to compare with
            
        Returns:
            True if the hashes are equal, False otherwise

        """
        if not isinstance(other, Hash):
            return False
        return self._bytes == other._bytes

    def __hash__(self) -> int:
        """Get the hash of the hash object.
        
        Returns:
            The hash of the hash bytes

        """
        return hash(self._bytes)


def hash_from_hex(hex_str: str) -> Hash:
    """Create a Hash from a hex string.
    
    Args:
        hex_str: The hex string (with or without 0x prefix)
        
    Returns:
        A Hash instance
        
    Raises:
        ValueError: If the hex string is invalid

    """
    if not Hash.PATTERN.match(hex_str):
        raise ValueError(f"Invalid hash format: {hex_str}")

    # Remove 0x prefix if present
    clean_hex = hex_str[2:] if hex_str.startswith("0x") else hex_str

    # Convert to bytes
    try:
        hash_bytes = bytes.fromhex(clean_hex)
        return Hash(hash_bytes)
    except ValueError as e:
        raise ValueError(f"Invalid hash hex: {hex_str}") from e
