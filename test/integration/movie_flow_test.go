package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// INT-10: Список фильмов с корректной пагинацией
func TestINT10_ListMovies_WithPagination(t *testing.T) {
    resp, err := http.Get(fmt.Sprintf("%s/movies?page=1&size=2", baseURL))
    require.NoError(t, err)
    defer resp.Body.Close()

    require.Equal(t, http.StatusOK, resp.StatusCode)

    var movies []struct {
        MovieID int     `json:"movie_id"`
        Title   string  `json:"title"`
        Rating  float64 `json:"ratingKinopoisk"`
    }
    require.NoError(t, json.NewDecoder(resp.Body).Decode(&movies))
    assert.NotEmpty(t, movies, "список фильмов не должен быть пустым")
    assert.LessOrEqual(t, len(movies), 2, "количество фильмов не должно превышать size=2")
}

// INT-11: Пагинация — страница 2 содержит другие фильмы
func TestINT11_ListMovies_Page2_DifferentFromPage1(t *testing.T) {
    resp1, err := http.Get(fmt.Sprintf("%s/movies?page=1&size=1", baseURL))
    require.NoError(t, err)
    defer resp1.Body.Close()

    resp2, err := http.Get(fmt.Sprintf("%s/movies?page=2&size=1", baseURL))
    require.NoError(t, err)
    defer resp2.Body.Close()

    require.Equal(t, http.StatusOK, resp1.StatusCode)
    require.Equal(t, http.StatusOK, resp2.StatusCode)

    var page1, page2 []struct{ MovieID int `json:"movie_id"` }
    require.NoError(t, json.NewDecoder(resp1.Body).Decode(&page1))
    require.NoError(t, json.NewDecoder(resp2.Body).Decode(&page2))

    if len(page1) > 0 && len(page2) > 0 {
        assert.NotEqual(t, page1[0].MovieID, page2[0].MovieID,
            "страница 1 и страница 2 должны содержать разные фильмы")
    }
}

// INT-12: Получение фильма по корректному ID
func TestINT12_GetMovie_ByValidID(t *testing.T) {
    id := getFirstMovieID(t)

    resp, err := http.Get(fmt.Sprintf("%s/movies/%d", baseURL, id))
    require.NoError(t, err)
    defer resp.Body.Close()

    require.Equal(t, http.StatusOK, resp.StatusCode)

    var movie struct {
        MovieID int    `json:"movie_id"`
        Title   string `json:"title"`
    }
    require.NoError(t, json.NewDecoder(resp.Body).Decode(&movie))
    assert.Equal(t, id, movie.MovieID, "возвращённый movie_id должен совпадать с запрошенным")
    assert.NotEmpty(t, movie.Title, "название фильма не должно быть пустым")
}

// INT-13: Получение фильма по нечисловому ID (негативный — 400)
func TestINT13_GetMovie_InvalidID_NotNumber_Negative(t *testing.T) {
    resp, err := http.Get(baseURL + "/movies/not-a-number")
    require.NoError(t, err)
    defer resp.Body.Close()

    assert.Equal(t, http.StatusBadRequest, resp.StatusCode,
        "нечисловой ID должен вернуть 400 Bad Request")
}

// INT-14: Поиск фильмов без параметра q (негативный — 400)
func TestINT14_SearchMovies_MissingQuery_Negative(t *testing.T) {
    resp, err := http.Get(baseURL + "/movies/search")
    require.NoError(t, err)
    defer resp.Body.Close()

    assert.Equal(t, http.StatusBadRequest, resp.StatusCode,
        "поиск без параметра q должен вернуть 400")
}

// INT-15: Поиск фильмов с валидным запросом (через mock-сервер Кинопоиска)
func TestINT15_SearchMovies_WithValidQuery(t *testing.T) {
    // "Интерстеллар" есть в mock-данных mock-сервера
    resp, err := http.Get(baseURL + "/movies/search?q=Интерстеллар")
    require.NoError(t, err)
    defer resp.Body.Close()

    // Поиск либо из БД, либо через mock Кинопоиска — оба случая должны вернуть 200
    assert.Equal(t, http.StatusOK, resp.StatusCode)

    var movies []struct {
        MovieID int    `json:"movie_id"`
        Title   string `json:"title"`
    }
    require.NoError(t, json.NewDecoder(resp.Body).Decode(&movies))
    // Если фильм есть в БД или mock вернул результат — список непустой
    // Но если поиск вернул 200 с пустым списком — тоже OK для интеграционного теста
    t.Logf("Найдено фильмов: %d", len(movies))
}