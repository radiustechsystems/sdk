"""Integration tests for the Radius SDK.

This module provides integration tests for the SDK, requiring a running Ethereum node.
"""

import unittest
import os
import sys
from pathlib import Path

# Add the parent directory to sys.path
sys.path.insert(0, str(Path(__file__).parent))
from helpers import (
    SIMPLE_STORAGE_ABI,
    SIMPLE_STORAGE_BIN,
    create_test_account,
    create_test_client,
    deploy_test_contract,
    get_funded_account,
    get_test_abi,
    get_test_bytecode,
    skip_if_insufficient_balance,
)


class IntegrationTest(unittest.IsolatedAsyncioTestCase):
    """Integration tests for the Radius SDK."""
    
    async def asyncSetUp(self) -> None:
        """Set up the test environment."""
        # Setup will be handled in the individual tests
        self.client = None
        self.funded_account = None
        self.setup_failed = False

    async def test_send_transaction(self) -> None:
        """Test sending a transaction between accounts."""
        # Create client
        endpoint = os.environ.get("RADIUS_ENDPOINT")
        private_key = os.environ.get("RADIUS_PRIVATE_KEY")
        if not endpoint or not private_key:
            self.skipTest("Missing environment variables")
            return

        self.client = await create_test_client(endpoint)
        if not self.client:
            self.skipTest("Failed to create client")
            return

        # Get funded account
        self.funded_account = await get_funded_account(self.client, private_key)
        if not self.funded_account:
            self.skipTest("Failed to get funded account")
            return

        # Verify account has sufficient balance
        initial_balance = await skip_if_insufficient_balance(self.funded_account, self.client)
        if not initial_balance:
            self.skipTest("Insufficient account balance")
            return

        # Create recipient account
        recipient = await create_test_account(self.client)
        self.assertIsNotNone(recipient, "Failed to create recipient account")
        self.assertIsNotNone(recipient.address(), "Recipient address should not be None")

        # Send a small amount
        amount = 100
        receipt = await self.funded_account.send(
            self.client,
            recipient.address(),
            amount
        )
        self.assertIsNotNone(receipt, "Transaction receipt should not be None")
        self.assertEqual(self.funded_account.address(), receipt.from_address, "Unexpected sender address")
        self.assertEqual(recipient.address(), receipt.to_address, "Unexpected recipient address")

        # Check sender balance
        sender_balance = await self.funded_account.balance(self.client)
        self.assertLessEqual(sender_balance, initial_balance - amount, "Sender balance should decrease by at least the amount sent")

        # Check recipient balance
        recipient_balance = await recipient.balance(self.client)
        self.assertEqual(recipient_balance, amount, "Recipient balance should equal the amount sent")
        
    async def test_simple_storage_contract(self) -> None:
        """Test deploying and interacting with the SimpleStorage contract."""
        # Create client
        endpoint = os.environ.get("RADIUS_ENDPOINT")
        private_key = os.environ.get("RADIUS_PRIVATE_KEY")
        if not endpoint or not private_key:
            self.skipTest("Missing environment variables")
            return

        self.client = await create_test_client(endpoint)
        if not self.client:
            self.skipTest("Failed to create client")
            return

        # Get funded account
        self.funded_account = await get_funded_account(self.client, private_key)
        if not self.funded_account:
            self.skipTest("Failed to get funded account")
            return

        # Deploy contract
        bytecode = get_test_bytecode(SIMPLE_STORAGE_BIN)
        abi = get_test_abi(SIMPLE_STORAGE_ABI)
        
        contract = await self.client.deploy_contract(
            self.funded_account.signer,
            bytecode,
            abi
        )
        self.assertIsNotNone(contract, "Contract should not be None")
        self.assertIsNotNone(contract.address(), "Contract address should not be None")
        
        # Set storage value
        value = 42
        set_receipt = await contract.execute(
            self.client,
            self.funded_account.signer,
            "set",
            value
        )
        self.assertIsNotNone(set_receipt, "Set method receipt should not be None")
        self.assertEqual(self.funded_account.address(), set_receipt.from_address, "Unexpected from address")
        self.assertEqual(contract.address(), set_receipt.to_address, "Unexpected to address")
        
        # Get storage value
        result = await contract.call(self.client, "get")
        self.assertEqual(1, len(result), "Result should contain 1 value")
        self.assertEqual(value, result[0], "Value should match what was set")
        
        # Update storage value
        new_value = 99
        update_receipt = await self.client.execute(
            contract,
            self.funded_account.signer,
            "set",
            new_value
        )
        self.assertIsNotNone(update_receipt, "Update receipt should not be None")
        self.assertTrue(update_receipt.status, "Transaction should succeed")
        
        # Verify updated value
        updated_result = await contract.call(self.client, "get")
        self.assertEqual(1, len(updated_result), "Result should contain 1 value")
        self.assertEqual(new_value, updated_result[0], "Value should match the updated value")


if __name__ == "__main__":
    unittest.main()
