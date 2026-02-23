package system

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUserJourneyE2E(t *testing.T) {
	baseURL := os.Getenv("API_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	email := fmt.Sprintf("e2e_%d@example.com", time.Now().UnixNano())
	password := "pass123"
	credentials := map[string]string{"email": email, "password": password}

	token, userID := registerAndLogin(t, baseURL, credentials)
	movieID := firstMovieID(t, baseURL)

	t.Run("profile available", func(t *testing.T) {
		req := authReq(t, http.MethodGet, baseURL+"/users/me", token, nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("watchlist lifecycle", func(t *testing.T) {
		resp := do(t, authReq(t, http.MethodPost, fmt.Sprintf("%s/users/%d/watchlist", baseURL, userID), token, map[string]int{"movie_id": movieID}))
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		resp.Body.Close()

		resp = do(t, authReq(t, http.MethodGet, fmt.Sprintf("%s/users/%d/watchlist", baseURL, userID), token, nil))
		require.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()

		resp = do(t, authReq(t, http.MethodDelete, fmt.Sprintf("%s/users/%d/watchlist/%d", baseURL, userID, movieID), token, nil))
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
		resp.Body.Close()
	})

	t.Run("ratings lifecycle", func(t *testing.T) {
		resp := do(t, authReq(t, http.MethodPost, fmt.Sprintf("%s/users/%d/ratings", baseURL, userID), token, map[string]int{"movie_id": movieID, "rating": 9}))
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		resp.Body.Close()

		resp = do(t, authReq(t, http.MethodGet, fmt.Sprintf("%s/users/%d/ratings", baseURL, userID), token, nil))
		require.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()

		resp = do(t, authReq(t, http.MethodDelete, fmt.Sprintf("%s/users/%d/ratings/%d", baseURL, userID, movieID), token, nil))
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
		resp.Body.Close()
	})

	t.Run("unauthorized access blocked", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, baseURL+"/users/me", nil)
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func do(t *testing.T, req *http.Request) *http.Response {
	t.Helper()
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

func registerAndLogin(t *testing.T, baseURL string, credentials map[string]string) (string, int64) {
	body, _ := json.Marshal(credentials)
	registerResp, err := http.Post(baseURL+"/auth/register", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer registerResp.Body.Close()
	require.Equal(t, http.StatusCreated, registerResp.StatusCode)

	var reg struct {
		UserID int64 `json:"user_id"`
	}
	require.NoError(t, json.NewDecoder(registerResp.Body).Decode(&reg))

	loginResp, err := http.Post(baseURL+"/auth/login", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer loginResp.Body.Close()
	require.Equal(t, http.StatusOK, loginResp.StatusCode)

	var login struct {
		AccessToken string `json:"access_token"`
	}
	require.NoError(t, json.NewDecoder(loginResp.Body).Decode(&login))
	return login.AccessToken, reg.UserID
}

func firstMovieID(t *testing.T, baseURL string) int {
	resp, err := http.Get(fmt.Sprintf("%s/movies?page=1&size=1", baseURL))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var movies []struct {
		MovieID int `json:"movie_id"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&movies))
	require.NotEmpty(t, movies)
	return movies[0].MovieID
}

func authReq(t *testing.T, method, url, token string, payload any) *http.Request {
	t.Helper()
	var body []byte
	if payload != nil {
		var err error
		body, err = json.Marshal(payload)
		require.NoError(t, err)
	}
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req
}
