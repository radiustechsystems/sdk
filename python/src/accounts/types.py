"""Account types for the Radius SDK.

This module provides type definitions for account-related functionality.
"""

from __future__ import annotations

from typing import Protocol

from src.common.address import Address
from src.common.receipt import Receipt


class AccountClient(Protocol):
    """Protocol for clients that can perform account-related operations."""

    async def balance_at(self, address: Address) -> int:
        """Get the balance of an address."""
        ...

    async def pending_nonce_at(self, address: Address) -> int:
        """Get the pending nonce for an address."""
        ...
        
    async def send(self, signer, recipient: Address, value: int) -> Receipt:
        """Send native tokens to an address."""
        ...
