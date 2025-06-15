package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// helper в начале файла (общий для обоих тестов)
func getFirstMovieID(t *testing.T) int {
	resp, err := http.Get(fmt.Sprintf("%s/movies?page=1&size=1", baseURL))
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var arr []struct {
		MovieID int `json:"movie_id"`
	}
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&arr))
	if len(arr) == 0 {
		t.Fatal("no movies available to test with")
	}
	return arr[0].MovieID
}

func TestRatingsFlow(t *testing.T) {
	client := http.DefaultClient

	// Получаем реальный movieID
	movieID := getFirstMovieID(t)

	// 1) Add rating
	rp := map[string]int{"movie_id": movieID, "rating": 8}
	b, _ := json.Marshal(rp)
	req, _ := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/users/%d/ratings", baseURL, userID),
		bytes.NewReader(b),
	)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// 2) Get ratings — ожидаем наш rating
	req2, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/users/%d/ratings", baseURL, userID),
		nil,
	)
	req2.Header.Set("Authorization", "Bearer "+token)
	resp2, err := client.Do(req2)
	assert.NoError(t, err)
	defer resp2.Body.Close()
	assert.Equal(t, http.StatusOK, resp2.StatusCode)

	var ratings []struct {
		MovieID int `json:"movie_id"`
		Rating  int `json:"rating"`
	}
	assert.NoError(t, json.NewDecoder(resp2.Body).Decode(&ratings))
	found := false
	for _, r := range ratings {
		if r.MovieID == movieID && r.Rating == 8 {
			found = true
			break
		}
	}
	assert.True(t, found, "added rating must be present")
}
