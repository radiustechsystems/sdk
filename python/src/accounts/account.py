"""Account implementation for the Radius SDK.

This module provides the Account class for interacting with accounts on the Radius platform.
"""

from __future__ import annotations

from typing import Optional

from src.accounts.options import AccountOption
from src.accounts.types import AccountClient
from src.auth.types import Signer
from src.common.address import Address
from src.common.receipt import Receipt
from src.common.transaction import SignedTransaction, Transaction
from src.providers.eth.types import BigNumberish, BytesLike


class Account:
    """Represents an account on the Radius platform."""

    @classmethod
    async def new(cls, *opts: AccountOption) -> Account:
        """Create a new account with the given options.
        
        Args:
            opts: Options for configuring the account
            
        Returns:
            A new account instance
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
    def signer(self) -> Signer:
        """Get the signer for the account."""
        if not self._signer:
            raise RuntimeError("Account has no signer")
        return self._signer

    def address(self) -> Address:
        """Get the address of the account."""
        if not self._signer:
            raise RuntimeError("Account has no signer")
        return self._signer.address()

    async def balance(self, client: AccountClient) -> int:
        """Get the balance of the account.
        
        Args:
            client: The client to use for the operation
            
        Returns:
            The account balance in wei
        """
        return await client.balance_at(self.address())

    async def nonce(self, client: AccountClient) -> int:
        """Get the next nonce for this account.

        Args:
            client: The client to use for the operation

        Returns:
            The next nonce to use
        """
        return await client.pending_nonce_at(self.address())

    async def send(
        self,
        client: AccountClient,
        recipient: Address,
        value: BigNumberish
    ) -> Receipt:
        """Send native tokens to an address.

        Args:
            client: The client to use for the operation
            recipient: The recipient address
            value: The amount to send in wei

        Returns:
            The transaction receipt
        """
        return await client.send(self.signer, recipient, value)

    async def sign_message(self, message: BytesLike) -> bytes:
        """Sign a message.

        Args:
            message: The message to sign

        Returns:
            The signature as bytes
        """
        signature = await self.signer.sign_message(message)

        # Convert signature dict to bytes if needed
        if isinstance(signature, dict):
            # Convert from r, s, v format to bytes
            r = signature.get('r', b'')
            s = signature.get('s', b'')
            v = signature.get('v', b'')

            if isinstance(r, int):
                r = r.to_bytes(32, 'big')
            if isinstance(s, int):
                s = s.to_bytes(32, 'big')
            if isinstance(v, int):
                v = v.to_bytes(1, 'big')

            return r + s + v

        # If it's already bytes, return as is
        return signature if isinstance(signature, bytes) else bytes(signature)

    async def sign_transaction(self, transaction: Transaction) -> SignedTransaction:
        """Sign a transaction.

        Args:
            transaction: The transaction to sign

        Returns:
            The signed transaction
        """
        return await self.signer.sign_transaction(transaction)
