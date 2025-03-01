"""Utility functions for the Radius SDK."""

from __future__ import annotations

import binascii
from typing import Optional


def bytecode_from_hex(hex_str: str) -> Optional[bytes]:
    """Convert a hex string to a byte slice.
    
    Args:
        hex_str: The hex string to convert (with or without 0x prefix)
        
    Returns:
        The byte representation of the hex string, or None if the input is invalid

    """
    if not hex_str:
        return None

    # Remove 0x prefix if present
    clean_hex = hex_str[2:] if hex_str.startswith("0x") else hex_str

    # Convert to bytes
    try:
        return bytes.fromhex(clean_hex)
    except (ValueError, binascii.Error):
        return None
