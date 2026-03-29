package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlexKeyyyy/movies-picker/internal/models"
	"github.com/AlexKeyyyy/movies-picker/internal/service"
	kp "github.com/AlexKeyyyy/movies-picker/pkg/kinopoisk"
	yt "github.com/AlexKeyyyy/movies-picker/pkg/youtube"
	"github.com/go-chi/chi/v5"
)

func TestMoviesHandler_SearchMovies_Success(t *testing.T) {
	repo := &mockRepo{
		searchMoviesFn: func(query string) ([]models.Movie, error) {
			if query != "matrix" {
				t.Fatalf("expected query matrix, got %s", query)
			}
			return []models.Movie{{ID: 1, Title: "The Matrix", Year: 1999}}, nil
		},
	}
	h := NewMoviesHandler(newTestService(repo, &mockKPClient{}, &mockYTClient{}))

	req := httptest.NewRequest(http.MethodGet, "/movies/search?q=matrix", nil)
	rr := httptest.NewRecorder()

	h.SearchMovies(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	got := decodeJSON[[]models.Movie](t, rr)
	if len(got) != 1 || got[0].Title != "The Matrix" {
		t.Fatalf("unexpected response: %+v", got)
	}
}

func TestMoviesHandler_SearchMovies_ServiceError(t *testing.T) {
	repo := &mockRepo{
		searchMoviesFn: func(query string) ([]models.Movie, error) {
			return nil, nil
		},
	}
	kpClient := &mockKPClient{
		searchByKeywordFn: func(keyword string, page int) ([]kp.Film, int, error) {
			return nil, 0, errors.New("kp failed")
		},
	}
	h := NewMoviesHandler(newTestService(repo, kpClient, &mockYTClient{}))

	req := httptest.NewRequest(http.MethodGet, "/movies/search?q=matrix", nil)
	rr := httptest.NewRecorder()

	h.SearchMovies(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}

func TestMoviesHandler_GetMovie_Success(t *testing.T) {
	repo := &mockRepo{
		getMovieByIDFn: func(id int64) (*models.Movie, error) {
			if id != 10 {
				t.Fatalf("expected id 10, got %d", id)
			}
			return &models.Movie{ID: 10, Title: "Interstellar", Year: 2014}, nil
		},
	}
	h := NewMoviesHandler(newTestService(repo, &mockKPClient{}, &mockYTClient{}))

	req := httptest.NewRequest(http.MethodGet, "/movies/10", nil)
	req = withURLParam(req, "id", "10")
	rr := httptest.NewRecorder()

	h.GetMovie(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	got := decodeJSON[models.Movie](t, rr)
	if got.ID != 10 || got.Title != "Interstellar" {
		t.Fatalf("unexpected response: %+v", got)
	}
}

func TestMoviesHandler_GetMovie_InvalidID(t *testing.T) {
	h := NewMoviesHandler(newTestService(&mockRepo{}, &mockKPClient{}, &mockYTClient{}))

	req := httptest.NewRequest(http.MethodGet, "/movies/bad", nil)
	req = withURLParam(req, "id", "bad")
	rr := httptest.NewRecorder()

	h.GetMovie(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestMoviesHandler_GetMovie_NotFound(t *testing.T) {
	repo := &mockRepo{
		getMovieByIDFn: func(id int64) (*models.Movie, error) {
			return nil, errors.New("not found")
		},
	}
	h := NewMoviesHandler(newTestService(repo, &mockKPClient{}, &mockYTClient{}))

	req := httptest.NewRequest(http.MethodGet, "/movies/10", nil)
	req = withURLParam(req, "id", "10")
	rr := httptest.NewRecorder()

	h.GetMovie(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestMoviesHandler_GetMovieReviews_Success(t *testing.T) {
	repo := &mockRepo{
		getMovieByIDFn: func(id int64) (*models.Movie, error) {
			return &models.Movie{ID: id, Title: "Dune"}, nil
		},
	}
	ytClient := &mockYTClient{
		searchReviewsFn: func(keyword string, maxResultsPerPage int) ([]yt.ReviewResult, error) {
			if keyword != "Dune" {
				t.Fatalf("expected keyword Dune, got %s", keyword)
			}
			return []yt.ReviewResult{{
				VideoID:      "abc",
				VideoURL:     "https://youtube.com/watch?v=abc",
				Title:        "Dune review",
				ChannelTitle: "Channel",
				ThumbnailURL: "https://img/abc.jpg",
			}}, nil
		},
	}
	h := NewMoviesHandler(newTestService(repo, &mockKPClient{}, ytClient))

	req := httptest.NewRequest(http.MethodGet, "/movies/1/reviews", nil)
	req = withURLParam(req, "id", "1")
	rr := httptest.NewRecorder()

	h.GetMovieReviews(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	got := decodeJSON[[]models.ReviewItem](t, rr)
	if len(got) != 1 || got[0].VideoID != "abc" {
		t.Fatalf("unexpected response: %+v", got)
	}
}

func TestMoviesHandler_GetMovieReviews_ServiceError(t *testing.T) {
	repo := &mockRepo{
		getMovieByIDFn: func(id int64) (*models.Movie, error) {
			return &models.Movie{ID: id, Title: "Dune"}, nil
		},
	}
	ytClient := &mockYTClient{
		searchReviewsFn: func(keyword string, maxResultsPerPage int) ([]yt.ReviewResult, error) {
			return nil, errors.New("youtube failed")
		},
	}
	h := NewMoviesHandler(newTestService(repo, &mockKPClient{}, ytClient))

	req := httptest.NewRequest(http.MethodGet, "/movies/1/reviews", nil)
	req = withURLParam(req, "id", "1")
	rr := httptest.NewRecorder()

	h.GetMovieReviews(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}

func TestMoviesHandler_ListMovies_Success(t *testing.T) {
	repo := &mockRepo{
		listMoviesFn: func(offset, limit int) ([]models.Movie, error) {
			if offset != 3 {
				t.Fatalf("expected offset 3, got %d", offset)
			}
			if limit != 3 {
				t.Fatalf("expected limit 3, got %d", limit)
			}
			return []models.Movie{{ID: 1, Title: "Movie 1"}}, nil
		},
	}
	h := NewMoviesHandler(newTestService(repo, &mockKPClient{}, &mockYTClient{}))

	req := httptest.NewRequest(http.MethodGet, "/movies?page=2&size=3", nil)
	rr := httptest.NewRecorder()

	h.ListMovies(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	got := decodeJSON[[]models.Movie](t, rr)
	if len(got) != 1 || got[0].Title != "Movie 1" {
		t.Fatalf("unexpected response: %+v", got)
	}
}

func TestMoviesHandler_ListPopular_Success(t *testing.T) {
	repo := &mockRepo{
		listPopularMoviesFn: func(limit int) ([]models.Movie, error) {
			if limit != 5 {
				t.Fatalf("expected limit 5, got %d", limit)
			}
			return []models.Movie{{ID: 99, Title: "Popular Movie"}}, nil
		},
	}
	h := NewMoviesHandler(newTestService(repo, &mockKPClient{}, &mockYTClient{}))

	req := httptest.NewRequest(http.MethodGet, "/movies/popular?limit=5", nil)
	rr := httptest.NewRecorder()

	h.ListPopular(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	got := decodeJSON[[]models.Movie](t, rr)
	if len(got) != 1 || got[0].ID != 99 {
		t.Fatalf("unexpected response: %+v", got)
	}
}

func withURLParam(req *http.Request, key, value string) *http.Request {
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add(key, value)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
}

func TestNewMoviesHandler(t *testing.T) {
	h := NewMoviesHandler(service.NewService(&mockRepo{}, nil, nil, "secret"))
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
	if h.svc == nil {
		t.Fatal("expected non-nil service")
	}
}

func TestMoviesHandler_SearchMovies_MissingQuery(t *testing.T) {
	h := NewMoviesHandler(service.NewService(&mockRepo{}, nil, nil, "secret"))

	req := httptest.NewRequest(http.MethodGet, "/movies/search", nil)
	rr := httptest.NewRecorder()

	h.SearchMovies(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestMoviesHandler_GetMovieReviews_InvalidID(t *testing.T) {
	h := NewMoviesHandler(service.NewService(&mockRepo{}, nil, nil, "secret"))

	req := httptest.NewRequest(http.MethodGet, "/movies/abc/reviews", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "abc")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	h.GetMovieReviews(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestMoviesHandler_ListMovies_DefaultPagination(t *testing.T) {
	repo := &mockRepo{
		listMoviesFn: func(offset, limit int) ([]models.Movie, error) {
			if offset != 0 || limit != 20 {
				t.Fatalf("expected offset=0 limit=20, got offset=%d limit=%d", offset, limit)
			}
			return []models.Movie{}, nil
		},
	}

	h := NewMoviesHandler(service.NewService(repo, nil, nil, "secret"))

	req := httptest.NewRequest(http.MethodGet, "/movies", nil)
	rr := httptest.NewRecorder()

	h.ListMovies(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestMoviesHandler_ListMovies_InvalidPageDefaults(t *testing.T) {
	repo := &mockRepo{
		listMoviesFn: func(offset, limit int) ([]models.Movie, error) {
			if offset != 0 || limit != 5 {
				t.Fatalf("expected offset=0 limit=5, got offset=%d limit=%d", offset, limit)
			}
			return []models.Movie{}, nil
		},
	}

	h := NewMoviesHandler(service.NewService(repo, nil, nil, "secret"))

	req := httptest.NewRequest(http.MethodGet, "/movies?page=abc&size=5", nil)
	rr := httptest.NewRecorder()

	h.ListMovies(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestMoviesHandler_ListMovies_InvalidSizeDefaults(t *testing.T) {
	repo := &mockRepo{
		listMoviesFn: func(offset, limit int) ([]models.Movie, error) {
			if offset != 20 || limit != 20 {
				t.Fatalf("expected offset=20 limit=20, got offset=%d limit=%d", offset, limit)
			}
			return []models.Movie{}, nil
		},
	}

	h := NewMoviesHandler(service.NewService(repo, nil, nil, "secret"))

	req := httptest.NewRequest(http.MethodGet, "/movies?page=2&size=abc", nil)
	rr := httptest.NewRecorder()

	h.ListMovies(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestMoviesHandler_ListMovies_ServiceError(t *testing.T) {
	repo := &mockRepo{
		listMoviesFn: func(offset, limit int) ([]models.Movie, error) {
			return nil, errors.New("db failed")
		},
	}

	h := NewMoviesHandler(service.NewService(repo, nil, nil, "secret"))

	req := httptest.NewRequest(http.MethodGet, "/movies?page=1&size=10", nil)
	rr := httptest.NewRecorder()

	h.ListMovies(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}

func TestMoviesHandler_ListPopular_DefaultLimit(t *testing.T) {
	repo := &mockRepo{
		listPopularMoviesFn: func(limit int) ([]models.Movie, error) {
			if limit != 10 {
				t.Fatalf("expected limit=10, got %d", limit)
			}
			return []models.Movie{}, nil
		},
	}

	h := NewMoviesHandler(service.NewService(repo, nil, nil, "secret"))

	req := httptest.NewRequest(http.MethodGet, "/movies/popular", nil)
	rr := httptest.NewRecorder()

	h.ListPopular(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestMoviesHandler_ListPopular_InvalidLimitDefaults(t *testing.T) {
	repo := &mockRepo{
		listPopularMoviesFn: func(limit int) ([]models.Movie, error) {
			if limit != 10 {
				t.Fatalf("expected limit=10, got %d", limit)
			}
			return []models.Movie{}, nil
		},
	}

	h := NewMoviesHandler(service.NewService(repo, nil, nil, "secret"))

	req := httptest.NewRequest(http.MethodGet, "/movies/popular?limit=abc", nil)
	rr := httptest.NewRecorder()

	h.ListPopular(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestMoviesHandler_ListPopular_ServiceError(t *testing.T) {
	repo := &mockRepo{
		listPopularMoviesFn: func(limit int) ([]models.Movie, error) {
			return nil, errors.New("db failed")
		},
	}

	h := NewMoviesHandler(service.NewService(repo, nil, nil, "secret"))

	req := httptest.NewRequest(http.MethodGet, "/movies/popular?limit=5", nil)
	rr := httptest.NewRecorder()

	h.ListPopular(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}
