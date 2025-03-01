"""Integration tests for the Radius SDK.

This module provides integration tests for the SDK, requiring a running Ethereum node.
"""

import asyncio
import unittest

from radius import (
    Hash,
    Transaction,
)
import os
import sys
from pathlib import Path

# Add the parent directory to sys.path
sys.path.insert(0, str(Path(__file__).parent))
from helpers import deploy_test_contract, setup_test_account, setup_test_client


class RadiusIntegrationTest(unittest.IsolatedAsyncioTestCase):
    """Base class for Radius integration tests.
    
    Provides setup and teardown for integration tests that require a live node.
    """

    async def asyncSetUp(self) -> None:
        """Set up the test environment.
        
        Creates a client and account for testing.
        
        Raises:
            unittest.SkipTest: If the required environment is not available

        """
        try:
            self.client = await setup_test_client()
            self.account = await setup_test_account(self.client)
        except Exception as e:
            self.skipTest(f"Test environment not available: {e}")

    async def asyncTearDown(self) -> None:
        """Clean up the test environment."""
        self.client = None
        self.account = None


class ClientTest(RadiusIntegrationTest):
    """Tests for the Client class."""

    async def test_chain_id(self) -> None:
        """Test getting the chain ID."""
        chain_id = await self.client.chain_id()
        self.assertIsInstance(chain_id, int)
        self.assertGreater(chain_id, 0)

    async def test_gas_price(self) -> None:
        """Test getting the gas price."""
        gas_price = await self.client.gas_price()
        self.assertIsInstance(gas_price, int)
        self.assertGreaterEqual(gas_price, 0)  # Radius may have zero gas price

    async def test_block_number(self) -> None:
        """Test getting the block number."""
        block_number = await self.client.block_number()
        self.assertIsInstance(block_number, int)
        self.assertGreaterEqual(block_number, 0)

    async def test_get_balance(self) -> None:
        """Test getting an account balance."""
        balance = await self.client.get_balance(self.account.address)
        self.assertIsInstance(balance, int)
        self.assertGreaterEqual(balance, 0)

    async def test_get_transaction_count(self) -> None:
        """Test getting a transaction count."""
        nonce = await self.client.get_transaction_count(self.account.address)
        self.assertIsInstance(nonce, int)
        self.assertGreaterEqual(nonce, 0)


class AccountTest(RadiusIntegrationTest):
    """Tests for the Account class."""

    async def test_get_balance(self) -> None:
        """Test getting the account balance."""
        balance = await self.account.get_balance(self.client)
        self.assertIsInstance(balance, int)
        self.assertGreaterEqual(balance, 0)

    async def test_sign_message(self) -> None:
        """Test signing a message."""
        message = b"Hello, Radius!"
        signature = await self.account.sign_message(message)

        self.assertIn("r", signature)
        self.assertIn("s", signature)
        self.assertIn("v", signature)

        self.assertEqual(len(signature["r"]), 32)
        self.assertEqual(len(signature["s"]), 32)
        self.assertEqual(len(signature["v"]), 1)

    async def test_send_transaction(self) -> None:
        """Test sending a simple transaction."""
        # Create a transaction to send a small amount of ETH to ourselves
        tx = Transaction(
            to=self.account.address,
            value=1000,
            gas_price=await self.client.gas_price(),
        )

        # Send the transaction
        tx_hash = await self.account.send_transaction(self.client, tx)

        # Verify the transaction hash
        self.assertIsInstance(tx_hash, Hash)

        # Wait for the receipt
        receipt = None
        for _ in range(5):  # Wait up to 5 seconds
            receipt = await self.client.get_transaction_receipt(tx_hash)
            if receipt:
                break
            await asyncio.sleep(1)

        # Verify the receipt
        self.assertIsNotNone(receipt)
        self.assertTrue(receipt.status)


class ContractTest(RadiusIntegrationTest):
    """Tests for the Contract class."""

    async def asyncSetUp(self) -> None:
        """Set up the test environment.
        
        Creates a client, account, and deploys a test contract.
        
        Raises:
            unittest.SkipTest: If the required environment is not available

        """
        await super().asyncSetUp()

        # Path to the SimpleStorage contract files
        contracts_dir = os.path.join(os.path.dirname(os.path.dirname(os.path.dirname(__file__))),
                                      "contracts", "solidity")

        # Load the ABI and bytecode
        try:
            with open(os.path.join(contracts_dir, "SimpleStorage.abi")) as f:
                abi_json = f.read()

            with open(os.path.join(contracts_dir, "SimpleStorage.bin")) as f:
                bytecode_hex = f.read().strip()

            # Deploy the test contract (without constructor arguments)
            self.contract = await deploy_test_contract(
                self.client, self.account, abi_json, bytecode_hex
            )
            
            # Set the initial value to 42
            tx_hash = await self.contract.execute(
                self.client, self.account, "set", 42,
                gas_price=await self.client.gas_price(),
            )
            
            # Wait for the transaction to be mined
            receipt = None
            for _ in range(5):  # Wait up to 5 seconds
                receipt = await self.client.get_transaction_receipt(tx_hash)
                if receipt:
                    break
                await asyncio.sleep(1)
                
            if not receipt or not receipt.status:
                self.skipTest("Failed to set initial value for contract")
        except Exception as e:
            self.skipTest(f"Could not deploy test contract: {e}")

    async def test_call(self) -> None:
        """Test calling a contract method."""
        # Call the get method
        result = await self.contract.call(self.client, "get")

        # Verify the result
        self.assertEqual(len(result), 1)
        self.assertEqual(result[0], 42)

    async def test_execute(self) -> None:
        """Test executing a contract method."""
        # Execute the set method
        new_value = 99
        tx_hash = await self.contract.execute(
            self.client, self.account, "set", new_value,
            gas_price=await self.client.gas_price(),
        )

        # Wait for the transaction to be mined
        receipt = None
        for _ in range(5):  # Wait up to 5 seconds
            receipt = await self.client.get_transaction_receipt(tx_hash)
            if receipt:
                break
            await asyncio.sleep(1)

        # Verify the transaction succeeded
        self.assertIsNotNone(receipt)
        self.assertTrue(receipt.status)

        # Check that the value was updated
        result = await self.contract.call(self.client, "get")
        self.assertEqual(result[0], new_value)


if __name__ == "__main__":
    unittest.main()
