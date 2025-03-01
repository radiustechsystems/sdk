"""Private key signer implementation for the Radius SDK.

This module provides a signer implementation that uses a private key for signing.
"""

from __future__ import annotations

import binascii
from typing import Dict, Union

from eth_account import Account as EthAccount
from eth_account.messages import encode_defunct

from src.auth.types import SignerClient
from src.common.address import Address
from src.common.hash import hash_from_hex
from src.common.transaction import SignedTransaction, Transaction
from src.providers.eth.utils import to_eth_transaction


class PrivateKeySigner:
    """A signer that uses an ECDSA private key for signing.
    
    This signer is used to sign transactions and messages using a private key.
    """

    def __init__(self, private_key: Union[str, bytes], client: SignerClient) -> None:
        """Initialize a new private key signer.
        
        Args:
            private_key: The private key (hex string with or without 0x prefix, or bytes)
            client: The client used to get chain information for signing
            
        Raises:
            ValueError: If the private key is invalid

        """
        self._client = client

        # Normalize the private key
        if isinstance(private_key, str):
            # Remove 0x prefix if present
            clean_key = private_key[2:] if private_key.startswith("0x") else private_key
            try:
                key_bytes = bytes.fromhex(clean_key)
            except binascii.Error as e:
                raise ValueError("Invalid private key format") from e
        else:
            key_bytes = private_key

        # Create the Ethereum account
        self._account = EthAccount.from_key(key_bytes)

        # Extract the address
        address_bytes = bytes.fromhex(self._account.address[2:])
        self._address = Address(address_bytes)

    @property
    def address(self) -> Address:
        """Get the address associated with this signer.
        
        Returns:
            The address

        """
        return self._address

    async def sign_transaction(self, tx: Transaction) -> SignedTransaction:
        """Sign a transaction.
        
        Args:
            tx: The transaction to sign
            
        Returns:
            The signed transaction
            
        Raises:
            RuntimeError: If the transaction is incomplete or cannot be signed

        """
        # Ensure the transaction has all required fields
        if tx.nonce is None:
            tx.nonce = await self._client.get_transaction_count(self.address)

        if tx.gas_price is None:
            raise RuntimeError("Transaction must have a gas price")

        if tx.gas_limit is None:
            raise RuntimeError("Transaction must have a gas limit")

        # Get the chain ID
        chain_id = await self._client.chain_id()

        # Convert to Ethereum transaction format
        eth_tx = to_eth_transaction(tx)
        eth_tx["chainId"] = chain_id

        # Sign the transaction
        signed = self._account.sign_transaction(eth_tx)

        # Create a SignedTransaction
        tx_hash = hash_from_hex(signed.hash.hex())
        raw_tx = signed.raw_transaction  # accessing as attribute 

        return SignedTransaction(
            tx_hash=tx_hash,
            raw_tx=raw_tx,
            tx=tx,
        )

    async def sign_message(self, message: bytes) -> Dict[str, bytes]:
        """Sign a message.
        
        Args:
            message: The message to sign
            
        Returns:
            A dictionary containing signature components
            
        Raises:
            RuntimeError: If unable to sign the message

        """
        # Encode the message with Ethereum message prefix
        encoded_message = encode_defunct(message)

        # Sign the message
        signed_message = self._account.sign_message(encoded_message)

        # Extract signature components
        return {
            "r": signed_message.r.to_bytes(32, byteorder="big"),
            "s": signed_message.s.to_bytes(32, byteorder="big"),
            "v": signed_message.v.to_bytes(1, byteorder="big"),
        }
