package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMovieEndpoints(t *testing.T) {
	t.Run("list movies with pagination", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/movies?page=1&size=5", baseURL))
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var movies []struct {
			MovieID int     `json:"movie_id"`
			Title   string  `json:"title"`
			Rating  float64 `json:"ratingKinopoisk"`
		}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&movies))
		assert.NotEmpty(t, movies)
		assert.LessOrEqual(t, len(movies), 5)
	})

	t.Run("get movie by existing id", func(t *testing.T) {
		id := getFirstMovieID(t)

		resp, err := http.Get(fmt.Sprintf("%s/movies/%d", baseURL, id))
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var m struct {
			MovieID int `json:"movie_id"`
		}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&m))
		assert.Equal(t, id, m.MovieID)
	})

	t.Run("get movie by invalid id", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/movies/not-a-number", baseURL))
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("search with missing query", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/movies/search", baseURL))
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
