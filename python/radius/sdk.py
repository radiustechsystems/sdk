"""SDK core module that re-exports all public components.

This module serves as the main entry point for the Radius SDK, providing access to all public
classes, functions, and types.
"""

from __future__ import annotations

# Import from src package through sys.path
from src.accounts.account import Account
from src.accounts.options import AccountOption, with_private_key, with_signer
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
    "new_abi",
    "new_account",
    "new_address",
    "new_clef_signer",
    "new_client",
    "new_contract",
    "new_private_key_signer",
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

# Helper functions to create objects
def new_abi(json_abi: str) -> ABI:
    """Create a new ABI from JSON string."""
    return abi_from_json(json_abi)
    
def new_address(hex_address: str) -> Address:
    """Create a new Address from hex string."""
    return address_from_hex(hex_address)

def new_private_key_signer(private_key: str | bytes, chain_id: BigNumberish) -> PrivateKeySigner:
    """Create a new PrivateKeySigner with the given private key and chain ID."""
    return PrivateKeySigner(private_key, chain_id)

async def new_client(url: str, *opts: ClientOption) -> Client:
    """Create a new Client with the given URL and options."""
    return await Client.new(url, *opts)

def new_contract(address: Address, abi: ABI) -> Contract:
    """Create a new Contract with the given address and ABI."""
    return Contract.new(address, abi)

async def new_account(*opts: AccountOption) -> Account:
    """Create a new Account with the given options."""
    return await Account.new(*opts)
