package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func authRequest(t *testing.T, method, url string, body any) *http.Request {
	t.Helper()

	var payload []byte
	if body != nil {
		var err error
		payload, err = json.Marshal(body)
		require.NoError(t, err)
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(payload))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req
}
