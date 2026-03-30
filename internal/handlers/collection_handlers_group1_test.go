package handlers_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlexKeyyyy/movies-picker/internal/handlers"
	"github.com/AlexKeyyyy/movies-picker/internal/models"
	kp "github.com/AlexKeyyyy/movies-picker/pkg/kinopoisk"
	yt "github.com/AlexKeyyyy/movies-picker/pkg/youtube"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthHandler_Register_InternalError(t *testing.T) {
	handler := handlers.NewAuthHandler(&flexibleAuthService{
		registerFn: func(string, string) (*models.User, error) {
			return nil, errors.New("unexpected")
		},
		loginFn: func(string, string) (string, error) {
			return "", nil
		},
	})
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(`{"email":"a@a.com","password":"pass"}`))
	rec := httptest.NewRecorder()

	handler.Register(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestAuthHandler_Register_ResponseContainsCreated(t *testing.T) {
	handler := handlers.NewAuthHandler(&flexibleAuthService{
		registerFn: func(email, password string) (*models.User, error) {
			return &models.User{ID: 2, Email: email, CreatedAt: "2026-03-30T10:00:00Z"}, nil
		},
		loginFn: func(string, string) (string, error) {
			return "", nil
		},
	})
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(`{"email":"a@a.com","password":"pass"}`))
	rec := httptest.NewRecorder()

	handler.Register(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	resp := decodeBody[map[string]any](t, rec)
	assert.Equal(t, "2026-03-30T10:00:00Z", resp["created"])
}

func TestAuthHandler_Login_ResponseMetadata(t *testing.T) {
	handler := handlers.NewAuthHandler(&flexibleAuthService{
		registerFn: func(string, string) (*models.User, error) {
			return nil, nil
		},
		loginFn: func(email, password string) (string, error) {
			assert.Equal(t, "demo@example.com", email)
			assert.Equal(t, "pass123", password)
			return "jwt-token", nil
		},
	})
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(`{"email":"demo@example.com","password":"pass123"}`))
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	resp := decodeBody[map[string]any](t, rec)
	assert.Equal(t, "jwt-token", resp["access_token"])
	assert.Equal(t, "Bearer", resp["token_type"])
	assert.Equal(t, float64(3600), resp["expires_in"])
}

