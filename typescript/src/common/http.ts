/**
 * A function interface that matches the Fetch API, for making HTTP requests.
 */
export interface HttpClient {
  (input: string | URL | Request, init?: RequestInit | undefined): Promise<Response>;
}
