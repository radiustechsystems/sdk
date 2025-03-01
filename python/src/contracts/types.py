"""Contract types for the Radius SDK.

This module provides type definitions for contract-related functionality.
"""

from __future__ import annotations

from typing import Protocol

from src.common.transaction import Transaction


class ContractClient(Protocol):
    """Protocol for clients that can perform contract-related operations.
    
    This protocol defines methods for calling contract methods and sending contract transactions.
    """

    async def call(self, tx: Transaction, block_identifier: str = "latest") -> bytes:
        """Call a contract method without sending a transaction.
        
        Args:
            tx: The transaction to call
            block_identifier: The block to execute the call on
            
        Returns:
            The result of the call
            
        Raises:
            RuntimeError: If the call fails

        """
        ...

    async def estimate_gas(self, tx: Transaction) -> int:
        """Estimate the gas required for a transaction.
        
        Args:
            tx: The transaction to estimate gas for
            
        Returns:
            The estimated gas amount
            
        Raises:
            RuntimeError: If unable to estimate the gas

        """
        ...
