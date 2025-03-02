"""Authentication and signing types for the Radius SDK.

This module provides protocols for signers and related authentication components.
"""

from __future__ import annotations

from typing import Union, Protocol

from src.common.address import Address
from src.common.hash import Hash
from src.common.transaction import SignedTransaction, Transaction
from src.providers.eth.types import BigNumberish, BytesLike


class SignerClient(Protocol):
    """Protocol for clients that can provide chain information for signing."""

    async def chain_id(self) -> int:
        """Get the chain ID of the connected network."""
        ...

    async def pending_nonce_at(self, address: Address) -> int:
        """Get the pending nonce for an address."""
        ...


class Signer(Protocol):
    """Protocol for transaction and message signers."""

    def address(self) -> Address:
        """Get the address associated with this signer."""
        ...

    def chain_id(self) -> BigNumberish:
        """Get the chain ID associated with this signer."""
        ...
        
    def hash(self, transaction: Transaction) -> Hash:
        """Hash a transaction for signing."""
        ...

    async def sign_transaction(self, transaction: Transaction) -> SignedTransaction:
        """Sign a transaction."""
        ...

    async def sign_message(self, message: BytesLike) -> bytes:
        """Sign a message."""
        ...
