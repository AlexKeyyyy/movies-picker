package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getFirstMovieID(t *testing.T) int {
	resp, err := http.Get(fmt.Sprintf("%s/movies?page=1&size=1", baseURL))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var arr []struct {
		MovieID int `json:"movie_id"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&arr))
	require.NotEmpty(t, arr, "no movies available to test with")
	return arr[0].MovieID
}

func TestRatingsFlow(t *testing.T) {
	client := http.DefaultClient
	movieID := getFirstMovieID(t)

	t.Run("create rating", func(t *testing.T) {
		req := authRequest(t, http.MethodPost, fmt.Sprintf("%s/users/%d/ratings", baseURL, userID), map[string]int{"movie_id": movieID, "rating": 8})
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("read created rating", func(t *testing.T) {
		req := authRequest(t, http.MethodGet, fmt.Sprintf("%s/users/%d/ratings", baseURL, userID), nil)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var ratings []struct {
			MovieID int `json:"movie_id"`
			Rating  int `json:"rating"`
		}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&ratings))
		found := false
		for _, r := range ratings {
			if r.MovieID == movieID && r.Rating == 8 {
				found = true
				break
			}
		}
		assert.True(t, found)
	})

	t.Run("delete rating", func(t *testing.T) {
		req := authRequest(t, http.MethodDelete, fmt.Sprintf("%s/users/%d/ratings/%d", baseURL, userID, movieID), nil)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("create rating with invalid payload", func(t *testing.T) {
		req := authRequest(t, http.MethodPost, fmt.Sprintf("%s/users/%d/ratings", baseURL, userID), map[string]string{"movie_id": "bad", "rating": "bad"})
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
