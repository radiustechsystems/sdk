import { HttpClient } from '../common';
import { Interceptor, Logf } from '../transport';

/**
 * A function that configures the Radius Client with a given set of options
 * Used for functional-style configuration of the client
 *
 * @param options The options object to modify with configuration settings
 */
export type ClientOption = (options: ClientOptions) => void;

/**
 * Configuration options for the Radius Client
 */
export interface ClientOptions {
  /**
   * Custom HTTP client implementation for making API requests
   * If not provided, the global fetch function will be used
   */
  httpClient?: HttpClient;

  /**
   * Response interceptor for modifying or monitoring JSON-RPC responses
   * Can be used to implement custom error handling or response transformation
   */
  interceptor?: Interceptor;

  /**
   * Logger function for debugging request/response cycles
   * Useful for troubleshooting connection or parsing issues
   */
  logger?: Logf;
}

/**
 * Creates an option to set a custom HTTP client for the Radius Client
 * By default, the global fetch function is used for HTTP requests
 *
 * @param httpClient Custom HTTP client implementing the HttpClient interface
 * @returns A ClientOption function that can be passed to Client.New()
 */
export function withHttpClient(httpClient: HttpClient): ClientOption {
  return (options: ClientOptions) => {
    options.httpClient = httpClient;
  };
}

/**
 * Creates an option to set a response interceptor for the Radius Client.
 * This can be used to log, modify, or analyze responses from the Radius server.
 * It's useful for debugging, testing, and to temporarily patch any issues in JSON-RPC responses.
 *
 * @param interceptor Function that can intercept and potentially modify JSON-RPC responses
 * @returns A ClientOption function that can be passed to Client.New()
 */
export function withInterceptor(interceptor: Interceptor): ClientOption {
  return (options: ClientOptions) => {
    options.interceptor = interceptor;
  };
}

/**
 * Creates an option to set a logger for the Radius Client.
 * This can be used to log JSON-RPC requests and responses for debugging or audit purposes.
 * The logger receives the raw request and response bodies for inspection.
 *
 * @param logger Function that logs messages with format strings and variable arguments
 * @returns A ClientOption function that can be passed to Client.New()
 */
export function withLogger(logger: Logf): ClientOption {
  return (options: ClientOptions) => {
    options.logger = logger;
  };
}
