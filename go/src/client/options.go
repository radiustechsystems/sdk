package client

import (
	"net/http"

	"github.com/radiustechsystems/sdk/go/src/transport"
)

// Option is a functional option for configuring a new Radius Client.
// It follows the functional options pattern for clean and extensible API configuration.
type Option func(*Options)

// Options contains configuration options for a new Radius Client.
// These options control how the client connects to and interacts with the Radius node.
type Options struct {
	// httpClient is the HTTP client used for making API requests
	httpClient *http.Client

	// interceptor is a function for modifying or monitoring JSON-RPC responses
	interceptor transport.Interceptor

	// logger is a function for debugging request/response cycles
	logger transport.Logf
}

// WithHTTPClient creates an option to set a custom HTTP client for the Radius Client.
// By default, the standard http.Client is used for HTTP requests.
//
// @param client Custom HTTP client implementing the http.Client interface
// @return An Option function that can be passed to New()
func WithHTTPClient(client *http.Client) Option {
	return func(o *Options) {
		o.httpClient = client
	}
}

// WithInterceptor creates an option to set a response interceptor for the Radius Client.
// This can be used to log, modify, or analyze responses from the Radius server.
// It's useful for debugging, testing, and to temporarily patch any issues in JSON-RPC responses.
//
// @param interceptor Function that can intercept and potentially modify JSON-RPC responses
// @return An Option function that can be passed to New()
func WithInterceptor(interceptor transport.Interceptor) Option {
	return func(o *Options) {
		o.interceptor = interceptor
	}
}

// WithLogger creates an option to set a logger for the Radius Client.
// This can be used to log JSON-RPC requests and responses for debugging or audit purposes.
// The logger receives the raw request and response bodies for inspection.
//
// @param logger Function that logs messages with format strings and variable arguments
// @return An Option function that can be passed to New()
func WithLogger(logger transport.Logf) Option {
	return func(o *Options) {
		o.logger = logger
	}
}
