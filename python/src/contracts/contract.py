"""Contract implementation for the Radius SDK.

This module provides the Contract class for interacting with smart contracts.
"""

from __future__ import annotations

from typing import Any, Callable, Dict, List, Optional, Tuple, Union, cast

from src.auth.types import Signer
from src.common.abi import ABI
from src.common.address import Address
from src.common.hash import Hash
from src.common.receipt import Receipt
from src.common.transaction import Transaction
from src.contracts.types import ContractClient
from src.common.event import Event


class Contract:
    """Represents a smart contract on the Radius platform."""

    def __init__(self, address: Address, abi: ABI) -> None:
        """Initialize a new contract."""
        self._address = address
        self._abi = abi

    @property
    def abi(self) -> ABI:
        """Get the ABI of the contract."""
        return self._abi

    def address(self) -> Address:
        """Get the address of the contract."""
        return self._address

    async def call(
        self,
        client: ContractClient,
        method: str,
        *args: Any,
        block_identifier: str = "latest",
    ) -> List[Any]:
        """Call a read-only contract method.
        
        Args:
            client: The client to use for the call
            method: The name of the method to call
            args: The arguments to pass to the method
            block_identifier: The block to execute the call on
            
        Returns:
            The decoded result of the method call
        """
        # Encode the function call
        data = self._abi.encode_function_data(method, *args)

        # Create a transaction for the call
        tx = Transaction(to=self._address, data=data)

        # Execute the call
        result = await self._eth_call(client, tx, block_identifier)

        # Decode the result
        return self._abi.decode_function_result(method, result)

    async def execute(
        self,
        client: ContractClient,
        signer: Signer,
        method: str,
        *args: Any,
        value: int = 0
    ) -> Receipt:
        """Execute a state-changing contract method.
        
        Args:
            client: The client to use for the transaction
            signer: The signer to use for the transaction
            method: The name of the method to call
            args: The arguments to pass to the method
            value: The amount of native currency to send with the transaction
            
        Returns:
            The transaction receipt
        """
        # Encode the function call
        data = self._abi.encode_function_data(method, *args)

        # Create a transaction for the execution
        tx = Transaction(
            to=self._address,
            data=data,
            value=value
        )
        
        # Get nonce if not provided
        if tx.nonce is None:
            tx.nonce = await client.pending_nonce_at(signer.address())
            
        # Estimate gas if not provided
        if tx.gas is None:
            gas = await client.estimate_gas(tx)
            tx.gas = int(gas * 1.2)  # Add 20% safety margin
            
        # Set gas price if not provided
        if tx.gas_price is None:
            try:
                tx.gas_price = await client.gas_price()
            except Exception:
                # Fallback to 0 gas price for Radius networks
                tx.gas_price = 0
            
        # Sign transaction
        signed_tx = await signer.sign_transaction(tx)
        
        # Send transaction and wait for receipt
        return await client.transact(signer, signed_tx)

    @classmethod
    def new(cls, address: Address, abi: ABI) -> Contract:
        """Create a new contract instance.
        
        Args:
            address: The address of the deployed contract
            abi: The ABI of the contract
            
        Returns:
            A new contract instance
        """
        return cls(address, abi)

    async def _eth_call(
        self, 
        client: ContractClient, 
        tx: Transaction, 
        block_identifier: str
    ) -> bytes:
        """Execute a call to the Ethereum client.
        
        Args:
            client: The client to use for the call
            tx: The transaction to execute
            block_identifier: The block to execute the call on
            
        Returns:
            The raw result data
        """
        # Convert the transaction to a format the node understands
        call_obj: Dict[str, Any] = {}

        if tx.to is not None:
            call_obj["to"] = tx.to.hex()

        if tx.data:
            call_obj["data"] = "0x" + tx.data.hex()

        if tx.value > 0:
            call_obj["value"] = hex(tx.value)

        # Make the call
        response = await client._call("eth_call", [call_obj, block_identifier])

        # Convert the response to bytes
        if isinstance(response, str) and response.startswith("0x"):
            return bytes.fromhex(response[2:])
        elif isinstance(response, str):
            return bytes.fromhex(response)
        else:
            return bytes()
