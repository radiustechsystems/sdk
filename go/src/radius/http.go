package radius

import (
	"bytes"
	"io"
	"net/http"
)

type Logger func(format string, args ...any)

type Interceptor func(reqBody string, resp *http.Response) (*http.Response, error)

type InterceptingRoundTripper struct {
	Interceptor Interceptor
	Log         Logger
	Proxied     http.RoundTripper
}

func (irt InterceptingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var err error

	// Clone the request body so it can be read again
	reqBody := parseRequestBody(req)

	if irt.Log != nil {
		irt.Log("Request to %s: %s", req.URL, reqBody)
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
	if irt.Log != nil {
		irt.Log("Response from %s: %s", req.URL, string(body))
	}

	// Set the response body back to its original state so it can be read again
	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	if irt.Interceptor != nil {
		return irt.Interceptor(reqBody, resp)
	}

	return resp, nil
}

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
