"""Authentication and signing types for the Radius SDK.

This module provides protocols for signers and related authentication components.
"""

from __future__ import annotations

from typing import Dict, Protocol

from src.common.address import Address
from src.common.transaction import SignedTransaction, Transaction


class SignerClient(Protocol):
    """Protocol for clients that can provide chain information for signing.
    
    This protocol defines methods required by signers to get chain information.
    """

    async def chain_id(self) -> int:
        """Get the chain ID of the connected blockchain.
        
        Returns:
            The chain ID as an integer
            
        Raises:
            RuntimeError: If unable to fetch the chain ID

        """
        ...

    async def get_transaction_count(self, address: Address) -> int:
        """Get the transaction count (nonce) for an address.
        
        Args:
            address: The address to get the transaction count for
            
        Returns:
            The transaction count as an integer
            
        Raises:
            RuntimeError: If unable to fetch the transaction count

        """
        ...


class Signer(Protocol):
    """Protocol for transaction and message signers.
    
    This protocol defines methods for signing transactions and messages.
    """

    @property
    def address(self) -> Address:
        """Get the address associated with this signer.
        
        Returns:
            The address

        """
        ...

    async def sign_transaction(self, tx: Transaction) -> SignedTransaction:
        """Sign a transaction.
        
        Args:
            tx: The transaction to sign
            
        Returns:
            The signed transaction
            
        Raises:
            RuntimeError: If unable to sign the transaction

        """
        ...

    async def sign_message(self, message: bytes) -> Dict[str, bytes]:
        """Sign a message.
        
        Args:
            message: The message to sign
            
        Returns:
            A dictionary containing signature components
            
        Raises:
            RuntimeError: If unable to sign the message

        """
        ...
