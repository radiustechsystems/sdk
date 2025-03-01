"""Receipt module for the Radius SDK.

This module provides the Receipt class and related utilities for transaction receipts.
"""

from __future__ import annotations

from dataclasses import dataclass
from typing import Any, Dict, List, Optional

from src.common.address import Address
from src.common.event import Event
from src.common.hash import Hash


@dataclass
class Receipt:
    """Represents a transaction receipt on the Radius blockchain.
    
    Transaction receipts contain information about the execution of a transaction.
    """

    transaction_hash: Hash
    """The hash of the transaction."""

    block_hash: Hash
    """The hash of the block containing the transaction."""

    block_number: int
    """The number of the block containing the transaction."""

    contract_address: Optional[Address] = None
    """The address of the contract created, if the transaction was a contract creation."""

    from_address: Optional[Address] = None
    """The address of the sender."""

    to_address: Optional[Address] = None
    """The address of the receiver."""

    gas_used: int = 0
    """The amount of gas used by the transaction."""

    status: bool = False
    """Whether the transaction was successful."""

    logs: List[Event] = None
    """The logs generated during the transaction."""

    def __post_init__(self) -> None:
        """Initialize default values for optional fields."""
        if self.logs is None:
            self.logs = []

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> Receipt:
        """Create a Receipt from a dictionary.
        
        Args:
            data: The dictionary containing receipt data
            
        Returns:
            A Receipt instance

        """
        from src.common.address import address_from_hex
        from src.common.hash import hash_from_hex

        # Process required fields
        tx_hash = hash_from_hex(data["transactionHash"])
        block_hash = hash_from_hex(data["blockHash"])
        block_number = int(data["blockNumber"], 16) if isinstance(data["blockNumber"], str) else int(data["blockNumber"])

        # Create the receipt
        receipt = cls(
            transaction_hash=tx_hash,
            block_hash=block_hash,
            block_number=block_number,
        )

        # Process optional fields
        if "contractAddress" in data and data["contractAddress"]:
            receipt.contract_address = address_from_hex(data["contractAddress"])

        if "from" in data and data["from"]:
            receipt.from_address = address_from_hex(data["from"])

        if "to" in data and data["to"]:
            receipt.to_address = address_from_hex(data["to"])

        if "gasUsed" in data:
            receipt.gas_used = int(data["gasUsed"], 16) if isinstance(data["gasUsed"], str) else int(data["gasUsed"])

        if "status" in data:
            # Status can be "0x1", "0x0", 1, 0
            status_val = data["status"]
            if isinstance(status_val, str):
                receipt.status = status_val == "0x1"
            else:
                receipt.status = bool(status_val)

        # Process logs
        if "logs" in data and isinstance(data["logs"], list):
            receipt.logs = [Event.from_dict(log) for log in data["logs"]]

        return receipt
