"""Address module for the Radius SDK.

This module provides the Address class and related utilities for handling blockchain addresses.
"""

from __future__ import annotations

import re
from typing import ClassVar


class Address:
    """Represents a Radius blockchain address.
    
    An address is a unique identifier for an account or contract on the Radius blockchain.
    """

    # Class constants
    ZERO: ClassVar[bytes] = bytes(20)
    PATTERN: ClassVar[re.Pattern] = re.compile(r"^(0x)?[0-9a-fA-F]{40}$")

    def __init__(self, address_bytes: bytes) -> None:
        """Initialize an address with its underlying bytes.
        
        Args:
            address_bytes: The underlying bytes of the address (20 bytes)
            
        Raises:
            ValueError: If address_bytes is not exactly 20 bytes long

        """
        if len(address_bytes) != 20:
            raise ValueError(f"Address must be exactly 20 bytes, got {len(address_bytes)}")
        self._bytes = address_bytes

    @property
    def bytes(self) -> bytes:
        """Get the raw bytes of the address.
        
        Returns:
            The raw bytes of the address

        """
        return self._bytes

    def hex(self) -> str:
        """Get the hex representation of the address.
        
        Returns:
            The hex representation of the address with 0x prefix

        """
        return "0x" + self._bytes.hex()
        
    def checksummed_hex(self) -> str:
        """Get the checksummed hex representation of the address.
        
        This returns the address in EIP-55 checksummed format, which is
        required by some libraries like web3.py.
        
        Returns:
            The checksummed hex representation of the address with 0x prefix

        """
        from eth_utils import to_checksum_address
        return to_checksum_address("0x" + self._bytes.hex())

    def __str__(self) -> str:
        """Get the string representation of the address.
        
        Returns:
            The hex representation of the address with 0x prefix

        """
        return self.hex()

    def __repr__(self) -> str:
        """Get the representation of the address.
        
        Returns:
            A string representation of the Address object

        """
        return f"Address({self.hex()})"

    def __eq__(self, other: object) -> bool:
        """Check if two addresses are equal.
        
        Args:
            other: The object to compare with
            
        Returns:
            True if the addresses are equal, False otherwise

        """
        if not isinstance(other, Address):
            return False
        return self._bytes == other._bytes

    def __hash__(self) -> int:
        """Get the hash of the address.
        
        Returns:
            The hash of the address bytes

        """
        return hash(self._bytes)


def address_from_hex(hex_str: str) -> Address:
    """Create an Address from a hex string.
    
    Args:
        hex_str: The hex string (with or without 0x prefix)
        
    Returns:
        An Address instance
        
    Raises:
        ValueError: If the hex string is invalid

    """
    if not Address.PATTERN.match(hex_str):
        raise ValueError(f"Invalid address format: {hex_str}")

    # Remove 0x prefix if present
    clean_hex = hex_str[2:] if hex_str.startswith("0x") else hex_str

    # Convert to bytes
    try:
        address_bytes = bytes.fromhex(clean_hex)
        return Address(address_bytes)
    except ValueError as e:
        raise ValueError(f"Invalid address hex: {hex_str}") from e


def zero_address() -> Address:
    """Create a zero address (0x0000000000000000000000000000000000000000).
    
    Returns:
        An Address instance representing the zero address

    """
    return Address(Address.ZERO)
