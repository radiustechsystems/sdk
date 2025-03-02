"""Contract implementation for the Radius SDK.

This module provides the Contract class for interacting with smart contracts.
"""

from __future__ import annotations

from typing import Any, List, Optional, Tuple

from src.accounts.account import Account
from src.common.abi import ABI
from src.common.address import Address
from src.common.hash import Hash
from src.common.transaction import Transaction
from src.contracts.types import ContractClient


class Contract:
    """Represents a smart contract on the Radius platform.
    
    The Contract class is used to call methods and send transactions to smart contracts.
    """

    def __init__(self, address: Address, abi: ABI) -> None:
        """Initialize a new contract.
        
        Args:
            address: The address of the deployed contract
            abi: The ABI of the contract

        """
        self._address = address
        self._abi = abi

    @property
    def address(self) -> Address:
        """Get the address of the contract.
        
        Returns:
            The contract address

        """
        return self._address

    @property
    def abi(self) -> ABI:
        """Get the ABI of the contract.
        
        Returns:
            The contract ABI

        """
        return self._abi

    async def call(
        self,
        client: ContractClient,
        method_name: str,
        *args: Any,
        block_identifier: str = "latest",
    ) -> List[Any]:
        """Call a read-only contract method.
        
        Args:
            client: The client to use for the call
            method_name: The name of the method to call
            args: The arguments to pass to the method
            block_identifier: The block to execute the call on
            
        Returns:
            The decoded result of the method call
            
        Raises:
            ValueError: If the method is not found in the ABI
            RuntimeError: If the call fails

        """
        # Encode the function call
        data = self._abi.encode_function_data(method_name, *args)

        # Create a transaction for the call
        tx = Transaction(to=self._address, data=data)

        # Execute the call
        result = await client.call(tx, block_identifier)

        # Decode the result
        return self._abi.decode_function_result(method_name, result)

    async def execute(
        self,
        client: ContractClient,
        account: Account,
        method_name: str,
        *args: Any,
        value: int = 0,
        gas_limit: Optional[int] = None,
        gas_price: Optional[int] = None,
    ) -> Hash:
        """Execute a state-changing contract method.
        
        Args:
            client: The client to use for the transaction
            account: The account to send the transaction from
            method_name: The name of the method to call
            args: The arguments to pass to the method
            value: The amount of native currency to send with the transaction
            gas_limit: The gas limit for the transaction
            gas_price: The gas price for the transaction
            
        Returns:
            The transaction hash
            
        Raises:
            ValueError: If the method is not found in the ABI
            RuntimeError: If the transaction fails

        """
        # Encode the function call
        data = self._abi.encode_function_data(method_name, *args)

        # Create a transaction for the execution
        tx = Transaction(
            to=self._address,
            data=data,
            value=value,
            gas_limit=gas_limit,
            gas_price=gas_price,
        )

        # Send the transaction
        return await account.send_transaction(client, tx)

    @classmethod
    async def deploy(
        cls,
        client: ContractClient,
        account: Account,
        abi: ABI,
        bytecode: bytes,
        *constructor_args: Any,
        value: int = 0,
        gas_limit: Optional[int] = None,
        gas_price: Optional[int] = None,
    ) -> Tuple[Hash, Contract]:
        """Deploy a new contract.
        
        Args:
            client: The client to use for the transaction
            account: The account to send the transaction from
            abi: The ABI of the contract
            bytecode: The bytecode of the contract
            constructor_args: The arguments to pass to the constructor
            value: The amount of native currency to send with the transaction
            gas_limit: The gas limit for the transaction
            gas_price: The gas price for the transaction
            
        Returns:
            A tuple of (transaction hash, Contract instance)
            
        Raises:
            ValueError: If the constructor is not found in the ABI
            RuntimeError: If the transaction fails

        """
        # Prepare the deployment data
        deploy_data = bytecode

        # If there are constructor arguments, encode them
        if constructor_args:
            # Find the constructor in the ABI
            constructor = None
            for item in abi._raw_abi:
                if item.get("type") == "constructor":
                    constructor = item
                    break

            if not constructor:
                if constructor_args:
                    raise ValueError("ABI does not contain a constructor, but arguments were provided")
            else:
                try:
                    # Use web3.py to encode constructor arguments
                    from web3 import Web3
                    
                    # Initialize Web3 with a null provider - we only need encoding functionality
                    w3 = Web3()
                    
                    # Create contract factory
                    contract_factory = w3.eth.contract(abi=abi._raw_abi, bytecode="0x" + bytecode.hex())
                    
                    # Get the data with encoded constructor arguments
                    constructor_instance = contract_factory.constructor(*constructor_args)
                    deploy_data_hex = constructor_instance.data_in_transaction
                    
                    # The data includes the bytecode plus encoded arguments
                    if deploy_data_hex.startswith('0x'):
                        deploy_data_hex = deploy_data_hex[2:]  # Remove 0x prefix
                    
                    deploy_data = bytes.fromhex(deploy_data_hex)
                except Exception as e:
                    raise ValueError(f"Failed to encode constructor arguments: {e}") from e

        # Create a transaction for the deployment
        tx = Transaction(
            data=deploy_data,
            value=value,
            gas_limit=gas_limit,
            gas_price=gas_price,
        )

        # Send the transaction
        tx_hash = await account.send_transaction(client, tx)

        # Return the transaction hash and a placeholder contract
        # The actual contract address will be available in the transaction receipt
        return tx_hash, cls(Address(bytes(20)), abi)
