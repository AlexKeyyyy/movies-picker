package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWatchlistFlow(t *testing.T) {
	client := http.DefaultClient
	movieID := getFirstMovieID(t)

	t.Run("add movie to watchlist", func(t *testing.T) {
		req := authRequest(t, http.MethodPost, fmt.Sprintf("%s/users/%d/watchlist", baseURL, userID), map[string]int{"movie_id": movieID})
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("watchlist contains added movie", func(t *testing.T) {
		req := authRequest(t, http.MethodGet, fmt.Sprintf("%s/users/%d/watchlist", baseURL, userID), nil)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var items []struct {
			MovieID int `json:"movie_id"`
		}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&items))
		found := false
		for _, it := range items {
			if it.MovieID == movieID {
				found = true
				break
			}
		}
		assert.True(t, found)
	})

	t.Run("remove movie from watchlist", func(t *testing.T) {
		req := authRequest(t, http.MethodDelete, fmt.Sprintf("%s/users/%d/watchlist/%d", baseURL, userID, movieID), nil)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("add to watchlist with invalid payload", func(t *testing.T) {
		req := authRequest(t, http.MethodPost, fmt.Sprintf("%s/users/%d/watchlist", baseURL, userID), map[string]string{"movie_id": "oops"})
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
