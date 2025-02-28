import { Signer } from '../auth';
import { Receipt } from '../common';
import { Contract } from './contract';

/**
 * Client interface for interacting with smart contracts on the Radius platform.
 * Implemented by the main Radius Client.
 */
export interface ContractClient {
  /**
   * Calls a read-only contract method without creating a transaction
   * @param contract Contract instance to interact with
   * @param method Name of the method to call on the contract
   * @param args Arguments to pass to the contract method
   * @returns Array of decoded return values from the contract method
   * @throws Error if the contract ABI is missing
   * @throws Error if the contract address is missing or zero
   * @throws Error if the contract method call fails
   */
  call(contract: Contract, method: string, ...args: unknown[]): Promise<unknown[]>;

  /**
   * Executes a contract method that modifies Radius state
   * @param contract Contract instance to interact with
   * @param signer The signer used to sign the transaction
   * @param method Name of the method to execute on the contract
   * @param args Arguments to pass to the contract method
   * @returns Transaction receipt after the method execution
   * @throws Error if the contract ABI is missing
   * @throws Error if the contract address is missing or zero
   * @throws Error if the transaction fails or is reverted
   * @throws Error if the transaction receipt is not returned
   */
  execute(contract: Contract, signer: Signer, method: string, ...args: unknown[]): Promise<Receipt>;
}
