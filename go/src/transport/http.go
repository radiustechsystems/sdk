// Package transport provides HTTP transport mechanisms for the Radius SDK.
// It includes interceptors and middleware for logging, debugging, and modifying
// JSON-RPC requests and responses.
package transport

import (
	"bytes"
	"io"
	"net/http"
)

// InterceptingRoundTripper is a http.RoundTripper implementation that intercepts HTTP requests and responses.
// It can be used to log, analyze, and even modify requests and responses between the Radius client and server.
// This is useful for debugging, testing, and to temporarily patch any issues in the JSON-RPC communication.
type InterceptingRoundTripper struct {
	// Interceptor is an optional function to intercept and modify responses
	Interceptor Interceptor

	// Logf is an optional logging function to record requests and responses
	Logf Logf

	// Proxied is the underlying RoundTripper that will actually send the request
	Proxied http.RoundTripper
}

// RoundTrip implements the http.RoundTripper interface for sending HTTP requests.
// It handles the interception, logging, and optional modification of requests and responses.
//
// @param req The HTTP request to send
// @return The HTTP response and nil error on success
// @return nil and error if the request fails or interceptor processing fails
func (irt InterceptingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var err error

	// Clone the request body so it can be read again
	reqBody := parseRequestBody(req)

	if irt.Logf != nil {
		irt.Logf("Request to %s: %s", req.URL, reqBody)
	}

	// Make the actual request
	resp, err := irt.Proxied.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// Clone the response body so it can be read again
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Log the response body
	if irt.Logf != nil {
		irt.Logf("Response from %s: %s", req.URL, string(body))
	}

	// Set the response body back to its original state so it can be read again
	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	if irt.Interceptor != nil {
		return irt.Interceptor(reqBody, resp)
	}

	return resp, nil
}

// parseRequestBody reads the request body and returns it as a string.
// It also resets the request body so it can be read again by subsequent handlers.
//
// @param req The HTTP request containing the body to parse
// @return The request body as a string, or empty string if body is nil or reading fails
func parseRequestBody(req *http.Request) string {
	if req.Body == nil {
		return ""
	}

	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		return ""
	}

	req.Body = io.NopCloser(bytes.NewBuffer(reqBody))

	return string(reqBody)
}