func TestMoviesHandler_SearchMovies_MissingQuery(t *testing.T) {
	handler := handlers.NewMoviesHandler(newHandlerService(nil, nil, nil))
	req := httptest.NewRequest(http.MethodGet, "/movies/search", nil)
	rec := httptest.NewRecorder()

	handler.SearchMovies(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestMoviesHandler_SearchMovies_Success(t *testing.T) {
	repo := &stubRepo{
		searchMoviesFn: func(query string) ([]models.Movie, error) {
			assert.Equal(t, "matrix", query)
			return []models.Movie{{ID: 1, Title: "Matrix"}}, nil
		},
	}
	handler := handlers.NewMoviesHandler(newHandlerService(repo, nil, nil))
	req := httptest.NewRequest(http.MethodGet, "/movies/search?q=matrix", nil)
	rec := httptest.NewRecorder()

	handler.SearchMovies(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	resp := decodeBody[[]models.Movie](t, rec)
	require.Len(t, resp, 1)
	assert.Equal(t, int64(1), resp[0].ID)
}

func TestMoviesHandler_SearchMovies_ServiceError(t *testing.T) {
	repo := &stubRepo{
		searchMoviesFn: func(string) ([]models.Movie, error) {
			return []models.Movie{}, nil
		},
	}
	kpClient := &stubKP{
		searchByKeywordFn: func(keyword string, page int) ([]kp.Film, int, error) {
			assert.Equal(t, "matrix", keyword)
			assert.Equal(t, 1, page)
			return nil, 0, errors.New("kp unavailable")
		},
	}
	handler := handlers.NewMoviesHandler(newHandlerService(repo, kpClient, nil))
	req := httptest.NewRequest(http.MethodGet, "/movies/search?q=matrix", nil)
	rec := httptest.NewRecorder()

	handler.SearchMovies(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestMoviesHandler_GetMovie_InvalidID(t *testing.T) {
	handler := handlers.NewMoviesHandler(newHandlerService(nil, nil, nil))
	req := requestWithRouteParam(httptest.NewRequest(http.MethodGet, "/movies/bad", nil), "id", "bad")
	rec := httptest.NewRecorder()

	handler.GetMovie(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestMoviesHandler_GetMovie_NotFound(t *testing.T) {
	repo := &stubRepo{
		getMovieByIDFn: func(id int64) (*models.Movie, error) {
			assert.Equal(t, int64(77), id)
			return nil, errors.New("not found")
		},
	}
	handler := handlers.NewMoviesHandler(newHandlerService(repo, nil, nil))
	req := requestWithRouteParam(httptest.NewRequest(http.MethodGet, "/movies/77", nil), "id", "77")
	rec := httptest.NewRecorder()

	handler.GetMovie(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestMoviesHandler_GetMovie_Success(t *testing.T) {
	repo := &stubRepo{
		getMovieByIDFn: func(id int64) (*models.Movie, error) {
			return &models.Movie{ID: id, Title: "Interstellar"}, nil
		},
	}
	handler := handlers.NewMoviesHandler(newHandlerService(repo, nil, nil))
	req := requestWithRouteParam(httptest.NewRequest(http.MethodGet, "/movies/5", nil), "id", "5")
	rec := httptest.NewRecorder()

	handler.GetMovie(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	resp := decodeBody[models.Movie](t, rec)
	assert.Equal(t, int64(5), resp.ID)
	assert.Equal(t, "Interstellar", resp.Title)
}

func TestMoviesHandler_GetMovieReviews_InvalidID(t *testing.T) {
	handler := handlers.NewMoviesHandler(newHandlerService(nil, nil, nil))
	req := requestWithRouteParam(httptest.NewRequest(http.MethodGet, "/movies/nope/reviews", nil), "id", "nope")
	rec := httptest.NewRecorder()

	handler.GetMovieReviews(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestMoviesHandler_GetMovieReviews_ServiceError(t *testing.T) {
	repo := &stubRepo{
		getMovieByIDFn: func(id int64) (*models.Movie, error) {
			return &models.Movie{ID: id, Title: "Dune"}, nil
		},
	}
	ytClient := &stubYT{
		searchReviewsFn: func(keyword string, max int) ([]yt.ReviewResult, error) {
			assert.Equal(t, "Dune", keyword)
			assert.Equal(t, 10, max)
			return nil, errors.New("yt failed")
		},
	}
	handler := handlers.NewMoviesHandler(newHandlerService(repo, nil, ytClient))
	req := requestWithRouteParam(httptest.NewRequest(http.MethodGet, "/movies/9/reviews", nil), "id", "9")
	rec := httptest.NewRecorder()

	handler.GetMovieReviews(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestMoviesHandler_GetMovieReviews_Success(t *testing.T) {
	repo := &stubRepo{
		getMovieByIDFn: func(id int64) (*models.Movie, error) {
			return &models.Movie{ID: id, Title: "Dune"}, nil
		},
	}
	ytClient := &stubYT{
		searchReviewsFn: func(string, int) ([]yt.ReviewResult, error) {
			return []yt.ReviewResult{{VideoID: "v1", Title: "Review"}}, nil
		},
	}
	handler := handlers.NewMoviesHandler(newHandlerService(repo, nil, ytClient))
	req := requestWithRouteParam(httptest.NewRequest(http.MethodGet, "/movies/9/reviews", nil), "id", "9")
	rec := httptest.NewRecorder()

	handler.GetMovieReviews(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	resp := decodeBody[[]models.ReviewItem](t, rec)
	require.Len(t, resp, 1)
	assert.Equal(t, "v1", resp[0].VideoID)
}

func TestMoviesHandler_ListMovies_DefaultsInvalidPagination(t *testing.T) {
	repo := &stubRepo{
		listMoviesFn: func(offset, limit int) ([]models.Movie, error) {
			assert.Equal(t, 0, offset)
			assert.Equal(t, 20, limit)
			return []models.Movie{}, nil
		},
	}
	handler := handlers.NewMoviesHandler(newHandlerService(repo, nil, nil))
	req := httptest.NewRequest(http.MethodGet, "/movies?page=bad&size=-1", nil)
	rec := httptest.NewRecorder()

	handler.ListMovies(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
}

func TestMoviesHandler_ListMovies_CustomPagination(t *testing.T) {
	repo := &stubRepo{
		listMoviesFn: func(offset, limit int) ([]models.Movie, error) {
			assert.Equal(t, 10, offset)
			assert.Equal(t, 10, limit)
			return []models.Movie{{ID: 2, Title: "Blade Runner"}}, nil
		},
	}
	handler := handlers.NewMoviesHandler(newHandlerService(repo, nil, nil))
	req := httptest.NewRequest(http.MethodGet, "/movies?page=2&size=10", nil)
	rec := httptest.NewRecorder()

	handler.ListMovies(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	resp := decodeBody[[]models.Movie](t, rec)
	require.Len(t, resp, 1)
	assert.Equal(t, int64(2), resp[0].ID)
}

func TestMoviesHandler_ListMovies_ServiceError(t *testing.T) {
	repo := &stubRepo{
		listMoviesFn: func(int, int) ([]models.Movie, error) {
			return nil, errors.New("db failed")
		},
	}
	handler := handlers.NewMoviesHandler(newHandlerService(repo, nil, nil))
	req := httptest.NewRequest(http.MethodGet, "/movies?page=1&size=20", nil)
	rec := httptest.NewRecorder()

	handler.ListMovies(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestMoviesHandler_ListMovies_ZeroPageDefaultsToFirstPage(t *testing.T) {
	repo := &stubRepo{
		listMoviesFn: func(offset, limit int) ([]models.Movie, error) {
			assert.Equal(t, 0, offset)
			assert.Equal(t, 5, limit)
			return []models.Movie{}, nil
		},
	}
	handler := handlers.NewMoviesHandler(newHandlerService(repo, nil, nil))
	req := httptest.NewRequest(http.MethodGet, "/movies?page=0&size=5", nil)
	rec := httptest.NewRecorder()

	handler.ListMovies(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestMoviesHandler_ListMovies_ZeroSizeDefaultsToTwenty(t *testing.T) {
	repo := &stubRepo{
		listMoviesFn: func(offset, limit int) ([]models.Movie, error) {
			assert.Equal(t, 0, offset)
			assert.Equal(t, 20, limit)
			return []models.Movie{}, nil
		},
	}
	handler := handlers.NewMoviesHandler(newHandlerService(repo, nil, nil))
	req := httptest.NewRequest(http.MethodGet, "/movies?page=1&size=0", nil)
	rec := httptest.NewRecorder()

	handler.ListMovies(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestMoviesHandler_ListPopular_DefaultLimitOnMissingValue(t *testing.T) {
	repo := &stubRepo{
		listPopularMoviesFn: func(limit int) ([]models.Movie, error) {
			assert.Equal(t, 10, limit)
			return []models.Movie{}, nil
		},
	}
	handler := handlers.NewMoviesHandler(newHandlerService(repo, nil, nil))
	req := httptest.NewRequest(http.MethodGet, "/movies/popular", nil)
	rec := httptest.NewRecorder()

	handler.ListPopular(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestMoviesHandler_ListPopular_DefaultLimitOnNegativeValue(t *testing.T) {
	repo := &stubRepo{
		listPopularMoviesFn: func(limit int) ([]models.Movie, error) {
			assert.Equal(t, 10, limit)
			return []models.Movie{}, nil
		},
	}
	handler := handlers.NewMoviesHandler(newHandlerService(repo, nil, nil))
	req := httptest.NewRequest(http.MethodGet, "/movies/popular?limit=-3", nil)
	rec := httptest.NewRecorder()

	handler.ListPopular(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestMoviesHandler_ListPopular_Success(t *testing.T) {
	repo := &stubRepo{
		listPopularMoviesFn: func(limit int) ([]models.Movie, error) {
			assert.Equal(t, 3, limit)
			return []models.Movie{{ID: 10, Title: "Tenet"}}, nil
		},
	}
	handler := handlers.NewMoviesHandler(newHandlerService(repo, nil, nil))
	req := httptest.NewRequest(http.MethodGet, "/movies/popular?limit=3", nil)
	rec := httptest.NewRecorder()

	handler.ListPopular(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	resp := decodeBody[[]models.Movie](t, rec)
	require.Len(t, resp, 1)
	assert.Equal(t, int64(10), resp[0].ID)
}

func TestMoviesHandler_ListPopular_ServiceError(t *testing.T) {
	repo := &stubRepo{
		listPopularMoviesFn: func(int) ([]models.Movie, error) {
			return nil, errors.New("repo failed")
		},
	}
	handler := handlers.NewMoviesHandler(newHandlerService(repo, nil, nil))
	req := httptest.NewRequest(http.MethodGet, "/movies/popular?limit=2", nil)
	rec := httptest.NewRecorder()

	handler.ListPopular(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestMoviesHandler_ListPopular_InvalidTextLimitDefaults(t *testing.T) {
	repo := &stubRepo{
		listPopularMoviesFn: func(limit int) ([]models.Movie, error) {
			assert.Equal(t, 10, limit)
			return []models.Movie{}, nil
		},
	}
	handler := handlers.NewMoviesHandler(newHandlerService(repo, nil, nil))
	req := httptest.NewRequest(http.MethodGet, "/movies/popular?limit=abc", nil)
	rec := httptest.NewRecorder()

	handler.ListPopular(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestMoviesHandler_SearchMovies_EmptyResultStillReturnsJSON(t *testing.T) {
	repo := &stubRepo{
		searchMoviesFn: func(query string) ([]models.Movie, error) {
			assert.Equal(t, "unknown", query)
			return []models.Movie{}, nil
		},
	}
	kpClient := &stubKP{
		searchByKeywordFn: func(string, int) ([]kp.Film, int, error) {
			return []kp.Film{}, 1, nil
		},
	}
	handler := handlers.NewMoviesHandler(newHandlerService(repo, kpClient, nil))
	req := httptest.NewRequest(http.MethodGet, "/movies/search?q=unknown", nil)
	rec := httptest.NewRecorder()

	handler.SearchMovies(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	resp := decodeBody[[]models.Movie](t, rec)
	assert.Empty(t, resp)
}

func TestMoviesHandler_GetMovieReviews_EmptyReviewsSuccess(t *testing.T) {
	repo := &stubRepo{
		getMovieByIDFn: func(id int64) (*models.Movie, error) {
			return &models.Movie{ID: id, Title: "Solaris"}, nil
		},
	}
	ytClient := &stubYT{
		searchReviewsFn: func(string, int) ([]yt.ReviewResult, error) {
			return []yt.ReviewResult{}, nil
		},
	}
	handler := handlers.NewMoviesHandler(newHandlerService(repo, nil, ytClient))
	req := requestWithRouteParam(httptest.NewRequest(http.MethodGet, "/movies/13/reviews", nil), "id", "13")
	rec := httptest.NewRecorder()

	handler.GetMovieReviews(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	resp := decodeBody[[]models.ReviewItem](t, rec)
	assert.Empty(t, resp)
}

func TestWatchlistHandler_GetWatchlist_Success(t *testing.T) {
	repo := &stubRepo{
		getWatchlistFn: func(userID int64) ([]models.WatchlistItem, error) {
			assert.Equal(t, int64(12), userID)
			return []models.WatchlistItem{{MovieID: 50, Title: "Alien"}}, nil
		},
	}
	handler := handlers.NewWatchlistHandler(newHandlerService(repo, nil, nil))
	req := requestWithUser(httptest.NewRequest(http.MethodGet, "/users/me/watchlist", nil), 12)
	rec := httptest.NewRecorder()

	handler.GetWatchlist(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	resp := decodeBody[[]models.WatchlistItem](t, rec)
	require.Len(t, resp, 1)
	assert.Equal(t, int64(50), resp[0].MovieID)
}
