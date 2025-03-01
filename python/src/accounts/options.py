"""Account options for the Radius SDK.

This module provides functional options for configuring accounts.
"""

from __future__ import annotations

from typing import Callable, TypeVar, Union

from src.accounts.types import AccountClient
from src.auth.privatekey.signer import PrivateKeySigner
from src.auth.types import Signer

T = TypeVar("T")

# Type for account options
AccountOption = Callable[[T], T]


def with_private_key(private_key: Union[str, bytes], client: AccountClient) -> AccountOption:
    """Create an account option that sets the private key signer.
    
    Args:
        private_key: The private key to use for signing
        client: The client to use for transaction operations
        
    Returns:
        An account option function
        
    Raises:
        ValueError: If the private key is invalid

    """
    def apply_option(account: T) -> T:
        signer = PrivateKeySigner(private_key, client)
        account._signer = signer  # type: ignore
        return account

    return apply_option


def with_signer(signer: Signer) -> AccountOption:
    """Create an account option that sets the signer.
    
    Args:
        signer: The signer to use
        
    Returns:
        An account option function

    """
    def apply_option(account: T) -> T:
        account._signer = signer  # type: ignore
        return account

    return apply_option
