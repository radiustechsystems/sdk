"""Ethereum utility functions for the Radius SDK.

This module provides utility functions for working with Ethereum data types and formats.
"""

from __future__ import annotations

from typing import Any, Dict, Union

from src.common.receipt import Receipt
from src.common.transaction import Transaction


def to_hex(value: Union[int, bytes, str]) -> str:
    """Convert a value to a hex string.
    
    Args:
        value: The value to convert
        
    Returns:
        The value as a hex string with 0x prefix

    """
    if isinstance(value, int):
        return hex(value)
    elif isinstance(value, bytes):
        return "0x" + value.hex()
    elif isinstance(value, str):
        if value.startswith("0x"):
            return value
        try:
            return "0x" + bytes.fromhex(value).hex()
        except ValueError:
            return "0x" + value.encode().hex()
    else:
        raise TypeError(f"Cannot convert {type(value)} to hex")


def from_hex(hex_str: str) -> bytes:
    """Convert a hex string to bytes.
    
    Args:
        hex_str: The hex string to convert
        
    Returns:
        The bytes representation of the hex string
        
    Raises:
        ValueError: If the hex string is invalid

    """
    if not hex_str.startswith("0x"):
        raise ValueError(f"Expected hex string to start with 0x, got {hex_str}")

    try:
        return bytes.fromhex(hex_str[2:])
    except ValueError as e:
        raise ValueError(f"Invalid hex string: {hex_str}") from e


def to_eth_transaction(tx: Transaction) -> Dict[str, Any]:
    """Convert a Radius Transaction to an Ethereum JSON-RPC transaction format.
    
    Args:
        tx: The Radius transaction to convert
        
    Returns:
        The Ethereum JSON-RPC transaction

    """
    result: Dict[str, Any] = {}

    if tx.to is not None:
        # Convert to checksummed Ethereum address format required by eth_account
        result["to"] = tx.to.checksummed_hex()

    if tx.data:
        result["data"] = to_hex(tx.data)

    if tx.value > 0:
        result["value"] = to_hex(tx.value)

    if tx.nonce is not None:
        result["nonce"] = to_hex(tx.nonce)

    if tx.gas_price is not None:
        result["gasPrice"] = to_hex(tx.gas_price)

    if tx.gas_limit is not None:
        result["gas"] = to_hex(tx.gas_limit)

    return result


def from_eth_receipt(eth_receipt: Dict[str, Any]) -> Receipt:
    """Convert an Ethereum JSON-RPC receipt to a Radius Receipt.
    
    Args:
        eth_receipt: The Ethereum receipt to convert
        
    Returns:
        The Radius Receipt

    """
    return Receipt.from_dict(eth_receipt)
