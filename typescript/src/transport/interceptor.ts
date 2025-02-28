import { JsonRpcPayload, JsonRpcResult, eth } from '../providers/eth';
import { Interceptor, Logf, RoundTripper } from './types';

/**
 * A RoundTripper implementation that intercepts HTTP requests and responses
 * Provides request logging and response modification capabilities
 * @implements {RoundTripper}
 */
export class InterceptingRoundTripper implements RoundTripper {
  /**
   * Creates a new InterceptingRoundTripper
   * @param interceptor Optional function to intercept and modify responses
   * @param logf Optional logging function to record requests and responses
   * @param proxied Underlying RoundTripper implementation (defaults to fetch-based implementation)
   */
  constructor(
    private readonly interceptor?: Interceptor,
    private readonly logf?: Logf,
    private readonly proxied: RoundTripper = new DefaultRoundTripper()
  ) {}

  /**
   * Sends a request and handles interception and logging of the response.
   *
   * The process flow is:
   * 1. Parse the request body
   * 2. Log the request if a logger is provided
   * 3. Send the request using the proxied RoundTripper
   * 4. Log the response if a logger is provided
   * 5. Pass the response through the interceptor if one is provided
   * 6. Return the final (potentially modified) response
   *
   * @param request The HTTP request to send
   * @returns The HTTP response, potentially modified by the interceptor
   * @throws Error if the request fails or if the interceptor throws an error
   */
  async roundTrip(request: Request): Promise<Response> {
    const reqBody = await this.parseRequestBody(request);

    if (this.logf) {
      this.logf('Request:', {
        url: request.url,
        method: request.method,
        body: reqBody,
      });
    }

    let response: Response;
    try {
      response = await this.proxied.roundTrip(request);
      const body = await response.clone().text();

      if (this.logf) {
        this.logf('Response:', {
          status: response.status,
          body,
        });
      }

      response = new Response(body, {
        status: response.status,
        statusText: response.statusText,
        headers: response.headers,
      });
    } catch (err) {
      if (this.logf) {
        this.logf('Request failed', {
          error: err instanceof Error ? err.message : String(err),
        });
      }
      throw err;
    }

    if (this.interceptor) {
      return this.interceptor(reqBody, response);
    }

    return response;
  }

  /**
   * Parse the body of a request, cloning the request to avoid modifying the original.
   * @param request The HTTP request containing the body to parse
   * @returns The request body as a string, or empty string if body is null
   * @throws Error if parsing the request body fails
   * @private
   */
  private async parseRequestBody(request: Request): Promise<string> {
    if (!request.body) {
      return '';
    }

    try {
      const clone = request.clone();
      return await clone.text();
    } catch (err) {
      throw new Error(`Failed to parse request body: ${err}`);
    }
  }
}

/**
 * A simple implementation of RoundTripper that uses the Fetch API.
 * This is used as the default implementation when no other RoundTripper is provided.
 * It serves as a thin wrapper around the standard fetch function to make it conform
 * to the RoundTripper interface.
 *
 * @implements {RoundTripper}
 * @private
 */
class DefaultRoundTripper implements RoundTripper {
  /**
   * Sends an HTTP request using the Fetch API and returns the response
   * @param request The HTTP request to send
   * @returns A Promise that resolves to the HTTP response
   * @throws Error if the fetch request fails
   */
  async roundTrip(request: Request): Promise<Response> {
    return fetch(request);
  }
}

/**
 * A JSON-RPC provider that uses an InterceptingRoundTripper to intercept and modify requests/responses
 * Extends ethers.js JsonRpcProvider to add logging and interception capabilities
 * @extends {eth.JsonRpcProvider}
 */
export class InterceptingProvider extends eth.JsonRpcProvider {
  /**
   * The round tripper used to send HTTP requests and receive responses
   * @private
   */
  private roundTripper: InterceptingRoundTripper;

  /**
   * Creates a new InterceptingProvider instance
   * @param url URL of the JSON-RPC endpoint (e.g., "http://localhost:8545")
   * @param roundTripper Optional custom InterceptingRoundTripper instance
   */
  constructor(url?: string, roundTripper?: InterceptingRoundTripper) {
    super(url);
    this.roundTripper = roundTripper || new InterceptingRoundTripper();
  }

  /**
   * Sends a JSON-RPC request via the InterceptingRoundTripper and returns the result
   * Overrides the _send method from JsonRpcProvider to use our custom transport
   * @param payload The JSON-RPC request payload (either a single request or a batch)
   * @returns An array of JSON-RPC response objects
   * @throws Error if the request fails or returns invalid JSON
   * @private
   */
  async _send(payload: JsonRpcPayload | Array<JsonRpcPayload>): Promise<Array<JsonRpcResult>> {
    // Get the connection URL
    const connection = this._getConnection();

    // Create the request
    const request = new Request(connection.url, {
      method: 'POST',
      headers: {
        'content-type': 'application/json',
      },
      body: JSON.stringify(payload),
    });

    // Send the request through our interceptor
    const response = await this.roundTripper.roundTrip(request);

    // Parse the response
    const result = await response.json();

    // Handle both single and batch responses
    return Array.isArray(result) ? result : [result];
  }
}
