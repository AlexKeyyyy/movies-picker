package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// INT-21: Добавление/обновление рейтинга фильма
func TestINT21_Ratings_CreateRating(t *testing.T) {
    movieID := getFirstMovieID(t)

    req := authRequest(t, http.MethodPost,
        fmt.Sprintf("%s/users/%d/ratings", baseURL, userID),
        map[string]int{"movie_id": movieID, "rating": 8},
    )
    resp := doRequest(t, req)
    defer resp.Body.Close()

    require.Equal(t, http.StatusCreated, resp.StatusCode,
        "создание рейтинга должно вернуть 201")

    var result struct {
        MovieID int `json:"movie_id"`
        Rating  int `json:"rating"`
    }
    require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
    assert.Equal(t, movieID, result.MovieID)
    assert.Equal(t, 8, result.Rating)
}

// INT-22: Получение списка рейтингов пользователя
func TestINT22_Ratings_GetList(t *testing.T) {
    movieID := getFirstMovieID(t)

    // Убеждаемся что рейтинг выставлен
    addReq := authRequest(t, http.MethodPost,
        fmt.Sprintf("%s/users/%d/ratings", baseURL, userID),
        map[string]int{"movie_id": movieID, "rating": 9},
    )
    addResp := doRequest(t, addReq)
    addResp.Body.Close()

    // Получаем список
    req := authRequest(t, http.MethodGet,
        fmt.Sprintf("%s/users/%d/ratings", baseURL, userID), nil)
    resp := doRequest(t, req)
    defer resp.Body.Close()

    require.Equal(t, http.StatusOK, resp.StatusCode)

    var ratings []struct {
        MovieID int `json:"movie_id"`
        Rating  int `json:"rating"`
    }
    require.NoError(t, json.NewDecoder(resp.Body).Decode(&ratings))
    assert.NotEmpty(t, ratings, "список рейтингов не должен быть пустым")

    found := false
    for _, r := range ratings {
        if r.MovieID == movieID {
            found = true
            assert.Equal(t, 9, r.Rating, "рейтинг должен быть последним выставленным (9)")
            break
        }
    }
    assert.True(t, found, "выставленный рейтинг должен присутствовать в списке")
}

// INT-23: Удаление рейтинга
func TestINT23_Ratings_DeleteRating(t *testing.T) {
    movieID := getFirstMovieID(t)

    // Добавляем рейтинг
    addReq := authRequest(t, http.MethodPost,
        fmt.Sprintf("%s/users/%d/ratings", baseURL, userID),
        map[string]int{"movie_id": movieID, "rating": 7},
    )
    addResp := doRequest(t, addReq)
    addResp.Body.Close()

    // Удаляем
    delReq := authRequest(t, http.MethodDelete,
        fmt.Sprintf("%s/users/%d/ratings/%d", baseURL, userID, movieID), nil)
    delResp := doRequest(t, delReq)
    defer delResp.Body.Close()

    assert.Equal(t, http.StatusNoContent, delResp.StatusCode,
        "удаление рейтинга должно вернуть 204")
}

// INT-24: Некорректный payload при создании рейтинга (негативный — 400)
func TestINT24_Ratings_InvalidPayload_Negative(t *testing.T) {
    req := authRequest(t, http.MethodPost,
        fmt.Sprintf("%s/users/%d/ratings", baseURL, userID),
        map[string]string{"movie_id": "bad", "rating": "bad"},
    )
    resp := doRequest(t, req)
    defer resp.Body.Close()

    assert.Equal(t, http.StatusBadRequest, resp.StatusCode,
        "некорректный payload должен вернуть 400")
}