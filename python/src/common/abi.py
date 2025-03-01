"""ABI module for the Radius SDK.

This module provides utilities for working with Ethereum Contract ABIs.
"""

from __future__ import annotations

import json
from typing import Any, Dict, List, Optional

from web3 import Web3

# Create a web3 instance with a null provider (we only need encoding functionality)
w3 = Web3()


class ABI:
    """Represents an Ethereum Contract ABI (Application Binary Interface).
    
    The ABI defines the methods and events of a smart contract and how to call them.
    """

    def __init__(self, raw_abi: List[Dict[str, Any]]) -> None:
        """Initialize an ABI with the underlying ABI structure.
        
        Args:
            raw_abi: The raw ABI structure

        """
        self._raw_abi = raw_abi
        self._methods: Dict[str, Dict[str, Any]] = {}
        self._events: Dict[str, Dict[str, Any]] = {}

        # Process the ABI to categorize methods and events
        for item in raw_abi:
            item_type = item.get("type")
            item_name = item.get("name")

            if not item_name:
                continue

            if item_type == "function":
                self._methods[item_name] = item
            elif item_type == "event":
                self._events[item_name] = item

    @property
    def methods(self) -> Dict[str, Dict[str, Any]]:
        """Get the methods defined in the ABI.
        
        Returns:
            A dictionary of method names to method definitions

        """
        return self._methods

    @property
    def events(self) -> Dict[str, Dict[str, Any]]:
        """Get the events defined in the ABI.
        
        Returns:
            A dictionary of event names to event definitions

        """
        return self._events

    def get_method(self, name: str) -> Optional[Dict[str, Any]]:
        """Get a method by name.
        
        Args:
            name: The name of the method
            
        Returns:
            The method definition, or None if not found

        """
        return self._methods.get(name)

    def get_event(self, name: str) -> Optional[Dict[str, Any]]:
        """Get an event by name.
        
        Args:
            name: The name of the event
            
        Returns:
            The event definition, or None if not found

        """
        return self._events.get(name)

    def encode_function_data(self, function_name: str, *args: Any) -> bytes:
        """Encode function data for a contract call.
        
        Args:
            function_name: The name of the function to call
            args: The arguments to pass to the function
            
        Returns:
            The encoded function data
            
        Raises:
            ValueError: If the function is not found in the ABI

        """
        method = self.get_method(function_name)
        if not method:
            raise ValueError(f"Function {function_name} not found in ABI")

        contract = w3.eth.contract(abi=self._raw_abi)
        method = getattr(contract.functions, function_name)
        
        # Use _encode_transaction_data to get the raw function call data
        # This gives us just the function signature + args without needing a transaction
        function_call = method(*args)
        encoded_data = function_call._encode_transaction_data()
        
        # Remove '0x' prefix if present and convert to bytes
        if encoded_data.startswith('0x'):
            encoded_data = encoded_data[2:]
        return bytes.fromhex(encoded_data)

    def decode_function_result(self, function_name: str, data: bytes) -> List[Any]:
        """Decode the result of a function call.
        
        Args:
            function_name: The name of the function
            data: The encoded result data
            
        Returns:
            The decoded result
            
        Raises:
            ValueError: If the function is not found in the ABI

        """
        method = self.get_method(function_name)
        if not method:
            raise ValueError(f"Function {function_name} not found in ABI")

        # For empty result, return empty list
        if not data or data == b'':
            return []
        
        # Create a contract instance
        contract = w3.eth.contract(abi=self._raw_abi)
        
        # Get the output types for decoding
        output_types = [output['type'] for output in method.get('outputs', [])]
        
        if len(output_types) == 0:
            return []
            
        # Use eth-abi to decode the result
        from eth_abi import decode
        
        try:
            # If we have a single output type, decode it directly
            if len(output_types) == 1:
                result = decode([output_types[0]], data)
                return [result[0]]
            else:
                # Multiple output types
                result = decode(output_types, data)
                return list(result)
        except Exception as e:
            # If decoding fails, try treating the data as a function return value
            # which might be padded to 32 bytes
            try:
                if len(output_types) == 1 and output_types[0] == 'uint256':
                    # For uint256, convert the 32 bytes to an integer
                    value = int.from_bytes(data, byteorder='big')
                    return [value]
                else:
                    # For other types, try decoding with padding
                    return [decode([t], data[-32:])[0] for t in output_types]
            except Exception as inner_e:
                raise ValueError(f"Failed to decode function result: {e}, {inner_e}") from e

    def __str__(self) -> str:
        """Get the string representation of the ABI.
        
        Returns:
            The ABI as a JSON string

        """
        return json.dumps(self._raw_abi)


def abi_from_json(json_str: str) -> Optional[ABI]:
    """Create an ABI from a JSON string.
    
    Args:
        json_str: The JSON string representing the ABI
        
    Returns:
        An ABI instance, or None if the JSON is invalid

    """
    try:
        abi_data = json.loads(json_str)
        if isinstance(abi_data, list):
            return ABI(abi_data)
        return None
    except json.JSONDecodeError:
        return None
