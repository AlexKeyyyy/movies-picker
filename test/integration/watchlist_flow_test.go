package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// INT-16: Добавление фильма в список «Смотреть позже»
func TestINT16_Watchlist_AddMovie(t *testing.T) {
    movieID := getFirstMovieID(t)

    req := authRequest(t, http.MethodPost,
        fmt.Sprintf("%s/users/%d/watchlist", baseURL, userID),
        map[string]int{"movie_id": movieID},
    )
    resp := doRequest(t, req)
    defer resp.Body.Close()

    assert.Equal(t, http.StatusCreated, resp.StatusCode,
        "добавление в watchlist должно вернуть 201")

    var result struct {
        UserID  int64 `json:"user_id"`
        MovieID int64 `json:"movie_id"`
    }
    require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
    assert.Equal(t, userID, result.UserID)
    assert.Equal(t, int64(movieID), result.MovieID)
}

// INT-17: Проверка наличия добавленного фильма в watchlist
func TestINT17_Watchlist_ContainsAddedMovie(t *testing.T) {
    movieID := getFirstMovieID(t)

    // Убеждаемся что фильм добавлен
    addReq := authRequest(t, http.MethodPost,
        fmt.Sprintf("%s/users/%d/watchlist", baseURL, userID),
        map[string]int{"movie_id": movieID},
    )
    addResp := doRequest(t, addReq)
    addResp.Body.Close()
    // Игнорируем статус — фильм мог быть уже добавлен в INT-16

    // Получаем список
    req := authRequest(t, http.MethodGet,
        fmt.Sprintf("%s/users/%d/watchlist", baseURL, userID), nil)
    resp := doRequest(t, req)
    defer resp.Body.Close()

    require.Equal(t, http.StatusOK, resp.StatusCode)

    var items []struct {
        MovieID int `json:"movie_id"`
    }
    require.NoError(t, json.NewDecoder(resp.Body).Decode(&items))

    found := false
    for _, item := range items {
        if item.MovieID == movieID {
            found = true
            break
        }
    }
    assert.True(t, found,
        "добавленный фильм должен присутствовать в watchlist")
}

// INT-18: Удаление фильма из watchlist
func TestINT18_Watchlist_RemoveMovie(t *testing.T) {
    movieID := getFirstMovieID(t)

    // Добавляем фильм (если ещё не добавлен)
    addReq := authRequest(t, http.MethodPost,
        fmt.Sprintf("%s/users/%d/watchlist", baseURL, userID),
        map[string]int{"movie_id": movieID},
    )
    addResp := doRequest(t, addReq)
    addResp.Body.Close()

    // Удаляем
    delReq := authRequest(t, http.MethodDelete,
        fmt.Sprintf("%s/users/%d/watchlist/%d", baseURL, userID, movieID), nil)
    delResp := doRequest(t, delReq)
    defer delResp.Body.Close()

    assert.Equal(t, http.StatusNoContent, delResp.StatusCode,
        "удаление из watchlist должно вернуть 204 No Content")
}

// INT-19: Обращение к watchlist без авторизации (негативный — 401)
func TestINT19_Watchlist_NoAuth_Negative(t *testing.T) {
    req, _ := http.NewRequest(http.MethodGet,
        fmt.Sprintf("%s/users/%d/watchlist", baseURL, userID), nil)

    resp, err := http.DefaultClient.Do(req)
    require.NoError(t, err)
    defer resp.Body.Close()

    assert.Equal(t, http.StatusUnauthorized, resp.StatusCode,
        "watchlist без токена должен вернуть 401")
}

// INT-20: Добавление в watchlist с некорректным JSON (негативный — 400)
func TestINT20_Watchlist_InvalidPayload_Negative(t *testing.T) {
    // movie_id должен быть числом, передаём строку
    req := authRequest(t, http.MethodPost,
        fmt.Sprintf("%s/users/%d/watchlist", baseURL, userID),
        map[string]string{"movie_id": "not-a-number"},
    )
    resp := doRequest(t, req)
    defer resp.Body.Close()

    assert.Equal(t, http.StatusBadRequest, resp.StatusCode,
        "некорректный JSON должен вернуть 400")
}