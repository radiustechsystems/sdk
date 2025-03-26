package test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func MockJSONRPCServer(t *testing.T, handlers map[string]func(params []interface{}) interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err, "Failed to read request body")
		defer func() {
			err = r.Body.Close()
			require.NoError(t, err, "Failed to close request body")
		}()

		var request struct {
			JSONRPC string        `json:"jsonrpc"`
			ID      interface{}   `json:"id"`
			Method  string        `json:"method"`
			Params  []interface{} `json:"params"`
		}
		err = json.Unmarshal(body, &request)
		require.NoError(t, err, "Failed to parse JSON-RPC request")

		response := struct {
			JSONRPC string      `json:"jsonrpc"`
			ID      interface{} `json:"id"`
			Result  interface{} `json:"result,omitempty"`
			Error   interface{} `json:"error,omitempty"`
		}{
			JSONRPC: "2.0",
			ID:      request.ID,
		}

		if handler, ok := handlers[request.Method]; ok {
			result := handler(request.Params)
			response.Result = result
		} else if request.Method == "eth_chainId" {
			// Return testnet chain ID by default
			response.Result = TestnetChainIDHex
		} else {
			response.Error = map[string]interface{}{
				"code":    -32601,
				"message": "Method not found",
			}
		}

		w.Header().Set("Content-Type", "application/json")
		responseBytes, err := json.Marshal(response)
		require.NoError(t, err, "Failed to marshal JSON-RPC response")
		_, err = w.Write(responseBytes)
		require.NoError(t, err, "Mock server failed to write response")
	}))
}

type mockRoundTripper struct {
	response *http.Response
	err      error
}

func (m mockRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return m.response, m.err
}
