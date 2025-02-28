import { beforeAll, describe, expect, test } from 'vitest';
import { Account, Client, Contract, ZeroAddress } from '../radius';
import {
  SIMPLE_STORAGE_ABI,
  SIMPLE_STORAGE_BIN,
  createTestAccount,
  createTestClient,
  getFundedAccount,
  getTestABI,
  getTestBytecode,
} from './helpers';

describe('Integration Tests', async () => {
  let client: Client;
  let fundedAccount: Account;
  let setupFailed = false;

  beforeAll(async () => {
    const endpoint: string = process.env.RADIUS_ENDPOINT;
    const privateKey: string = process.env.RADIUS_PRIVATE_KEY;
    if (!endpoint || !privateKey) {
      setupFailed = true;
      console.log('Skipping tests due to missing environment variables');
      return;
    }

    client = await createTestClient(endpoint);
    if (!client) {
      setupFailed = true;
      console.log('Skipping tests due to failed client creation');
      return;
    }

    fundedAccount = await getFundedAccount(client, privateKey);
    if (!fundedAccount) {
      setupFailed = true;
      console.log('Skipping tests due to failed account creation');
      return;
    }
  });

  test('Send transaction between accounts', async () => {
    if (setupFailed) {
      return;
    }

    // Create recipient account
    const recipient = await createTestAccount(client);
    const amount = BigInt(100);

    // Get initial balance
    const initialBalance = await fundedAccount.balance(client);
    console.log(`Initial balance: ${initialBalance}`);

    // Send ETH from test account to recipient
    const receipt = await fundedAccount.send(client, recipient.address(), amount);
    expect(receipt, 'Transaction receipt should be defined').toBeDefined();
    expect(receipt.from, 'Transaction from address should match sender').toEqual(
      fundedAccount.address()
    );
    expect(receipt.to, 'Transaction to address should match recipient').toEqual(
      recipient.address()
    );
    expect(receipt.value, 'Transaction value should match sent amount').toBe(amount);

    // Check sender balance
    const senderBalance = await fundedAccount.balance(client);
    console.log(`Sender balance after transfer: ${senderBalance}`);
    expect(senderBalance, 'Sender balance should be the initial balance minus amount').toBe(
      initialBalance - amount
    );

    // Check recipient balance
    const recipientBalance = await recipient.balance(client);
    console.log(`Recipient balance: ${recipientBalance}`);
    expect(recipientBalance, 'Recipient balance should equal sent amount').toBe(amount);
  });

  test('Deploy and interact with SimpleStorage contract', async () => {
    if (setupFailed) {
      return;
    }

    // Deploy contract
    const bytecode = getTestBytecode(SIMPLE_STORAGE_BIN);
    const abi = getTestABI(SIMPLE_STORAGE_ABI);
    const contract: Contract = await client.deployContract(fundedAccount.signer, bytecode, abi);
    expect(contract.address(), 'Contract address should be defined after deployment').toBeDefined();
    expect(contract.address(), 'Contract address should not be zero address').not.toEqual(
      ZeroAddress()
    );

    console.log(`Contract deployed at: ${contract.address().hex()}`);

    // Call Set method
    const expectedValue = BigInt(42);
    const setReceipt = await contract.execute(client, fundedAccount.signer, 'set', expectedValue);
    expect(setReceipt, 'Set method receipt should be defined').toBeDefined();
    expect(setReceipt.from, 'Set method from address should match sender').toEqual(
      fundedAccount.address()
    );
    expect(setReceipt.to, 'Set method to address should match contract').toEqual(
      contract.address()
    );

    // Call Get method
    const result = await contract.call(client, 'get');
    expect(result, 'Get method should return an array').toHaveLength(1);
    expect(result[0], 'Get method should return the set value').toBe(expectedValue);

    console.log(`Retrieved value from contract: ${result[0]}`);
  });
});
