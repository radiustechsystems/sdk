/**
 * Exports provider types and utilities for use in the Radius SDK.
 *
 * This module serves as an abstraction layer for Ethereum-compatible libraries.
 * While the Radius SDK has its own concrete implementations of core data structures
 * like addresses and transactions, this module maps those structures to Ethereum types
 * to leverage the functionality provided by established Ethereum libraries.
 *
 * This approach allows the SDK to benefit from well-tested Ethereum libraries
 * while maintaining its own independent implementations that are optimized for
 * Radius's custom database architecture for parallel transaction processing.
 */
export * from './types';
export * from './utils';
