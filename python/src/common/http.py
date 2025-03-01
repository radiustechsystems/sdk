"""HTTP client for the Radius SDK.

This module provides the HTTP client implementation for making requests to the blockchain.
"""

from __future__ import annotations

import json
from typing import Any, Dict, Optional, Protocol, TypeVar

import requests
from requests.adapters import HTTPAdapter
from urllib3.util.retry import Retry

from src.transport.types import Interceptor

T = TypeVar("T")


class HttpClient(Protocol):
    """Protocol for HTTP clients used in the Radius SDK."""

    def post(self, url: str, data: Dict[str, Any]) -> Dict[str, Any]:
        """Send a POST request.
        
        Args:
            url: The URL to send the request to
            data: The data to send in the request body
            
        Returns:
            The response body as a dictionary
            
        Raises:
            requests.RequestException: If the request fails

        """
        ...


class DefaultHttpClient:
    """Default implementation of the HttpClient interface.
    
    Provides HTTP functionality with request/response interception, retry, and logging.
    """

    def __init__(
        self,
        interceptor: Optional[Interceptor] = None,
        max_retries: int = 3,
        timeout: int = 30,
    ) -> None:
        """Initialize a new HTTP client.
        
        Args:
            interceptor: Optional request/response interceptor
            max_retries: Maximum number of retries for failed requests
            timeout: Request timeout in seconds

        """
        self._session = requests.Session()
        self._interceptor = interceptor
        self._timeout = timeout

        # Configure retry strategy
        retry_strategy = Retry(
            total=max_retries,
            backoff_factor=0.5,
            status_forcelist=[500, 502, 503, 504],
            allowed_methods=["POST"],
        )
        adapter = HTTPAdapter(max_retries=retry_strategy)
        self._session.mount("http://", adapter)
        self._session.mount("https://", adapter)

    def post(self, url: str, data: Dict[str, Any]) -> Dict[str, Any]:
        """Send a POST request.
        
        Args:
            url: The URL to send the request to
            data: The data to send in the request body
            
        Returns:
            The response body as a dictionary
            
        Raises:
            requests.RequestException: If the request fails
            ValueError: If the response is not valid JSON

        """
        # Apply request interceptor if available
        request_data = data
        if self._interceptor:
            request_data = self._interceptor.intercept_request(url, data)

        # Send the request
        response = self._session.post(
            url,
            json=request_data,
            timeout=self._timeout,
            headers={"Content-Type": "application/json"},
        )

        # Raise an exception for HTTP errors
        response.raise_for_status()

        # Parse the response
        try:
            response_data = response.json()
        except json.JSONDecodeError as e:
            raise ValueError(f"Invalid JSON response: {response.text}") from e

        # Apply response interceptor if available
        if self._interceptor:
            response_data = self._interceptor.intercept_response(url, request_data, response_data)

        return response_data
