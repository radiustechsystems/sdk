"""Account types for the Radius SDK.

This module provides type definitions for account-related functionality.
"""

from __future__ import annotations

from typing import Protocol

from src.common.address import Address
from src.common.hash import Hash
from src.common.transaction import Transaction


class AccountClient(Protocol):
    """Protocol for clients that can perform account-related operations.
    
    This protocol defines methods for getting balances, sending transactions, etc.
    """

    async def get_balance(self, address: Address) -> int:
        """Get the balance of an address.
        
        Args:
            address: The address to get the balance for
            
        Returns:
            The balance in wei
            
        Raises:
            RuntimeError: If unable to retrieve the balance

        """
        ...

    async def get_transaction_count(self, address: Address) -> int:
        """Get the transaction count (nonce) for an address.
        
        Args:
            address: The address to get the transaction count for
            
        Returns:
            The transaction count
            
        Raises:
            RuntimeError: If unable to retrieve the transaction count

        """
        ...

    async def send_raw_transaction(self, raw_tx: bytes) -> Hash:
        """Send a raw transaction to the blockchain.
        
        Args:
            raw_tx: The raw transaction bytes
            
        Returns:
            The transaction hash
            
        Raises:
            RuntimeError: If unable to send the transaction

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
