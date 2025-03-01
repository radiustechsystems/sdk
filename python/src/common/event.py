"""Event module for the Radius SDK.

This module provides the Event class and related utilities for blockchain events.
"""

from __future__ import annotations

from dataclasses import dataclass
from typing import Any, Dict, List

from src.common.address import Address
from src.common.hash import Hash


@dataclass
class Event:
    """Represents an event emitted during a transaction on the Radius blockchain.
    
    Events are emitted by smart contracts and contain information about state changes.
    """

    address: Address
    """The address of the contract that emitted the event."""

    topics: List[bytes]
    """The topics of the event."""

    data: bytes
    """The data of the event."""

    block_number: int
    """The number of the block containing the event."""

    transaction_hash: Hash
    """The hash of the transaction that emitted the event."""

    transaction_index: int
    """The index of the transaction in the block."""

    block_hash: Hash
    """The hash of the block containing the event."""

    log_index: int
    """The index of the log in the block."""

    removed: bool = False
    """Whether the log was removed due to a chain reorganization."""

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> Event:
        """Create an Event from a dictionary.
        
        Args:
            data: The dictionary containing event data
            
        Returns:
            An Event instance

        """
        from src.common.address import address_from_hex
        from src.common.hash import hash_from_hex

        # Process address
        address = address_from_hex(data["address"])

        # Process topics
        topics: List[bytes] = []
        for topic in data.get("topics", []):
            if isinstance(topic, str) and topic.startswith("0x"):
                topics.append(bytes.fromhex(topic[2:]))
            elif isinstance(topic, bytes):
                topics.append(topic)

        # Process data
        event_data = b""
        if "data" in data and data["data"] and data["data"] != "0x":
            data_hex = data["data"]
            if isinstance(data_hex, str) and data_hex.startswith("0x"):
                event_data = bytes.fromhex(data_hex[2:])
            elif isinstance(data_hex, bytes):
                event_data = data_hex

        # Process block number
        block_number = int(data["blockNumber"], 16) if isinstance(data["blockNumber"], str) else int(data["blockNumber"])

        # Process transaction hash
        tx_hash = hash_from_hex(data["transactionHash"])

        # Process transaction index
        tx_index = int(data["transactionIndex"], 16) if isinstance(data["transactionIndex"], str) else int(data["transactionIndex"])

        # Process block hash
        block_hash = hash_from_hex(data["blockHash"])

        # Process log index
        log_index = int(data["logIndex"], 16) if isinstance(data["logIndex"], str) else int(data["logIndex"])

        # Process removed flag
        removed = bool(data.get("removed", False))

        return cls(
            address=address,
            topics=topics,
            data=event_data,
            block_number=block_number,
            transaction_hash=tx_hash,
            transaction_index=tx_index,
            block_hash=block_hash,
            log_index=log_index,
            removed=removed,
        )
