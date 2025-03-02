"""Account options for the Radius SDK.

This module provides functional options for configuring accounts.
"""

from __future__ import annotations

from typing import Callable, TypeVar, Union

from src.auth.types import Signer
from src.providers.eth.types import BigNumberish

T = TypeVar("T")

# Type for account options
AccountOption = Callable[[T], T]


def with_private_key(private_key: Union[str, bytes], chain_id: BigNumberish) -> AccountOption:
    """Create an account option that sets the private key signer.
    
    Args:
        private_key: The private key to use for signing
        chain_id: The chain ID to use for signing
        
    Returns:
        An account option function
    """
    from src.auth.privatekey.signer import PrivateKeySigner
    
    def apply_option(account: T) -> T:
        signer = PrivateKeySigner(private_key, chain_id)
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
