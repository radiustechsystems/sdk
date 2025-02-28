/**
 * Type definitions for Radius smart contract operations.
 *
 * This module maps Radius SDK's own data structures to Ethereum-compatible types.
 * While the SDK has its own concrete implementations of core concepts like addresses,
 * transactions, and network communication, we use these type mappings to leverage
 * the functionality provided by Ethereum libraries.
 *
 * The SDK maintains its own independent implementations optimized for Radius's architecture,
 * but maps to these Ethereum types when needed to take advantage of well-tested libraries
 * and functionality.
 */

import {
  // Types for numeric values, used for amounts, gas prices, etc.
  type BigNumberish,
  // Types for byte representations, used for binary data
  type BytesLike,
  // Fee data type (gas prices, etc.)
  type FeeData,
  // Network request types
  type FetchRequest,
  // Event log filtering
  type Filter,
  type FilterByBlockHash,
  // ABI and contract interface types
  type Interface,
  // JSON-RPC protocol types
  type JsonRpcPayload,
  type JsonRpcProvider,
  type JsonRpcResult,
  // Network information types
  type Network,
  // Provider types for Ethereum communication
  type Provider,
  // Cryptographic types
  type Signature,
  // Transaction types
  type TransactionRequest,
  type TransactionResponse,
  type Wallet,
} from 'ethers';

// Re-export all types for use in the Radius SDK
export type {
  /**
   * Represents numeric values that can be converted to big integers.
   * Used for amounts, gas prices, and other numeric values in Radius.
   * Can be a string, number, or bigint.
   */
  BigNumberish,
  /**
   * Represents data that can be converted to byte arrays.
   * Used for transaction data, method parameters, and other binary data.
   * Can be a string, byte array, or other binary-convertible format.
   */
  BytesLike,
  /**
   * HTTP request configuration for communicating with Radius.
   * Used internally by providers to make JSON-RPC calls.
   */
  FetchRequest,
  /**
   * Smart contract ABI (Application Binary Interface) definition.
   * Used to encode and decode interactions with Radius smart contracts.
   */
  Interface,
  /**
   * JSON-RPC request payload structure.
   * Used for communication with Radius JSON-RPC endpoints.
   */
  JsonRpcPayload,
  /**
   * Provider implementation for JSON-RPC based communication.
   * Used to connect to and interact with Radius JSON-RPC endpoints.
   */
  JsonRpcProvider,
  /**
   * JSON-RPC response result structure.
   * Used to parse Radius node responses.
   */
  JsonRpcResult,
  /**
   * Network information for Radius.
   * Contains details about the connected Radius environment.
   */
  Network,
  /**
   * Base interface for Radius communication.
   * Provides methods to interact with the Radius system.
   */
  Provider,
  /**
   * Structure for a transaction request to be sent to the network.
   * Contains transaction parameters like to, from, value, etc.
   */
  TransactionRequest,
  /**
   * Structure for a transaction response received from the network.
   * Contains transaction details and hash.
   */
  TransactionResponse,
  /**
   * Event log filter criteria for querying blockchain events.
   * Allows filtering by topics, address, and block range.
   */
  Filter,
  /**
   * Event log filter criteria for querying blockchain events in a specific block.
   * Filters by block hash instead of block range.
   */
  FilterByBlockHash,
  /**
   * Fee data containing gas price and priority fee information.
   * Used for estimating transaction costs.
   */
  FeeData,
  /**
   * Cryptographic signature for transactions and messages.
   * Used to verify authenticity and authorization in Radius.
   */
  Signature,
  /**
   * Account with private key for signing.
   * Used to create signed transactions and messages for Radius.
   */
  Wallet,
};
