"""Private key signer implementation for the Radius SDK.

This module provides a signer implementation that uses a private key for signing.
"""

from __future__ import annotations

import binascii
from typing import Union

from eth_account import Account as EthAccount
from eth_account.messages import encode_defunct

from src.auth.types import SignerClient
from src.common.address import Address
from src.common.hash import Hash, hash_from_hex
from src.common.transaction import SignedTransaction, Transaction
from src.providers.eth.types import BigNumberish, BytesLike
from src.providers.eth.utils import to_eth_transaction


class PrivateKeySigner:
    """A signer that uses an ECDSA private key for signing."""

    def __init__(self, private_key: Union[str, bytes], chain_id: BigNumberish) -> None:
        """Initialize a new private key signer.
        
        Args:
            private_key: The private key (hex string with or without 0x prefix, or bytes)
            chain_id: The chain ID to use for signing
        """
        self._chain_id = chain_id

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

    def address(self) -> Address:
        """Get the address associated with this signer."""
        return self._address

    def chain_id(self) -> BigNumberish:
        """Get the chain ID associated with this signer."""
        return self._chain_id
        
    def hash(self, transaction: Transaction) -> Hash:
        """Hash a transaction for signing."""
        # Convert to Ethereum transaction format
        eth_tx = to_eth_transaction(transaction)
        eth_tx["chainId"] = int(self._chain_id)
        
        # Get the hash
        tx_hash = self._account._keys.keccak(eth_tx.encode('utf-8'))
        return Hash(tx_hash)

    async def sign_transaction(self, transaction: Transaction) -> SignedTransaction:
        """Sign a transaction."""
        # Ensure gas_price is set - required by Ethereum transaction format
        if transaction.gas_price is None:
            transaction.gas_price = 0  # Default to 0 gas price for Radius networks
            
        # Convert to Ethereum transaction format
        eth_tx = to_eth_transaction(transaction)
        eth_tx["chainId"] = int(self._chain_id)

        # Sign the transaction
        signed = self._account.sign_transaction(eth_tx)

        # Create a SignedTransaction
        tx_hash = hash_from_hex(signed.hash.hex())
        raw_tx = signed.raw_transaction

        return SignedTransaction(
            tx_hash=tx_hash,
            raw_tx=raw_tx,
            tx=transaction,
        )

    async def sign_message(self, message: BytesLike) -> bytes:
        """Sign a message."""
        # Ensure message is bytes
        if isinstance(message, str):
            if message.startswith('0x'):
                message_bytes = bytes.fromhex(message[2:])
            else:
                message_bytes = message.encode('utf-8')
        else:
            message_bytes = bytes(message)
            
        # Encode the message with Ethereum message prefix
        encoded_message = encode_defunct(message_bytes)

        # Sign the message
        signed_message = self._account.sign_message(encoded_message)

        # Combine r, s, v into a single signature
        r = signed_message.r.to_bytes(32, byteorder="big")
        s = signed_message.s.to_bytes(32, byteorder="big")
        v = signed_message.v.to_bytes(1, byteorder="big")
        
        return r + s + v
