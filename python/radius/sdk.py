"""SDK core module that re-exports all public components.

This module serves as the main entry point for the Radius SDK, providing access to all public
classes, functions, and types.
"""

from __future__ import annotations

# Import from src package through sys.path
from src.accounts.account import Account
from src.accounts.options import AccountOption, with_private_key, with_signer
from src.auth.clef.signer import ClefSigner
from src.auth.privatekey.signer import PrivateKeySigner
from src.auth.types import Signer, SignerClient
from src.client.client import Client
from src.client.options import ClientOption, with_http_client, with_interceptor, with_logger
from src.common.abi import ABI, abi_from_json
from src.common.address import Address, address_from_hex, zero_address
from src.common.constants import MAX_GAS
from src.common.event import Event
from src.common.hash import Hash, hash_from_hex
from src.common.http import HttpClient
from src.common.receipt import Receipt
from src.common.transaction import SignedTransaction, Transaction
from src.common.utils import bytecode_from_hex
from src.contracts.contract import Contract
from src.providers.eth.types import BigNumberish, BytesLike
from src.transport.types import Interceptor, Logf

# Classes re-exports
__all__ = [
    # Classes
    "ABI",
    "Account",
    "Address",
    "Client",
    "ClefSigner",
    "Contract",
    "Event",
    "Hash",
    "PrivateKeySigner",
    "Receipt",
    "Transaction",
    # Constants
    "MAX_GAS",
    # Functions
    "abi_from_json",
    "address_from_hex",
    "bytecode_from_hex",
    "hash_from_hex",
    "new_account",
    "new_clef_signer",
    "new_client",
    "new_contract",
    "with_http_client",
    "with_interceptor",
    "with_logger",
    "with_private_key",
    "with_signer",
    "zero_address",
    # Types
    "AccountOption",
    "BigNumberish",
    "BytesLike",
    "ClientOption",
    "HttpClient",
    "Interceptor",
    "Logf",
    "Signer",
    "SignerClient",
    "SignedTransaction",
]


def new_clef_signer(address: Address, client: SignerClient, clef_url: str) -> ClefSigner:
    """Create a new ClefSigner with the given address, client, and Clef URL.
    
    ClefSigner provides a way to sign transactions using Clef as an external signing service.
    
    Args:
        address: The address to use for signing
        client: The Radius client to use for transaction-related operations
        clef_url: The URL of the Clef server
        
    Returns:
        A new ClefSigner instance
        
    Raises:
        ValueError: If unable to connect to the Clef server

    """
    return ClefSigner(address, client, clef_url)


async def new_client(url: str, *opts: ClientOption) -> Client:
    """Create a new Client with the given URL and options.
    
    The client is the main entry point for interacting with the Radius platform.
    
    Args:
        url: The URL of the Radius JSON-RPC endpoint
        opts: Additional options for the client configuration
        
    Returns:
        A new Client instance
        
    Raises:
        ValueError: If the client cannot be created or cannot connect to Radius

    """
    return await Client.new(url, *opts)


def new_contract(address: Address, abi: ABI) -> Contract:
    """Create a new Contract with the given address and ABI.
    
    The Contract object allows interaction with smart contracts deployed on Radius.
    
    Args:
        address: The contract address on Radius
        abi: The contract ABI defining its methods and events
        
    Returns:
        A new Contract instance

    """
    return Contract(address, abi)


async def new_account(*opts: AccountOption) -> Account:
    """Create a new Account with the given options.
    
    Accounts are used to represent Radius accounts and manage their keys.
    
    Args:
        opts: Account options for configuring the account (private key, signer, etc.)
        
    Returns:
        A new Account instance
        
    Raises:
        ValueError: If the account cannot be created with the provided options

    """
    return await Account.new(*opts)
