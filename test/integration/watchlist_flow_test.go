package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWatchlistFlow(t *testing.T) {
	client := http.DefaultClient

	// Получаем реальный movieID
	movieID := getFirstMovieID(t)

	// 1) Add to watchlist
	payload := map[string]int{"movie_id": movieID}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/users/%d/watchlist", baseURL, userID),
		bytes.NewReader(b),
	)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// 2) Get watchlist — должен быть наш movieID
	req2, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/users/%d/watchlist", baseURL, userID),
		nil,
	)
	req2.Header.Set("Authorization", "Bearer "+token)
	resp2, err := client.Do(req2)
	assert.NoError(t, err)
	defer resp2.Body.Close()
	assert.Equal(t, http.StatusOK, resp2.StatusCode)

	var items []struct {
		MovieID int `json:"movie_id"`
	}
	assert.NoError(t, json.NewDecoder(resp2.Body).Decode(&items))
	found := false
	for _, it := range items {
		if it.MovieID == movieID {
			found = true
			break
		}
	}
	assert.True(t, found, "added movie must appear in watchlist")

	// 3) Remove from watchlist
	req3, _ := http.NewRequest(
		"DELETE",
		fmt.Sprintf("%s/users/%d/watchlist/%d", baseURL, userID, movieID),
		nil,
	)
	req3.Header.Set("Authorization", "Bearer "+token)
	resp3, err := client.Do(req3)
	assert.NoError(t, err)
	defer resp3.Body.Close()
	assert.Equal(t, http.StatusNoContent, resp3.StatusCode)
}
