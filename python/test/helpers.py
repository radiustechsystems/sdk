"""Test helpers for the Radius SDK.

This module provides utility functions and classes for testing the SDK.
"""

import asyncio
import os
from typing import Any, TypeVar

from radius import (
    Account,
    Client,
    Contract,
    new_account,
    new_client,
    with_private_key,
)

T = TypeVar("T")


async def setup_test_client() -> Client:
    """Set up a test client connected to a local node.
    
    Returns:
        A connected client
        
    Raises:
        RuntimeError: If unable to connect to a test node

    """
    url = os.environ.get("RADIUS_ENDPOINT", "")

    try:
        return await new_client(url)
    except Exception as e:
        raise RuntimeError(f"Failed to connect to test node at {url}: {e}") from e


async def setup_test_account(client: Client) -> Account:
    """Set up a test account with a private key.
    
    Args:
        client: The client to use
        
    Returns:
        An initialized account
        
    Raises:
        RuntimeError: If unable to create a test account

    """
    private_key = os.environ.get("RADIUS_PRIVATE_KEY", "")

    try:
        return await new_account(with_private_key(private_key, client))
    except Exception as e:
        raise RuntimeError(f"Failed to create test account: {e}") from e


async def deploy_test_contract(
    client: Client, account: Account, abi_json: str, bytecode_hex: str, *constructor_args: Any
) -> Contract:
    """Deploy a test contract.
    
    Args:
        client: The client to use
        account: The account to deploy from
        abi_json: The contract ABI as a JSON string
        bytecode_hex: The contract bytecode as a hex string
        constructor_args: Arguments for the contract constructor
        
    Returns:
        The deployed contract
        
    Raises:
        RuntimeError: If unable to deploy the contract

    """
    from radius import abi_from_json, bytecode_from_hex

    # Parse the ABI and bytecode
    abi = abi_from_json(abi_json)
    if not abi:
        raise RuntimeError("Invalid ABI JSON")

    bytecode = bytecode_from_hex(bytecode_hex)
    if not bytecode:
        raise RuntimeError("Invalid bytecode hex")

    # Deploy the contract with gas price
    gas_price = await client.gas_price()
    tx_hash, contract = await Contract.deploy(
        client, account, abi, bytecode, *constructor_args,
        gas_price=gas_price
    )

    # Wait for the receipt
    receipt = None
    for _ in range(30):  # Wait up to 30 seconds
        receipt = await client.get_transaction_receipt(tx_hash)
        if receipt:
            break
        await asyncio.sleep(1)

    if not receipt:
        raise RuntimeError("Contract deployment timed out")

    if not receipt.status:
        raise RuntimeError("Contract deployment failed")

    if not receipt.contract_address:
        raise RuntimeError("Contract deployment did not return an address")

    # Return the deployed contract
    return Contract(receipt.contract_address, abi)
