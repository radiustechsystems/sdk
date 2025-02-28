/**
 * Radius SDK for TypeScript/JavaScript
 *
 * This is the main entry point for the Radius SDK. It exports all the classes,
 * types, and functions needed to interact with the Radius platform.
 *
 * Example usage:
 * ```typescript
 * import { NewClient, NewAccount, Address, withPrivateKey } from '@radiustechsystems/sdk';
 *
 * // Create a new client connected to a Radius node
 * const client = await NewClient('https://testnet.radius.xyz');
 *
 * // Create a new account with a private key
 * const account = await NewAccount(withPrivateKey('0x1234...'));
 *
 * // Check the account balance
 * const balance = await account.balance(client);
 * ```
 *
 * @module radius-sdk
 */

// Export everything from the SDK module for a unified interface
export * from './sdk';
