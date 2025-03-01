"""Transport layer types for the Radius SDK.

This module provides types and protocols for the transport layer.
"""

from __future__ import annotations

from typing import Any, Callable, Dict, Protocol


class Interceptor(Protocol):
    """Protocol for request/response interceptors.
    
    Interceptors can modify requests before they are sent and responses before they are returned.
    """

    def intercept_request(self, url: str, request: Dict[str, Any]) -> Dict[str, Any]:
        """Intercept and potentially modify a request before it is sent.
        
        Args:
            url: The URL the request is being sent to
            request: The request data
            
        Returns:
            The modified request data

        """
        ...

    def intercept_response(
        self, url: str, request: Dict[str, Any], response: Dict[str, Any]
    ) -> Dict[str, Any]:
        """Intercept and potentially modify a response before it is returned.
        
        Args:
            url: The URL the request was sent to
            request: The request data
            response: The response data
            
        Returns:
            The modified response data

        """
        ...


# Type for logging functions
Logf = Callable[[str, Any], None]
