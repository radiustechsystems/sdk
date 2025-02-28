import { BigNumberish, eth } from '../providers/eth';
import { ABI } from './abi';
import { Address } from './address';
import { Event } from './event';
import { Hash } from './hash';
import { Receipt } from './receipt';

/**
 * Creates a new ABI (Application Binary Interface) from a JSON string
 * @param json ABI definition in JSON string format
 * @returns A new ABI instance, or undefined if the JSON is invalid
 */
export function abiFromJSON(json: string): ABI | undefined {
  try {
    return new ABI(json);
  } catch {
    return undefined;
  }
}

/**
 * Creates an Address from a hex string
 * @param hex Hex string with or without 0x prefix
 * @returns Address instance
 * @throws Error if the hex string is invalid
 */
export function addressFromHex(hex: string): Address {
  const cleanHex = hex.startsWith('0x') ? hex : `0x${hex}`;
  return new Address(eth.getBytes(cleanHex));
}

/**
 * Converts a hex string to a byte array
 * @param s Hex string (with or without 0x prefix)
 * @returns Byte array representation of the hex string, or undefined if the string is not valid hex
 */
export function bytecodeFromHex(s: string): Uint8Array | undefined {
  try {
    const cleanHex = s.startsWith('0x') ? s.slice(2) : s;
    return eth.getBytes(`0x${cleanHex}`);
  } catch {
    return undefined;
  }
}

/**
 * Converts a Radius Address to an Ethereum Address
 * @param address Radius Address
 * @returns Ethereum Address, or undefined if the input is undefined
 */
export function ethAddressFromRadiusAddress(address?: Address): string | undefined {
  if (!address) {
    return undefined;
  }
  return address.ethAddress();
}

/**
 * Converts Ethereum logs to Radius events
 * @param logs Ethereum logs
 * @returns Array of Radius events
 */
// biome-ignore lint/suspicious/noExplicitAny: Ethers.js does not export a type for Log
export function eventsFromEthLogs(logs: any[]): Event[] {
  return logs.map((log) => new Event(log.topics[0], {}, log.data));
}

/**
 * Creates a Hash from a hexadecimal string
 * @param hex The hexadecimal string (with or without 0x prefix)
 * @returns A new Hash instance
 * @throws Error if the hex string is invalid
 */
export function hashFromHex(hex: string): Hash {
  const cleanHex = hex.startsWith('0x') ? hex : `0x${hex}`;
  return new Hash(eth.getBytes(cleanHex));
}

/**
 * Creates a new Radius receipt from an Ethereum receipt
 * @param receipt Ethereum receipt
 * @param from Sender address
 * @param to Recipient address
 * @param value Transaction value
 * @returns Radius receipt
 */
export function receiptFromEthReceipt(
  // biome-ignore lint/suspicious/noExplicitAny: Ethers.js does not export a type for TransactionReceipt
  receipt: any,
  from: Address,
  to: Address = new Address(zeroAddress()),
  value?: BigNumberish
): Receipt {
  return new Receipt(
    from,
    to,
    new Address(receipt.contractAddress ?? zeroAddress()),
    new Hash(receipt.hash),
    receipt.gasUsed,
    receipt.status ?? 0,
    // biome-ignore lint/suspicious/noExplicitAny: Ethers.js does not export a type for Log
    eventsFromEthLogs(receipt.logs?.map((log: any) => log.toJSON()) ?? []),
    value
  );
}

/**
 * Creates a zero address (0x0000000000000000000000000000000000000000)
 * Used as a default value or to represent the zero address in the Ethereum ecosystem
 * @returns An Address instance representing the zero address
 */
export function zeroAddress(): Address {
  return new Address('0x'.padEnd(22, '0'));
}
