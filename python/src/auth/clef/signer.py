"""Clef signer implementation for the Radius SDK.

This module provides a signer implementation that uses Clef for signing.
"""

from __future__ import annotations

import json
from typing import Dict, Union

import requests

from src.common.address import Address
from src.common.hash import Hash, hash_from_hex
from src.common.transaction import SignedTransaction, Transaction
from src.providers.eth.types import BigNumberish, BytesLike
from src.providers.eth.utils import to_eth_transaction


class ClefSigner:
    """A signer that uses the Clef external signing service."""

    def __init__(self, address: Address, chain_id: BigNumberish, clef_url: str) -> None:
        """Initialize a new Clef signer.
        
        Args:
            address: The address associated with this signer
            chain_id: The chain ID for signing
            clef_url: The URL of the Clef service
        """
        self._address = address
        self._chain_id = chain_id
        self._clef_url = clef_url

        # Validate that we can connect to Clef
        try:
            self._post("account_list", [])
        except Exception as e:
            raise RuntimeError(f"Unable to connect to Clef at {clef_url}: {e}") from e

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
        
        # Use Keccak hash function
        import hashlib
        from eth_utils import keccak
        
        # Convert transaction to bytes and hash
        tx_bytes = json.dumps(eth_tx).encode('utf-8')
        tx_hash = keccak(tx_bytes)
        
        return Hash(tx_hash)

    async def sign_transaction(self, transaction: Transaction) -> SignedTransaction:
        """Sign a transaction using Clef."""
        # Convert to Ethereum transaction format
        eth_tx = to_eth_transaction(transaction)
        eth_tx["chainId"] = int(self._chain_id)
        eth_tx["from"] = self._address.hex()

        # Sign the transaction with Clef
        response = self._post("account_signTransaction", [eth_tx])

        if "error" in response:
            raise RuntimeError(f"Clef signing failed: {response['error']}")

        result = response.get("result")
        if not result:
            raise RuntimeError("Clef returned empty result")

        # Extract the signed transaction data
        raw_tx = bytes.fromhex(result["raw"][2:]) if result["raw"].startswith("0x") else bytes.fromhex(result["raw"])
        tx_hash = hash_from_hex(result["hash"])

        return SignedTransaction(
            tx_hash=tx_hash,
            raw_tx=raw_tx,
            tx=transaction,
        )

    async def sign_message(self, message: BytesLike) -> bytes:
        """Sign a message using Clef."""
        # Ensure message is bytes
        if isinstance(message, str):
            if message.startswith('0x'):
                message_bytes = bytes.fromhex(message[2:])
            else:
                message_bytes = message.encode('utf-8')
        else:
            message_bytes = bytes(message)
            
        # Convert message to hex
        data = "0x" + message_bytes.hex()

        # Sign the message with Clef
        response = self._post("account_signData", ["data/plain", self._address.hex(), data])

        if "error" in response:
            raise RuntimeError(f"Clef signing failed: {response['error']}")

        result = response.get("result")
        if not result or not isinstance(result, str):
            raise RuntimeError("Clef returned invalid result")

        # Extract signature
        signature = bytes.fromhex(result[2:]) if result.startswith("0x") else bytes.fromhex(result)

        if len(signature) != 65:
            raise RuntimeError("Invalid signature length")

        return signature

    def _post(self, method: str, params: list) -> Dict:
        """Send a JSON-RPC request to Clef."""
        payload = {
            "id": 1,
            "jsonrpc": "2.0",
            "method": method,
            "params": params,
        }

        try:
            response = requests.post(
                self._clef_url,
                data=json.dumps(payload),
                headers={"Content-Type": "application/json"},
            )
            response.raise_for_status()
            return response.json()
        except requests.RequestException as e:
            raise RuntimeError(f"Failed to communicate with Clef: {e}") from e
        except json.JSONDecodeError as e:
            raise RuntimeError(f"Invalid JSON response from Clef: {e}") from e
