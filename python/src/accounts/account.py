"""Account implementation for the Radius SDK.

This module provides the Account class for interacting with accounts on the Radius platform.
"""

from __future__ import annotations

from typing import Optional

from src.accounts.options import AccountOption
from src.accounts.types import AccountClient
from src.auth.types import Signer
from src.common.address import Address
from src.common.constants import MAX_GAS
from src.common.hash import Hash
from src.common.transaction import SignedTransaction, Transaction


class Account:
    """Represents an account on the Radius platform.
    
    An account is used to manage an address, check balances, and send transactions.
    """

    @classmethod
    async def new(cls, *opts: AccountOption) -> Account:
        """Create a new account with the given options.
        
        Args:
            opts: Options for configuring the account
            
        Returns:
            A new account instance
            
        Raises:
            ValueError: If required options are missing

        """
        account = cls()

        # Apply options
        for opt in opts:
            account = opt(account)

        if not account._signer:
            raise ValueError("Account requires a signer (use with_private_key or with_signer)")

        return account

    def __init__(self) -> None:
        """Initialize a new account."""
        self._signer: Optional[Signer] = None

    @property
    def address(self) -> Address:
        """Get the address of the account.
        
        Returns:
            The account address
            
        Raises:
            RuntimeError: If the account is not properly initialized

        """
        if not self._signer:
            raise RuntimeError("Account has no signer")
        return self._signer.address

    @property
    def signer(self) -> Signer:
        """Get the signer for the account.
        
        Returns:
            The account signer
            
        Raises:
            RuntimeError: If the account is not properly initialized

        """
        if not self._signer:
            raise RuntimeError("Account has no signer")
        return self._signer

    async def get_balance(self, client: AccountClient) -> int:
        """Get the balance of the account.
        
        Args:
            client: The client to use for the operation
            
        Returns:
            The account balance in wei
            
        Raises:
            RuntimeError: If unable to retrieve the balance

        """
        return await client.get_balance(self.address)

    async def sign_transaction(self, tx: Transaction) -> SignedTransaction:
        """Sign a transaction.
        
        Args:
            tx: The transaction to sign
            
        Returns:
            The signed transaction
            
        Raises:
            RuntimeError: If unable to sign the transaction

        """
        return await self.signer.sign_transaction(tx)

    async def sign_message(self, message: bytes) -> dict:
        """Sign a message.
        
        Args:
            message: The message to sign
            
        Returns:
            The signature components
            
        Raises:
            RuntimeError: If unable to sign the message

        """
        return await self.signer.sign_message(message)

    async def send_transaction(
        self,
        client: AccountClient,
        tx: Transaction,
        wait_for_receipt: bool = False,
    ) -> Hash:
        """Send a transaction from this account.
        
        Args:
            client: The client to use for the operation
            tx: The transaction to send
            wait_for_receipt: Whether to wait for the transaction receipt
            
        Returns:
            The transaction hash
            
        Raises:
            RuntimeError: If unable to send the transaction

        """
        # Complete the transaction if needed
        if tx.nonce is None:
            tx.nonce = await client.get_transaction_count(self.address)

        if tx.gas_limit is None:
            # Estimate gas and add a 10% buffer
            estimate = await client.estimate_gas(tx)
            tx.gas_limit = min(MAX_GAS, int(estimate * 1.1))

        # Sign the transaction
        signed_tx = await self.sign_transaction(tx)

        # Send the transaction
        return await client.send_raw_transaction(signed_tx.raw_tx)
