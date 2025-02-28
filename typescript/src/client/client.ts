import { TransactionReceipt, TransactionResponse } from 'ethers';
import { Signer } from '../auth';
import {
  ABI,
  Address,
  HttpClient,
  MAX_GAS,
  Receipt,
  SignedTransaction,
  Transaction,
  receiptFromEthReceipt,
  zeroAddress,
} from '../common';
import { Contract } from '../contracts';
import { BigNumberish, Provider, eth } from '../providers/eth';
import { InterceptingProvider, InterceptingRoundTripper } from '../transport';
import { ClientOption, ClientOptions } from './options';

/**
 * Client used to interact with the Radius platform.
 * This is the main entry point for working with the Radius ecosystem.
 * It provides methods for account management, contract deployment and interaction,
 * transaction handling, and querying Radius state.
 */
export class Client {
  /**
   * The Ethereum JSON-RPC provider used to communicate with Radius
   * @private
   */
  private readonly ethClient: Provider;

  /**
   * The HTTP client used for making API requests
   * @private
   */
  private readonly _httpClient: HttpClient;

  /**
   * Creates a new Radius Client instance
   * @param provider The Ethereum provider to use for Radius communication
   * @param httpClient Optional HTTP client to use for API requests
   */
  constructor(provider: Provider, httpClient?: HttpClient) {
    this.ethClient = provider;
    this._httpClient = httpClient ?? globalThis.fetch;
  }

  /**
   * Create a new Radius Client with the given URL and ClientOption(s).
   * @param url URL of the Radius node
   * @param opts ClientOption(s)
   * @returns New Radius Client
   * @throws Error if the client cannot be created
   */
  static async New(url: string, ...opts: ClientOption[]): Promise<Client> {
    const options: ClientOptions = {
      httpClient: globalThis.fetch,
    };

    for (const opt of opts) {
      opt(options);
    }

    // Create a new provider with the given URL and an optional HTTP interceptor and logger
    const provider =
      options.logger || options.interceptor
        ? new InterceptingProvider(
            url,
            new InterceptingRoundTripper(options.interceptor, options.logger, {
              roundTrip: (req) => (options.httpClient ? options.httpClient(req) : fetch(req)),
            })
          )
        : new eth.JsonRpcProvider(url);

    // Ensure the provider is connected to the network
    try {
      await provider.getNetwork();
    } catch (error) {
      throw new Error(
        `Failed to create Radius client: ${error instanceof Error ? error.message : String(error)}`
      );
    }

    return new Client(provider, options.httpClient);
  }

  /**
   * Gets the balance of an account in wei
   * @param address The address to check the balance for
   * @returns The account balance in wei as a bigint
   * @throws Error if the balance cannot be retrieved from the network
   */
  async balanceAt(address: Address): Promise<bigint> {
    try {
      return this.ethClient.getBalance(address.ethAddress());
    } catch (error) {
      throw new Error(
        `Failed to get balance: ${error instanceof Error ? error.message : String(error)}`
      );
    }
  }

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
  async call(contract: Contract, method: string, ...args: unknown[]): Promise<unknown[]> {
    if (!contract.abi) {
      throw new Error('Contract ABI is required');
    }
    if (!contract.address() || contract.address().equals(zeroAddress())) {
      throw new Error('Contract address is required');
    }

    const data = contract.abi.pack(method, ...args);
    const params = new TxParams(data, undefined, contract.address()); // No signer needed here
    const tx = await this.prepareTx(params);

    let resultData: Uint8Array;
    try {
      const result = await this.ethClient.call({
        to: tx.to?.ethAddress(),
        data: tx.data ? eth.hexlify(tx.data) : undefined,
        value: tx.value,
      });
      resultData = eth.getBytes(result);
    } catch (error) {
      throw new Error(
        `Failed to call contract method: ${error instanceof Error ? error.message : String(error)}`
      );
    }

    return contract.abi.unpack(method, resultData);
  }

