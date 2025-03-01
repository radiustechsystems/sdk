"""Clef signer implementation for the Radius SDK.

This module provides a signer implementation that uses Clef for signing.
"""

from __future__ import annotations

import json
from typing import Dict

import requests

from src.auth.types import SignerClient
from src.common.address import Address
from src.common.hash import hash_from_hex
from src.common.transaction import SignedTransaction, Transaction
from src.providers.eth.utils import to_eth_transaction


class ClefSigner:
    """A signer that uses the Clef external signing service.
    
    This signer is used to sign transactions and messages using Clef.
    """

    def __init__(self, address: Address, client: SignerClient, clef_url: str) -> None:
        """Initialize a new Clef signer.
        
        Args:
            address: The address associated with this signer
            client: The client used to get chain information for signing
            clef_url: The URL of the Clef service
            
        Raises:
            RuntimeError: If unable to connect to the Clef service

        """
        self._address = address
        self._client = client
        self._clef_url = clef_url

        # Validate that we can connect to Clef
        try:
            self._post("account_list", [])
        except Exception as e:
            raise RuntimeError(f"Unable to connect to Clef at {clef_url}: {e}") from e

    @property
    def address(self) -> Address:
        """Get the address associated with this signer.
        
        Returns:
            The address

        """
        return self._address

    async def sign_transaction(self, tx: Transaction) -> SignedTransaction:
        """Sign a transaction using Clef.
        
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
        eth_tx["from"] = self.address.hex()

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
            tx=tx,
        )

    async def sign_message(self, message: bytes) -> Dict[str, bytes]:
        """Sign a message using Clef.
        
        Args:
            message: The message to sign
            
        Returns:
            A dictionary containing signature components
            
        Raises:
            RuntimeError: If unable to sign the message

        """
        # Convert message to hex
        data = "0x" + message.hex()

        # Sign the message with Clef
        response = self._post("account_signData", ["data/plain", self.address.hex(), data])

        if "error" in response:
            raise RuntimeError(f"Clef signing failed: {response['error']}")

        result = response.get("result")
        if not result or not isinstance(result, str):
            raise RuntimeError("Clef returned invalid result")

        # Extract signature components
        signature = bytes.fromhex(result[2:]) if result.startswith("0x") else bytes.fromhex(result)

        if len(signature) != 65:
            raise RuntimeError("Invalid signature length")

        r = signature[:32]
        s = signature[32:64]
        v = signature[64:65]

        return {
            "r": r,
            "s": s,
            "v": v,
        }

    def _post(self, method: str, params: list) -> Dict:
        """Send a JSON-RPC request to Clef.
        
        Args:
            method: The JSON-RPC method to call
            params: The parameters to pass to the method
            
        Returns:
            The JSON-RPC response
            
        Raises:
            RuntimeError: If the request fails

        """
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
