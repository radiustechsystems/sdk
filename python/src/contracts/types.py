"""Contract types for the Radius SDK.

This module provides type definitions for contract-related functionality.
"""

from __future__ import annotations

from typing import Any, List, Protocol, TYPE_CHECKING

from src.auth.types import Signer
from src.common.address import Address
from src.common.hash import Hash
from src.common.receipt import Receipt
from src.common.transaction import SignedTransaction, Transaction

if TYPE_CHECKING:
    from src.contracts.contract import Contract


class ContractClient(Protocol):
    """Protocol for clients that can perform contract-related operations."""

    async def pending_nonce_at(self, address: Address) -> int:
        """Get the pending nonce for an address."""
        ...

    async def estimate_gas(self, tx: Transaction) -> int:
        """Estimate the gas required for a transaction."""
        ...

    async def send_raw_transaction(self, raw_tx: bytes) -> Hash:
        """Send a raw transaction to the network."""
        ...

    async def _call(self, method: str, params: List[Any]) -> Any:
        """Make a JSON-RPC call to the node."""
        ...

    async def transact(self, signer: Signer, tx: SignedTransaction) -> Receipt:
        """Send a signed transaction and wait for receipt."""
        ...
