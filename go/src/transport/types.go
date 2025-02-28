// Package transport provides HTTP transport mechanisms for the Radius SDK.
// It includes interceptors and middleware for logging, debugging, and modifying
// JSON-RPC requests and responses.
package transport

import "net/http"

// Logf is a logging function interface that accepts a format string and arguments.
// It follows the standard fmt.Printf style interface pattern in Go.
//
// @param format The format string with placeholders
// @param args The values to substitute into the format string
type Logf func(format string, args ...any)

// Interceptor is a function interface used to intercept and modify HTTP requests and responses.
// This allows for custom handling, validation, or manipulation of JSON-RPC calls.
//
// @param reqBody The stringified JSON-RPC request body
// @param resp The HTTP response from the JSON-RPC server
// @return A potentially modified response or the original response
// @return An error if interceptor processing fails
type Interceptor func(reqBody string, resp *http.Response) (*http.Response, error)
