"""Transaction module for the Radius SDK.

This module provides the Transaction and SignedTransaction classes and related utilities.
"""

from __future__ import annotations

from dataclasses import dataclass
from typing import Any, Dict, Optional

from src.common.address import Address
from src.common.hash import Hash


@dataclass
class Transaction:
    """Represents a transaction on the Radius platform.
    
    Transactions are used to transfer value, deploy contracts, or execute contract methods.
    """

    to: Optional[Address] = None
    """The recipient of the transaction."""

    data: bytes = b""
    """The data payload of the transaction."""

    value: int = 0
    """The amount of native currency to send with the transaction."""

    nonce: Optional[int] = None
    """The nonce of the transaction."""

    gas_price: Optional[int] = None
    """The gas price for the transaction."""

    gas_limit: Optional[int] = None
    """The gas limit for the transaction."""

    def to_dict(self) -> Dict[str, Any]:
        """Convert the transaction to a dictionary.
        
        Returns:
            The transaction as a dictionary

        """
        result: Dict[str, Any] = {}

        if self.to is not None:
            result["to"] = self.to.hex()

        if self.data:
            result["data"] = "0x" + self.data.hex()

        if self.value > 0:
            result["value"] = self.value

        if self.nonce is not None:
            result["nonce"] = self.nonce

        if self.gas_price is not None:
            result["gasPrice"] = self.gas_price

        if self.gas_limit is not None:
            result["gas"] = self.gas_limit

        return result

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> Transaction:
        """Create a Transaction from a dictionary.
        
        Args:
            data: The dictionary containing transaction data
            
        Returns:
            A Transaction instance

        """
        tx = cls()

        if "to" in data and data["to"]:
            from src.common.address import address_from_hex
            tx.to = address_from_hex(data["to"])

        if "data" in data and data["data"]:
            data_hex = data["data"]
            if isinstance(data_hex, str) and data_hex.startswith("0x"):
                tx.data = bytes.fromhex(data_hex[2:])
            elif isinstance(data_hex, bytes):
                tx.data = data_hex

        if "value" in data:
            tx.value = int(data["value"])

        if "nonce" in data:
            tx.nonce = int(data["nonce"])

        if "gasPrice" in data:
            tx.gas_price = int(data["gasPrice"])

        if "gas" in data:
            tx.gas_limit = int(data["gas"])

        return tx


@dataclass
class SignedTransaction:
    """Represents a signed transaction on the Radius platform.
    
    A signed transaction includes the original transaction data plus signature information.
    """

    tx_hash: Hash
    """The hash of the transaction."""

    raw_tx: bytes
    """The raw signed transaction data."""

    tx: Transaction
    """The original transaction."""