  /**
   * Get the chain ID of the connected Radius network
   * @returns Chain ID of the Radius network
   * @throws Error if the chain ID cannot be retrieved
   */
  async chainID(): Promise<bigint> {
    try {
      const network = await this.ethClient.getNetwork();
      return network.chainId;
    } catch (error) {
      throw new Error(
        `Failed to get chain ID: ${error instanceof Error ? error.message : String(error)}`
      );
    }
  }

  /**
   * Get the bytecode of a contract
   * @param address Address of the contract
   * @returns Bytecode of the contract
   */
  async codeAt(address: Address): Promise<Uint8Array> {
    return eth.getBytes(await this.ethClient.getCode(address.ethAddress()));
  }

  /**
   * Deploy a smart contract.
   * @param signer The signer that should be used to sign the transaction
   * @param bytecode Bytecode of the contract
   * @param abi ABI of the contract
   * @param args Arguments for the contract constructor
   * @throws Error if the contract bytecode is not provided
   * @throws Error if the contract deployment fails
   */
  async deployContract(
    signer: Signer,
    bytecode: Uint8Array,
    abi: ABI,
    ...args: unknown[]
  ): Promise<Contract> {
    if (signer === undefined) {
      throw new Error('Signer is required for deploying contracts');
    }

    const data = bytecode;
    if (args.length > 0) {
      const encodedConstructorArgs = abi.pack('', ...args);
      data.set(encodedConstructorArgs, bytecode.length);
    }

    let receipt: Receipt;
    try {
      const params = new TxParams(data, signer);
      receipt = await this.prepareAndSendTx(params);
    } catch (error) {
      throw new Error(
        `Failed to deploy contract: ${error instanceof Error ? error.message : String(error)}`
      );
    }
    if (!receipt) {
      throw new Error('Contract deployment failed: no receipt returned');
    }
    if (receipt.status !== 1) {
      throw new Error(
        `Failed to deploy contract: status ${receipt.status}, transaction hash ${receipt.txHash}`
      );
    }

    return new Contract(receipt.contractAddress, abi);
  }

  /**
   * Estimate gas for a transaction.
   * @param tx Transaction to estimate gas for
   */
  async estimateGas(tx: Transaction): Promise<bigint> {
    const estimate = await this.ethClient.estimateGas({
      to: tx.to?.ethAddress(),
      data: tx.data ? eth.hexlify(tx.data) : undefined,
      value: tx.value,
    });

    // Apply 20% safety margin
    const margin = estimate / BigInt('5');
    const gas = estimate + margin;

    // Cap at MAX_GAS
    return gas > MAX_GAS ? MAX_GAS : gas;
  }

  /**
   * Execute a contract state-changing method.
   * @param contract Contract to execute
   * @param signer The signer that should be used to sign the transaction
   * @param method Method to execute
   * @param args Arguments for the method
   * @returns Receipt of the transaction
   * @throws Error if the transaction fails
   * @throws Error if the transaction receipt is not returned
   */
  async execute(
    contract: Contract,
    signer: Signer,
    method: string,
    ...args: unknown[]
  ): Promise<Receipt> {
    if (!contract.abi) {
      throw new Error('Contract ABI is required');
    }
    if (!contract.address() || contract.address().equals(zeroAddress())) {
      throw new Error('Contract address is required');
    }

    const data = contract.abi.pack(method, ...args);

    return this.prepareAndSendTx({
      signer,
      to: contract.address(),
      data,
      value: BigInt('0'),
    });
  }

  /**
   * Get the HTTP client used by the client.
   * @returns HTTP client
   */
  httpClient(): HttpClient {
    return this._httpClient;
  }

  /**
   * Get the next nonce for an account.
   * @param address Address of the account
   * @returns Nonce of the account
   */
  async pendingNonceAt(address: Address): Promise<number> {
    return this.ethClient.getTransactionCount(address.ethAddress(), 'pending');
  }

  /**
   * Send value to an account.
   * @param signer The signer that should be used to sign the transaction
   * @param recipient Address of the recipient
   * @param value Value to send
   * @returns Receipt of the transaction
   * @throws Error if the transaction fails
   * @throws Error if the transaction receipt is not returned
   */
  async send(signer: Signer, recipient: Address, value: BigNumberish): Promise<Receipt> {
    const data = new Uint8Array();
    const params = new TxParams(data, signer, recipient, value); // Nonce gets set in prepareTx
    const receipt = await this.prepareAndSendTx(params);

    if (!receipt.status) {
      throw new Error('Transaction failed');
    }

    return receipt;
  }

