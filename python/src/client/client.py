"""Client implementation for the Radius SDK.

This module provides the main client for interacting with the Radius platform.
"""

from __future__ import annotations

import asyncio
import time
from typing import Any, Dict, List, Optional, TYPE_CHECKING

from src.auth.types import Signer, SignerClient
from src.client.options import ClientOption
from src.common.abi import ABI
from src.common.address import Address
from src.common.hash import Hash, hash_from_hex
from src.common.http import DefaultHttpClient, HttpClient
from src.common.receipt import Receipt
from src.common.transaction import SignedTransaction, Transaction
from src.common.utils import bytecode_from_hex
from src.providers.eth.types import BigNumberish
from src.providers.eth.utils import from_eth_receipt
from src.transport.types import Interceptor, Logf

if TYPE_CHECKING:
    from src.contracts.contract import Contract


class Client(SignerClient):
    """Main client for interacting with the Radius platform."""

    @classmethod
    async def new(cls, url: str, *opts: ClientOption) -> Client:
        """Create a new client.
        
        Args:
            url: The URL of the Radius JSON-RPC endpoint
            opts: Options for configuring the client
            
        Returns:
            A new client instance
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
        """Initialize a new client."""
        self._url = url
        self._http_client = DefaultHttpClient()
        self._logger: Optional[Logf] = None
        self._interceptor: Optional[Interceptor] = None

    def http_client(self) -> HttpClient:
        """Get the HTTP client used by this client.
        
        Returns:
            The HTTP client
        """
        return self._http_client

    async def chain_id(self) -> int:
        """Get the chain ID of the connected network.
        
        Returns:
            The chain ID
        """
        try:
            response = await self._call("eth_chainId", [])
            return int(response, 16) if isinstance(response, str) else int(response)
        except Exception as e:
            raise RuntimeError(f"Failed to get chain ID: {e}") from e
            
    async def block_number(self) -> int:
        """Get the latest block number.
        
        Returns:
            The latest block number
        """
        try:
            response = await self._call("eth_blockNumber", [])
            return int(response, 16) if isinstance(response, str) else int(response)
        except Exception as e:
            raise RuntimeError(f"Failed to get block number: {e}") from e
            
    async def gas_price(self) -> int:
        """Get the current gas price.
        
        Returns:
            The gas price in wei
        """
        try:
            response = await self._call("eth_gasPrice", [])
            return int(response, 16) if isinstance(response, str) else int(response)
        except Exception as e:
            raise RuntimeError(f"Failed to get gas price: {e}") from e

    async def balance_at(self, address: Address) -> int:
        """Get the balance of an address.
        
        Args:
            address: The address to get the balance for
            
        Returns:
            The balance in wei
        """
        try:
            response = await self._call("eth_getBalance", [address.hex(), "latest"])
            return int(response, 16) if isinstance(response, str) else int(response)
        except Exception as e:
            raise RuntimeError(f"Failed to get balance for {address.hex()}: {e}") from e

    async def pending_nonce_at(self, address: Address) -> int:
        """Get the pending nonce for an address.
        
        Args:
            address: The address to get the nonce for
            
        Returns:
            The nonce
        """
        try:
            response = await self._call("eth_getTransactionCount", [address.hex(), "pending"])
            return int(response, 16) if isinstance(response, str) else int(response)
        except Exception as e:
            raise RuntimeError(f"Failed to get nonce for {address.hex()}: {e}") from e

    async def code_at(self, address: Address) -> bytes:
        """Get the code stored at an address.
        
        Args:
            address: The address to get code for
            
        Returns:
            The contract bytecode
        """
        try:
            response = await self._call("eth_getCode", [address.hex(), "latest"])
            
            if response == "0x" or not response:
                return bytes()
                
            return bytecode_from_hex(response)
        except Exception as e:
            raise RuntimeError(f"Failed to get code at {address.hex()}: {e}") from e

    async def call(
        self, contract: "Contract", method: str, *args: Any
    ) -> List[Any]:
        """Call a read-only contract method.
        
        Args:
            contract: The contract to call
            method: The method name to call
            args: The arguments to pass to the method
            
        Returns:
            The decoded result of the call
        """
        return await contract.call(self, method, *args)

    async def execute(
        self,
        contract: "Contract",
        signer: Signer,
        method: str,
        *args: Any,
        value: int = 0
    ) -> Receipt:
        """Execute a state-changing contract method.
        
        Args:
            contract: The contract to interact with
            signer: The signer to use for sending the transaction
            method: The name of the method to call
            args: The arguments to pass to the method
            value: The amount of native currency to send with the transaction
            
        Returns:
            The transaction receipt
        """
        # Encode the function call
        data = contract.abi.encode_function_data(method, *args)
        
        # Create a transaction
        tx = Transaction(
            to=contract.address(),
            data=data,
            value=value
        )
        
        # Get nonce if not provided
        if tx.nonce is None:
            tx.nonce = await self.pending_nonce_at(signer.address())
            
        # Estimate gas if not provided
        if tx.gas is None:
            gas = await self.estimate_gas(tx)
            tx.gas = int(gas * 1.2)  # Add 20% safety margin
            
        # Set gas price if not provided
        if tx.gas_price is None:
            try:
                tx.gas_price = await self.gas_price()
            except Exception:
                # Fallback to 1 gwei if we can't get the current gas price
                tx.gas_price = 1000000000
        
        # Sign and send transaction
        signed_tx = await signer.sign_transaction(tx)
        return await self.transact(signer, signed_tx)

    async def estimate_gas(self, tx: Transaction) -> int:
        """Estimate the gas required for a transaction.
        
        Args:
            tx: The transaction to estimate gas for
            
        Returns:
            The estimated gas amount
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

    async def send(
        self,
        signer: Signer,
        recipient: Address,
        value: BigNumberish
    ) -> Receipt:
        """Send native tokens to an address.
        
        Args:
            signer: The signer to use for sending
            recipient: The recipient address
            value: The amount to send in wei
            
        Returns:
            The transaction receipt
        """
        # Create transaction
        tx = Transaction(
            to=recipient,
            value=int(value),
            data=bytes()
        )
        
        # Get nonce
        if tx.nonce is None:
            tx.nonce = await self.pending_nonce_at(signer.address())
            
        # Estimate gas if not provided
        if tx.gas is None:
            gas = await self.estimate_gas(tx)
            tx.gas = int(gas * 1.2)  # Add 20% safety margin
            
        # Set gas price if not provided
        if tx.gas_price is None:
            # Get the current gas price from the network
            try:
                tx.gas_price = await self.gas_price()
            except Exception:
                # Fallback to 1 gwei if we can't get the current gas price
                tx.gas_price = 1000000000
            
        # Sign transaction
        signed_tx = await signer.sign_transaction(tx)
            
        # Send transaction and wait for receipt
        return await self.transact(signer, signed_tx)

    async def transact(self, signer: Signer, tx: SignedTransaction) -> Receipt:
        """Send a signed transaction and wait for receipt.
        
        Args:
            signer: The signer that signed the transaction
            tx: The signed transaction
            
        Returns:
            The transaction receipt
        """
        tx_hash = await self.send_raw_transaction(tx.raw_tx)
        return await self._wait_for_transaction(tx_hash)

    async def deploy_contract(
        self, 
        signer: Signer, 
        bytecode: bytes, 
        abi: ABI,
        *args: Any,
        value: int = 0
    ) -> "Contract":
        """Deploy a contract to Radius.
        
        Args:
            signer: The signer to use for deployment
            bytecode: The contract bytecode
            abi: The contract ABI
            args: The constructor arguments
            value: Native tokens to send with deployment
            
        Returns:
            The deployed contract instance
        """
        if not bytecode:
            raise ValueError("Contract bytecode is required")
            
        # Encode constructor arguments if any
        data = bytecode
        if args and abi.constructor:
            encoded_args = abi.encode_constructor_arguments(*args)
            data = bytecode + encoded_args
            
        # Create transaction
        tx = Transaction(
            to=None,  # Contract creation
            data=data,
            value=value
        )
        
        # Get nonce
        if tx.nonce is None:
            tx.nonce = await self.pending_nonce_at(signer.address())
            
        # Estimate gas if not provided
        if tx.gas is None:
            gas = await self.estimate_gas(tx)
            tx.gas = int(gas * 1.2)  # Add 20% safety margin
            
        # Set gas price if not provided
        if tx.gas_price is None:
            try:
                tx.gas_price = await self.gas_price()
            except Exception:
                # Fallback to 1 gwei if we can't get the current gas price
                tx.gas_price = 1000000000
        
        # Sign and send transaction
        signed_tx = await signer.sign_transaction(tx)
        receipt = await self.transact(signer, signed_tx)
        
        if receipt.contract_address is None:
            raise RuntimeError("Contract deployment succeeded but no contract address was returned")
        
        # Import here to avoid circular imports
        from src.contracts.contract import Contract
        
        # Create contract instance
        return Contract(receipt.contract_address, abi)

    async def send_raw_transaction(self, raw_tx: bytes) -> Hash:
        """Send a raw transaction to Radius.
        
        Args:
            raw_tx: The raw transaction bytes
            
        Returns:
            The transaction hash
        """
        try:
            hex_tx = "0x" + raw_tx.hex()
            response = await self._call("eth_sendRawTransaction", [hex_tx])
            return hash_from_hex(response)
        except Exception as e:
            raise RuntimeError(f"Failed to send raw transaction: {e}") from e

    async def _wait_for_transaction(
        self, 
        tx_hash: Hash, 
        timeout_secs: int = 60,
        poll_interval_ms: int = 500
    ) -> Receipt:
        """Wait for a transaction to be mined.
        
        Args:
            tx_hash: Transaction hash to wait for
            timeout_secs: Maximum time to wait in seconds
            poll_interval_ms: Time between polls in milliseconds
            
        Returns:
            The transaction receipt
        """
        start_time = time.time()
        poll_interval = poll_interval_ms / 1000
        
        while True:
            receipt = await self._get_transaction_receipt(tx_hash)
            
            if receipt is not None:
                if receipt.status == 0:
                    raise RuntimeError(f"Transaction {tx_hash.hex()} failed")
                return receipt
                
            if time.time() - start_time > timeout_secs:
                raise RuntimeError(f"Timeout waiting for receipt of transaction {tx_hash.hex()}")
                
            await asyncio.sleep(poll_interval)

    async def _get_transaction_receipt(self, tx_hash: Hash) -> Optional[Receipt]:
        """Get the receipt for a transaction.
        
        Args:
            tx_hash: The transaction hash
            
        Returns:
            The transaction receipt, or None if the transaction is not yet mined
        """
        try:
            response = await self._call("eth_getTransactionReceipt", [tx_hash.hex()])

            if not response:
                return None

            return from_eth_receipt(response)
        except Exception as e:
            raise RuntimeError(f"Failed to get receipt for {tx_hash.hex()}: {e}") from e

    # Internal method for making RPC calls
    async def _call(self, method: str, params: List[Any]) -> Any:
        """Make a JSON-RPC call to the Radius endpoint."""
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
