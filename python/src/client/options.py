"""Client options for the Radius SDK.

This module provides functional options for configuring the Radius client.
"""

from __future__ import annotations

from typing import Callable, TypeVar

from src.common.http import DefaultHttpClient, HttpClient
from src.transport.interceptor import LoggingInterceptor
from src.transport.types import Interceptor, Logf

T = TypeVar("T")

# Type for client options
ClientOption = Callable[[T], T]


def with_http_client(http_client: HttpClient) -> ClientOption:
    """Create a client option that sets the HTTP client.
    
    Args:
        http_client: The HTTP client to use for requests
        
    Returns:
        A client option function

    """
    def apply_option(client: T) -> T:
        client._http_client = http_client  # type: ignore
        return client

    return apply_option


def with_interceptor(interceptor: Interceptor) -> ClientOption:
    """Create a client option that sets the request/response interceptor.
    
    Args:
        interceptor: The interceptor to use
        
    Returns:
        A client option function

    """
    def apply_option(client: T) -> T:
        http_client = DefaultHttpClient(interceptor=interceptor)
        client._http_client = http_client  # type: ignore
        return client

    return apply_option


def with_logger(logger: Logf) -> ClientOption:
    """Create a client option that sets up logging for requests and responses.
    
    Args:
        logger: The logging function to use
        
    Returns:
        A client option function

    """
    def apply_option(client: T) -> T:
        interceptor = LoggingInterceptor(logger)
        http_client = DefaultHttpClient(interceptor=interceptor)
        client._http_client = http_client  # type: ignore
        return client

    return apply_option