  /**
   * Send a signed transaction to Radius.
   * @param signer The signer that should be used to sign the transaction
   * @param tx Signed transaction
   * @returns Receipt of the transaction
   * @throws Error if the transaction fails
   * @throws Error if the transaction receipt is not returned
   */
  async transact(signer: Signer, tx: SignedTransaction): Promise<Receipt> {
    let response: TransactionResponse;
    try {
      response = await this.ethClient.broadcastTransaction(eth.hexlify(tx.serialized));
    } catch (error) {
      throw new Error(
        `Failed to send transaction: ${error instanceof Error ? error.message : String(error)}`
      );
    }

    let receipt: TransactionReceipt | null;
    try {
      receipt = await response.wait();
    } catch (error) {
      throw new Error(
        `Failed to get transaction receipt: ${
          error instanceof Error ? error.message : String(error)
        }`
      );
    }

    if (!receipt) {
      throw new Error('Failed to get transaction receipt: no receipt returned');
    }
    if (receipt.status !== 1) {
      throw new Error(
        `Transaction failed: status ${receipt.status}, transaction hash ${receipt.hash}`
      );
    }

    const from: Address = signer.address();
    const to: Address = receipt.to ? new Address(receipt.to) : zeroAddress();
    const value: BigNumberish | undefined = tx.value;

    return receiptFromEthReceipt(receipt, from, to, value);
  }

  /**
   * Prepare a transaction with correct nonce and gas.
   * @private
   * @param params Transaction parameters
   * @returns Prepared transaction
   */
  private async prepareTx(params: TxParams): Promise<Transaction> {
    const gas: bigint = BigInt('0');
    const gasPrice = BigInt('0');

    // Get the pending nonce for the signer address, if necessary
    const nonce = params.signer ? await this.pendingNonceAt(params.signer.address()) : undefined;

    // Must set Transaction.to value to undefined if it is the zero address
    if (params.to === zeroAddress()) {
      params.to = undefined;
    }

    // Create the initial transaction used to estimate gas
    const tx = new Transaction(params.data, gas, gasPrice, nonce, params.to, params.value);

    // Estimate gas cost for the transaction
    tx.gas = await this.estimateGas(tx);

    return tx;
  }

  /**
   * Prepare and send a transaction.
   * @private
   * @param params Transaction parameters
   * @returns Receipt of the transaction
   * @throws Error if the account is not provided
   * @throws Error if the transaction fails
   * @throws Error if the transaction receipt is not returned
   */
  private async prepareAndSendTx(params: TxParams): Promise<Receipt> {
    if (!params.signer) {
      throw new Error('Signer is required for sending transactions');
    }

    const tx = await this.prepareTx(params);
    const signedTx = await params.signer.signTransaction(tx);

    return this.transact(params.signer, signedTx);
  }
}

/**
 * Parameters used to prepare a transaction.
 * This is an internal class used by the Client for transaction preparation.
 */
class TxParams {
  /**
   * The transaction data (bytecode for contract creation or method call data)
   * @private
   */
  data: Uint8Array;

  /**
   * The signer used to sign the transaction
   * @private
   */
  signer?: Signer;

  /**
   * The destination address for the transaction (undefined for contract creation)
   * @private
   */
  to?: Address;

  /**
   * The amount of native currency to send with the transaction
   * @private
   */
  value?: BigNumberish;

  /**
   * Creates a new transaction parameters object
   * @param data The transaction data (bytecode or method call data)
   * @param signer Optional signer to sign the transaction
   * @param to Optional destination address (undefined for contract creation)
   * @param value Optional amount of native currency to send
   */
  constructor(data: Uint8Array, signer?: Signer, to?: Address, value?: BigNumberish) {
    this.data = data;
    this.signer = signer;
    this.to = to;
    this.value = value;
  }
}
