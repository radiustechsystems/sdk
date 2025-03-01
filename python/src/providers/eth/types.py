"""Ethereum type definitions for the Radius SDK.

This module provides type definitions and utilities for working with Ethereum types.
"""

from __future__ import annotations

from typing import Any, Dict, Union

# Type for values that can be treated as big numbers
BigNumberish = Union[int, str, bytes]

# Type for values that can be treated as byte arrays
BytesLike = Union[bytes, str, bytearray, memoryview]

# Ethereum JSON-RPC request
EthRequest = Dict[str, Any]

# Ethereum JSON-RPC response
EthResponse = Dict[str, Any]
