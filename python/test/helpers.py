"""Test helpers for the Radius SDK.

This module provides utility functions and classes for testing the SDK.
"""

import os
import secrets
from typing import Any, Optional

from radius import (
    ABI,
    Account,
    Client,
    Contract,
    abi_from_json,
    bytecode_from_hex,
    new_account,
    new_client,
    with_logger,
    with_private_key,
)

# Constants for testing
MIN_TEST_ACCOUNT_BALANCE = 1000000000000000000  # 1 ETH
SIMPLE_STORAGE_ABI = """[{"inputs":[],"name":"get","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"x","type":"uint256"}],"name":"set","outputs":[],"stateMutability":"nonpayable","type":"function"}]"""
SIMPLE_STORAGE_BIN = "6080604052348015600e575f5ffd5b5060a580601a5f395ff3fe6080604052348015600e575f5ffd5b50600436106030575f3560e01c806360fe47b11460345780636d4ce63c146045575b5f5ffd5b6043603f3660046059565b5f55565b005b5f5460405190815260200160405180910390f35b5f602082840312156068575f5ffd5b503591905056fea26469706673582212207655d86666fa8aa75666db8416e0f5db680914358a57e84aa369d9250218247f64736f6c634300081c0033"


async def skip_if_no_endpoint() -> Optional[str]:
    """Skip the test if the RADIUS_ENDPOINT environment variable is not set.
    
    Returns:
        The endpoint URL or None if not set
    """
    endpoint = os.environ.get("RADIUS_ENDPOINT", "")
    if not endpoint:
        print("RADIUS_ENDPOINT environment variable not set")
        return None
    return endpoint


async def skip_if_no_private_key() -> Optional[str]:
    """Skip the test if the RADIUS_PRIVATE_KEY environment variable is not set.
    
    Returns:
        The private key or None if not set
    """
    private_key = os.environ.get("RADIUS_PRIVATE_KEY", "")
    if not private_key:
        print("RADIUS_PRIVATE_KEY environment variable not set")
        return None
    return private_key


async def create_test_client(endpoint: str = None) -> Optional[Client]:
    """Create a test client for interacting with Radius.
    
    Args:
        endpoint: The URL of the Radius JSON-RPC endpoint
        
    Returns:
        A connected client or None if creation failed
    """
    if not endpoint:
        endpoint = await skip_if_no_endpoint()
        if not endpoint:
            return None

    try:
        return await new_client(endpoint, with_logger(print))
    except Exception as e:
        print(f"Failed to create test client: {e}")
        return None


async def create_test_account(client: Client) -> Optional[Account]:
    """Create a test account with a random private key.
    
    Args:
        client: The client to use
        
    Returns:
        A new account or None if creation failed
    """
    try:
        # Generate a random private key (32 bytes)
        private_key_bytes = secrets.token_bytes(32)
        private_key_hex = private_key_bytes.hex()
        
        chain_id = await client.chain_id()
        return await new_account(with_private_key(private_key_hex, chain_id))
    except Exception as e:
        print(f"Failed to create test account: {e}")
        return None


async def get_funded_account(client: Client, private_key: str = None) -> Optional[Account]:
    """Get a funded account from environment variables.
    
    Args:
        client: The client to use
        private_key: Optional private key to use
        
    Returns:
        A funded account or None if not available
    """
    if not private_key:
        private_key = await skip_if_no_private_key()
        if not private_key:
            return None

    try:
        chain_id = await client.chain_id()
        account = await new_account(with_private_key(private_key, chain_id))
        
        # Check if the account has sufficient balance
        balance = await account.balance(client)
        if balance < MIN_TEST_ACCOUNT_BALANCE:
            print(f"Insufficient account balance: {balance}")
            return None
            
        return account
    except Exception as e:
        print(f"Failed to create funded account: {e}")
        return None


async def skip_if_insufficient_balance(account: Account, client: Client) -> Optional[int]:
    """Skip if the account has insufficient balance.
    
    Args:
        account: The account to check
        client: The client to use
        
    Returns:
        The account balance or None if insufficient
    """
    try:
        balance = await account.balance(client)
        if balance < MIN_TEST_ACCOUNT_BALANCE:
            print(f"Insufficient account balance: {balance}")
            return None
        return balance
    except Exception as e:
        print(f"Failed to get account balance: {e}")
        return None


def get_test_abi(abi_string: str = SIMPLE_STORAGE_ABI) -> ABI:
    """Get a test ABI from a string.
    
    Args:
        abi_string: The ABI JSON string
        
    Returns:
        An ABI object
    """
    abi = abi_from_json(abi_string)
    if not abi:
        raise RuntimeError("Invalid ABI JSON")
    return abi


def get_test_bytecode(bytecode_hex: str = SIMPLE_STORAGE_BIN) -> bytes:
    """Get test bytecode from a hex string.
    
    Args:
        bytecode_hex: The hex string
        
    Returns:
        The bytecode as bytes
    """
    bytecode = bytecode_from_hex(bytecode_hex)
    if not bytecode:
        raise RuntimeError("Invalid bytecode hex")
    return bytecode


async def deploy_test_contract(
    client: Client, 
    account: Account, 
    abi_json: str = SIMPLE_STORAGE_ABI, 
    bytecode_hex: str = SIMPLE_STORAGE_BIN, 
    *constructor_args: Any
) -> Optional[Contract]:
    """Deploy a test contract.
    
    Args:
        client: The client to use
        account: The account to deploy from
        abi_json: The contract ABI as a JSON string
        bytecode_hex: The contract bytecode as a hex string
        constructor_args: Arguments for the contract constructor
        
    Returns:
        The deployed contract or None if deployment failed
    """
    try:
        # Parse the ABI and bytecode
        abi = get_test_abi(abi_json)
        bytecode = get_test_bytecode(bytecode_hex)

        # Deploy the contract
        contract = await client.deploy_contract(
            account.signer,
            bytecode,
            abi,
            *constructor_args
        )
        return contract
    except Exception as e:
        print(f"Contract deployment failed: {e}")
        return None
