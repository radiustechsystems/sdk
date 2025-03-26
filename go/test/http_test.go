package test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/radiustechsystems/sdk/go/src/radius"
)

const testURL = "https://radiustech.xyz"

func TestInterceptingRoundTripper(t *testing.T) {
	t.Run("Basic usage", func(t *testing.T) {
		expectedResult := "test response"
		mockResp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(expectedResult)),
		}

		mockTransport := mockRoundTripper{
			response: mockResp,
			err:      nil,
		}

		transport := radius.InterceptingRoundTripper{
			Proxied: mockTransport,
		}

		req, _ := http.NewRequest("GET", testURL, bytes.NewBufferString("test request"))

		resp, err := transport.RoundTrip(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, expectedResult, string(body))
	})

	t.Run("With logging", func(t *testing.T) {
		expectedResult := "test response"
		mockResp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(expectedResult)),
		}

		mockTransport := mockRoundTripper{
			response: mockResp,
			err:      nil,
		}

		var logMessages []string
		logger := func(format string, args ...any) {
			logMessages = append(logMessages, fmt.Sprintf(format, args...))
		}

		transport := radius.InterceptingRoundTripper{
			Proxied: mockTransport,
			Log:     logger,
		}

		requestBody := "test request"
		req, _ := http.NewRequest("GET", testURL, bytes.NewBufferString(requestBody))

		resp, err := transport.RoundTrip(req)
		require.NoError(t, err)
		assert.Len(t, logMessages, 2, "should have logged request and response")

		expectedReqLog := fmt.Sprintf("Request to %s: %s", req.URL, requestBody)
		assert.Contains(t, logMessages[0], expectedReqLog)

		expectedRespLog := fmt.Sprintf("Response from %s: %s", req.URL, expectedResult)
		assert.Contains(t, logMessages[1], expectedRespLog)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, expectedResult, string(body))
	})

	t.Run("With interceptor", func(t *testing.T) {
		originalResponseBody := "original response"
		mockResp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(originalResponseBody)),
		}

		mockTransport := mockRoundTripper{
			response: mockResp,
			err:      nil,
		}

		modifiedResponseBody := "modified response"
		interceptor := func(_ string, resp *http.Response) (*http.Response, error) {
			body, _ := io.ReadAll(resp.Body)
			assert.Equal(t, originalResponseBody, string(body))
			modifiedResp := &http.Response{
				StatusCode: http.StatusTeapot,
				Body:       io.NopCloser(bytes.NewBufferString(modifiedResponseBody)),
			}
			return modifiedResp, nil
		}

		transport := radius.InterceptingRoundTripper{
			Proxied:     mockTransport,
			Interceptor: interceptor,
		}

		req, _ := http.NewRequest("GET", testURL, bytes.NewBufferString("test request"))

		resp, err := transport.RoundTrip(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusTeapot, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, modifiedResponseBody, string(body))
	})

	t.Run("With proxied error", func(t *testing.T) {
		expectedErr := errors.New("transport error")
		mockTransport := mockRoundTripper{
			response: nil,
			err:      expectedErr,
		}

		transport := radius.InterceptingRoundTripper{
			Proxied: mockTransport,
		}

		req, _ := http.NewRequest("GET", testURL, nil)

		resp, err := transport.RoundTrip(req)
		assert.Error(t, expectedErr, fmt.Sprintf("expected error: %v", err))
		assert.Nil(t, resp)
	})

	t.Run("With interceptor error", func(t *testing.T) {
		mockResp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString("test")),
		}

		mockTransport := mockRoundTripper{
			response: mockResp,
			err:      nil,
		}

		expectedErr := errors.New("interceptor error")
		interceptor := func(_ string, _ *http.Response) (*http.Response, error) {
			return nil, expectedErr
		}

		transport := radius.InterceptingRoundTripper{
			Proxied:     mockTransport,
			Interceptor: interceptor,
		}

		req, _ := http.NewRequest("GET", testURL, nil)

		resp, err := transport.RoundTrip(req)
		assert.Error(t, expectedErr, fmt.Sprintf("expected error: %v", err))
		assert.Nil(t, resp)
	})
}
