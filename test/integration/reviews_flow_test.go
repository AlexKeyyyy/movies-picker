package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// INT-25: Получение обзоров фильма (через mock YouTube-сервер)
//
// Этот тест демонстрирует изоляцию от внешнего YouTube API:
// API-сервер настроен через YOUTUBE_BASE_URL на наш mock-сервер (порт 9090).
// Реальный YouTube.com не вызывается — все ответы приходят от mock-сервера.
func TestINT25_GetMovieReviews_MockYouTube(t *testing.T) {
    movieID := getFirstMovieID(t)

    resp, err := http.Get(fmt.Sprintf("%s/movies/%d/reviews", baseURL, movieID))
    require.NoError(t, err)
    defer resp.Body.Close()

    // Mock YouTube всегда возвращает данные — поэтому ожидаем 200
    require.Equal(t, http.StatusOK, resp.StatusCode,
        "получение обзоров должно вернуть 200 (данные из mock YouTube)")

    var reviews []struct {
        VideoID      string `json:"video_id"`
        VideoURL     string `json:"video_url"`
        Title        string `json:"title"`
        ChannelTitle string `json:"channel_title"`
        ThumbnailURL string `json:"thumbnail_url"`
    }
    require.NoError(t, json.NewDecoder(resp.Body).Decode(&reviews))
    assert.NotEmpty(t, reviews,
        "mock YouTube должен вернуть хотя бы один обзор")

    // Проверяем структуру первого обзора
    if len(reviews) > 0 {
        assert.NotEmpty(t, reviews[0].VideoID)
        assert.NotEmpty(t, reviews[0].VideoURL)
        assert.NotEmpty(t, reviews[0].Title)
        t.Logf("Первый обзор: VideoID=%s, Title=%s", reviews[0].VideoID, reviews[0].Title)
    }
}

// INT-26: Обзоры для несуществующего фильма (негативный — 404/500)
func TestINT26_GetMovieReviews_InvalidMovieID_Negative(t *testing.T) {
    resp, err := http.Get(baseURL + "/movies/9999999/reviews")
    require.NoError(t, err)
    defer resp.Body.Close()

    // Несуществующий фильм: либо 404 (фильм не найден), либо 500 (ошибка сервера)
    assert.True(t,
        resp.StatusCode == http.StatusNotFound ||
            resp.StatusCode == http.StatusInternalServerError,
        "запрос обзоров несуществующего фильма должен вернуть 404 или 500, получено: %d",
        resp.StatusCode,
    )
}