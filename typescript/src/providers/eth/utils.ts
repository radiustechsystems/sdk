/**
 * Utility functions and classes for Radius smart contract operations.
 *
 * This module provides access to Ethereum library utilities for data conversion,
 * cryptographic functions, and network communication. While the Radius SDK has its
 * own implementations of core functionality, this module allows the SDK to leverage
 * well-tested Ethereum utilities when beneficial.
 *
 * The SDK maps its own data structures to Ethereum-compatible formats only when
 * needed to utilize these utility functions, maintaining independence while
 * benefiting from established libraries.
 */
import {
  concat as ethConcat,
  ethers,
  getAddress as ethGetAddress,
  getBytes as ethGetBytes,
  hexlify as ethHexlify,
  keccak256 as ethKeccak256,
  toNumber as ethToNumber,
  toUtf8Bytes as ethToUtf8Bytes,
} from 'ethers';

/**
 * Collection of Ethereum utilities mapped for use in the Radius SDK.
 * These functions provide wrapped access to Ethereum library functionality
 * with built-in error handling to maintain loose coupling between the SDK and
 * its dependencies. Each function catches and transforms errors from the underlying
 * Ethereum library, providing more consistent error handling throughout the SDK.
 */
export const eth = {
  /**
   * Concatenates multiple byte arrays into a single array
   * @throws Error if the concatenation fails
   */
  concat: (...args: Parameters<typeof ethConcat>) => {
    try {
      return ethConcat(...args);
    } catch (error) {
      throw new Error(
        `Failed to concatenate byte arrays: ${
          error instanceof Error ? error.message : String(error)
        }`
      );
    }
  },

  /**
   * Converts a value to a checksummed Ethereum address
   * @throws Error if the address is invalid
   */
  getAddress: (address: string) => {
    try {
      return ethGetAddress(address);
    } catch (error) {
      throw new Error(
        `Invalid Ethereum address format: ${error instanceof Error ? error.message : String(error)}`
      );
    }
  },

  /**
   * Converts a value to a Uint8Array of bytes
   * @throws Error if the conversion fails
   */
  getBytes: (value: string | Uint8Array) => {
    try {
      return ethGetBytes(value);
    } catch (error) {
      throw new Error(
        `Failed to convert to bytes: ${error instanceof Error ? error.message : String(error)}`
      );
    }
  },

  /**
   * Converts a value to a hexadecimal string
   * @throws Error if the conversion fails
   */
  hexlify: (value: string | Uint8Array) => {
    try {
      return ethHexlify(value);
    } catch (error) {
      throw new Error(
        `Failed to convert to hex: ${error instanceof Error ? error.message : String(error)}`
      );
    }
  },

  /**
   * Computes the Keccak-256 cryptographic hash of a value
   * @throws Error if the hashing fails
   */
  keccak256: (value: string | Uint8Array) => {
    try {
      return ethKeccak256(value);
    } catch (error) {
      throw new Error(
        `Failed to compute Keccak-256 hash: ${
          error instanceof Error ? error.message : String(error)
        }`
      );
    }
  },

  /**
   * Converts a numeric value to a JavaScript number
   * @throws Error if the conversion fails
   */
  toNumber: (value: bigint | number | string) => {
    try {
      return ethToNumber(value);
    } catch (error) {
      throw new Error(
        `Failed to convert to number: ${error instanceof Error ? error.message : String(error)}`
      );
    }
  },

  /**
   * Converts a string to its UTF-8 byte representation
   * @throws Error if the conversion fails
   */
  toUtf8Bytes: (text: string) => {
    try {
      return ethToUtf8Bytes(text);
    } catch (error) {
      throw new Error(
        `Failed to convert to UTF-8 bytes: ${
          error instanceof Error ? error.message : String(error)
        }`
      );
    }
  },

  /**
   * Request class for making HTTP requests to Radius JSON-RPC endpoints
   */
  FetchRequest: ethers.FetchRequest,

  /**
   * Class for working with contract ABIs (Application Binary Interfaces)
   * Provides wrapped constructor with error handling
   */
  Interface: class extends ethers.Interface {
    constructor(...args: ConstructorParameters<typeof ethers.Interface>) {
      try {
        super(...args);
      } catch (error) {
        throw new Error(
          `Failed to create Interface: ${error instanceof Error ? error.message : String(error)}`
        );
      }
    }
  },

  /**
   * Provider class that connects to Radius JSON-RPC endpoints
   * Methods are wrapped with error handling
   */
  JsonRpcProvider: class extends ethers.JsonRpcProvider {
    constructor(...args: ConstructorParameters<typeof ethers.JsonRpcProvider>) {
      try {
        super(...args);
      } catch (error) {
        throw new Error(
          `Failed to create JsonRpcProvider: ${
            error instanceof Error ? error.message : String(error)
          }`
        );
      }
    }

    // Override key methods with error handling
    override async getBalance(address: string, blockTag?: string | number): Promise<bigint> {
      try {
        return await super.getBalance(address, blockTag);
      } catch (error) {
        throw new Error(
          `Failed to get balance: ${error instanceof Error ? error.message : String(error)}`
        );
      }
    }

    override async getNetwork(): Promise<ethers.Network> {
      try {
        return await super.getNetwork();
      } catch (error) {
        throw new Error(
          `Failed to get network: ${error instanceof Error ? error.message : String(error)}`
        );
      }
    }

    override async getBlockNumber(): Promise<number> {
      try {
        return await super.getBlockNumber();
      } catch (error) {
        throw new Error(
          `Failed to get block number: ${error instanceof Error ? error.message : String(error)}`
        );
      }
    }

    // Using getFeeData instead of getGasPrice (which is deprecated)
    async getFeeData(): Promise<ethers.FeeData> {
      try {
        return await super.getFeeData();
      } catch (error) {
        throw new Error(
          `Failed to get fee data: ${error instanceof Error ? error.message : String(error)}`
        );
      }
    }

    override async estimateGas(tx: ethers.TransactionRequest): Promise<bigint> {
      try {
        return await super.estimateGas(tx);
      } catch (error) {
        throw new Error(
          `Failed to estimate gas: ${error instanceof Error ? error.message : String(error)}`
        );
      }
    }

    override async call(tx: ethers.TransactionRequest): Promise<string> {
      try {
        return await super.call(tx);
      } catch (error) {
        throw new Error(
          `Failed to call contract: ${error instanceof Error ? error.message : String(error)}`
        );
      }
    }

    override async getCode(address: string): Promise<string> {
      try {
        return await super.getCode(address);
      } catch (error) {
        throw new Error(
          `Failed to get code: ${error instanceof Error ? error.message : String(error)}`
        );
      }
    }

    override async getTransaction(hash: string): Promise<ethers.TransactionResponse | null> {
      try {
        return await super.getTransaction(hash);
      } catch (error) {
        throw new Error(
          `Failed to get transaction: ${error instanceof Error ? error.message : String(error)}`
        );
      }
    }

    override async getTransactionReceipt(hash: string): Promise<ethers.TransactionReceipt | null> {
      try {
        return await super.getTransactionReceipt(hash);
      } catch (error) {
        throw new Error(
          `Failed to get transaction receipt: ${
            error instanceof Error ? error.message : String(error)
          }`
        );
      }
    }

    override async getLogs(
      filter: ethers.Filter | ethers.FilterByBlockHash
    ): Promise<ethers.Log[]> {
      try {
        return await super.getLogs(filter);
      } catch (error) {
        throw new Error(
          `Failed to get logs: ${error instanceof Error ? error.message : String(error)}`
        );
      }
    }

    override async waitForTransaction(hash: string): Promise<ethers.TransactionReceipt | null> {
      try {
        return await super.waitForTransaction(hash);
      } catch (error) {
        throw new Error(
          `Failed to wait for transaction: ${
            error instanceof Error ? error.message : String(error)
          }`
        );
      }
    }
  },

  /**
   * Class representing an event log from a transaction
   */
  Log: ethers.Log,

  /**
   * Class representing a cryptographic signature
   * Provides wrapped constructor with error handling
   */
  Signature: class extends ethers.Signature {
    constructor(...args: ConstructorParameters<typeof ethers.Signature>) {
      try {
        super(...args);
      } catch (error) {
        throw new Error(
          `Failed to create Signature: ${error instanceof Error ? error.message : String(error)}`
        );
      }
    }
  },

  /**
   * Class for creating and verifying cryptographic signatures
   * Provides wrapped constructor with error handling
   */
  SigningKey: class extends ethers.SigningKey {
    constructor(...args: ConstructorParameters<typeof ethers.SigningKey>) {
      try {
        super(...args);
      } catch (error) {
        throw new Error(
          `Failed to create SigningKey: ${error instanceof Error ? error.message : String(error)}`
        );
      }
    }

    override sign(digest: string | Uint8Array): ethers.Signature {
      try {
        return super.sign(digest);
      } catch (error) {
        throw new Error(
          `Failed to sign message: ${error instanceof Error ? error.message : String(error)}`
        );
      }
    }
  },

  /**
   * Class representing an Ethereum transaction
   */
  Transaction: ethers.Transaction,

  /**
   * Class representing a receipt for a mined transaction
   */
  TransactionReceipt: ethers.TransactionReceipt,

  /**
   * Class for managing private/public key pairs and signing operations
   * Provides wrapped constructor with error handling
   */
  Wallet: class extends ethers.Wallet {
    constructor(...args: ConstructorParameters<typeof ethers.Wallet>) {
      try {
        super(...args);
      } catch (error) {
        throw new Error(
          `Failed to create Wallet: ${error instanceof Error ? error.message : String(error)}`
        );
      }
    }

    override connect(provider: ethers.Provider): ethers.Wallet {
      try {
        return super.connect(provider);
      } catch (error) {
        throw new Error(
          `Failed to connect wallet to provider: ${
            error instanceof Error ? error.message : String(error)
          }`
        );
      }
    }

    override signTransaction(tx: ethers.TransactionRequest): Promise<string> {
      try {
        return super.signTransaction(tx);
      } catch (error) {
        throw new Error(
          `Failed to sign transaction: ${error instanceof Error ? error.message : String(error)}`
        );
      }
    }

    override signMessage(message: string | Uint8Array): Promise<string> {
      try {
        return super.signMessage(message);
      } catch (error) {
        throw new Error(
          `Failed to sign message: ${error instanceof Error ? error.message : String(error)}`
        );
      }
    }
  },
};
