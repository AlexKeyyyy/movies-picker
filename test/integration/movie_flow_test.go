package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMovieEndpoints(t *testing.T) {
	// ListMovies
	resp, err := http.Get(fmt.Sprintf("%s/movies?page=1&size=5", baseURL))
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var movies []struct {
		MovieID int     `json:"movie_id"`
		Title   string  `json:"title"`
		Rating  float64 `json:"ratingKinopoisk"`
	}
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&movies))
	assert.LessOrEqual(t, len(movies), 5)

	// GetMovieByID
	id := movies[0].MovieID
	resp2, err2 := http.Get(fmt.Sprintf("%s/movies/%d", baseURL, id))
	assert.NoError(t, err2)
	defer resp2.Body.Close()
	assert.Equal(t, http.StatusOK, resp2.StatusCode)

	var m struct {
		MovieID int `json:"movie_id"`
	}
	assert.NoError(t, json.NewDecoder(resp2.Body).Decode(&m))
	assert.Equal(t, id, m.MovieID)
}
