"""Client implementation for the Radius SDK.

This module provides the main client for interacting with the Radius platform.
"""

from __future__ import annotations

from typing import Any, Dict, List, Optional

from src.auth.types import SignerClient
from src.client.options import ClientOption
from src.common.address import Address
from src.common.hash import Hash, hash_from_hex
from src.common.http import DefaultHttpClient
from src.common.receipt import Receipt
from src.common.transaction import Transaction
from src.providers.eth.utils import from_eth_receipt


class Client(SignerClient):
    """Main client for interacting with the Radius platform.
    
    The client provides methods for querying the Radius platform and sending transactions.
    """

    @classmethod
    async def new(cls, url: str, *opts: ClientOption) -> Client:
        """Create a new client.
        
        Args:
            url: The URL of the Radius JSON-RPC endpoint
            opts: Options for configuring the client
            
        Returns:
            A new client instance
            
        Raises:
            ValueError: If the URL is invalid or the client cannot connect to Radius

        """
        client = cls(url)

        # Apply options
        for opt in opts:
            client = opt(client)

        # Test the connection
        try:
            await client.chain_id()
        except Exception as e:
            raise ValueError(f"Failed to connect to node at {url}: {e}") from e

        return client

    def __init__(self, url: str) -> None:
        """Initialize a new client.
        
        Args:
            url: The URL of the Radius JSON-RPC endpoint

        """
        self._url = url
        self._http_client = DefaultHttpClient()

    async def chain_id(self) -> int:
        """Get the chain ID of Radius.
        
        Returns:
            The chain ID
            
        Raises:
            RuntimeError: If unable to retrieve the chain ID

        """
        try:
            response = await self._call("eth_chainId", [])
            return int(response, 16) if isinstance(response, str) else int(response)
        except Exception as e:
            raise RuntimeError(f"Failed to get chain ID: {e}") from e

    async def gas_price(self) -> int:
        """Get the current gas price.
        
        Returns:
            The gas price in wei
            
        Raises:
            RuntimeError: If unable to retrieve the gas price

        """
        try:
            response = await self._call("eth_gasPrice", [])
            return int(response, 16) if isinstance(response, str) else int(response)
        except Exception as e:
            raise RuntimeError(f"Failed to get gas price: {e}") from e

    async def block_number(self) -> int:
        """Get the latest block number.
        
        Returns:
            The latest block number
            
        Raises:
            RuntimeError: If unable to retrieve the block number

        """
        try:
            response = await self._call("eth_blockNumber", [])
            return int(response, 16) if isinstance(response, str) else int(response)
        except Exception as e:
            raise RuntimeError(f"Failed to get block number: {e}") from e

    async def get_balance(self, address: Address) -> int:
        """Get the balance of an address.
        
        Args:
            address: The address to get the balance for
            
        Returns:
            The balance in wei
            
        Raises:
            RuntimeError: If unable to retrieve the balance

        """
        try:
            response = await self._call("eth_getBalance", [address.hex(), "latest"])
            return int(response, 16) if isinstance(response, str) else int(response)
        except Exception as e:
            raise RuntimeError(f"Failed to get balance for {address.hex()}: {e}") from e

    async def get_transaction_count(self, address: Address) -> int:
        """Get the transaction count (nonce) for an address.
        
        Args:
            address: The address to get the transaction count for
            
        Returns:
            The transaction count
            
        Raises:
            RuntimeError: If unable to retrieve the transaction count

        """
        try:
            response = await self._call("eth_getTransactionCount", [address.hex(), "latest"])
            return int(response, 16) if isinstance(response, str) else int(response)
        except Exception as e:
            raise RuntimeError(f"Failed to get transaction count for {address.hex()}: {e}") from e

    async def send_raw_transaction(self, raw_tx: bytes) -> Hash:
        """Send a raw transaction to Radius.
        
        Args:
            raw_tx: The raw transaction bytes
            
        Returns:
            The transaction hash
            
        Raises:
            RuntimeError: If unable to send the transaction

        """
        try:
            hex_tx = "0x" + raw_tx.hex()
            response = await self._call("eth_sendRawTransaction", [hex_tx])
            return hash_from_hex(response)
        except Exception as e:
            raise RuntimeError(f"Failed to send raw transaction: {e}") from e

    async def get_transaction_receipt(self, tx_hash: Hash) -> Optional[Receipt]:
        """Get the receipt for a transaction.
        
        Args:
            tx_hash: The transaction hash
            
        Returns:
            The transaction receipt, or None if the transaction is not yet mined
            
        Raises:
            RuntimeError: If unable to retrieve the receipt

        """
        try:
            response = await self._call("eth_getTransactionReceipt", [tx_hash.hex()])

            if not response:
                return None

            return from_eth_receipt(response)
        except Exception as e:
            raise RuntimeError(f"Failed to get receipt for {tx_hash.hex()}: {e}") from e

    async def call(
        self, tx: Transaction, block_identifier: str = "latest"
    ) -> bytes:
        """Call a contract method without sending a transaction.
        
        Args:
            tx: The transaction to call
            block_identifier: The block to execute the call on ("latest", "pending", etc.)
            
        Returns:
            The result of the call
            
        Raises:
            RuntimeError: If the call fails

        """
        try:
            # Convert the transaction to a format the node understands
            call_obj: Dict[str, Any] = {}

            if tx.to is not None:
                call_obj["to"] = tx.to.hex()

            if tx.data:
                call_obj["data"] = "0x" + tx.data.hex()

            if tx.value > 0:
                call_obj["value"] = hex(tx.value)

            # Make the call
            response = await self._call("eth_call", [call_obj, block_identifier])

            # Convert the response to bytes
            if isinstance(response, str) and response.startswith("0x"):
                return bytes.fromhex(response[2:])
            elif isinstance(response, str):
                return bytes.fromhex(response)
            else:
                return bytes([])
        except Exception as e:
            raise RuntimeError(f"Contract call failed: {e}") from e

    async def estimate_gas(self, tx: Transaction) -> int:
        """Estimate the gas required for a transaction.
        
        Args:
            tx: The transaction to estimate gas for
            
        Returns:
            The estimated gas amount
            
        Raises:
            RuntimeError: If unable to estimate the gas

        """
        try:
            # Convert the transaction to a format the node understands
            call_obj: Dict[str, Any] = {}

            if tx.to is not None:
                call_obj["to"] = tx.to.hex()

            if tx.data:
                call_obj["data"] = "0x" + tx.data.hex()

            if tx.value > 0:
                call_obj["value"] = hex(tx.value)

            # Make the call
            response = await self._call("eth_estimateGas", [call_obj])

            # Convert the response to an integer
            return int(response, 16) if isinstance(response, str) else int(response)
        except Exception as e:
            raise RuntimeError(f"Gas estimation failed: {e}") from e

    async def _call(self, method: str, params: List[Any]) -> Any:
        """Make a JSON-RPC call to the Radius endpoint.
        
        Args:
            method: The JSON-RPC method to call
            params: The parameters for the method
            
        Returns:
            The result of the call
            
        Raises:
            RuntimeError: If the call fails

        """
        request = {
            "jsonrpc": "2.0",
            "id": 1,
            "method": method,
            "params": params,
        }

        response = self._http_client.post(self._url, request)

        if "error" in response:
            error = response["error"]
            message = error.get("message", "Unknown error")
            raise RuntimeError(f"JSON-RPC error: {message}")

        if "result" not in response:
            raise RuntimeError("JSON-RPC response missing 'result' field")

        return response["result"]
