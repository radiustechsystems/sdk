"""Interceptor implementations for the Radius SDK.

This module provides concrete implementations of the Interceptor protocol.
"""

from __future__ import annotations

import json
from typing import Any, Dict

from src.transport.types import Logf


class LoggingInterceptor:
    """An interceptor that logs requests and responses.
    
    This interceptor logs the details of requests and responses for debugging purposes.
    """

    def __init__(self, logger: Logf) -> None:
        """Initialize a new logging interceptor.
        
        Args:
            logger: A function that will be called to log messages

        """
        self._logger = logger

    def intercept_request(self, url: str, request: Dict[str, Any]) -> Dict[str, Any]:
        """Log the request and return it unchanged.
        
        Args:
            url: The URL the request is being sent to
            request: The request data
            
        Returns:
            The request data unchanged

        """
        self._logger("Request to %s: %s", url, json.dumps(request, indent=2))
        return request

    def intercept_response(
        self, url: str, request: Dict[str, Any], response: Dict[str, Any]
    ) -> Dict[str, Any]:
        """Log the response and return it unchanged.
        
        Args:
            url: The URL the request was sent to
            request: The request data
            response: The response data
            
        Returns:
            The response data unchanged

        """
        self._logger("Response from %s: %s", url, json.dumps(response, indent=2))
        return response
