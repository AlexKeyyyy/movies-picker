package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthFlow(t *testing.T) {
	// проверяем /users/me
	req, _ := http.NewRequest("GET", baseURL+"/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
