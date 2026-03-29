package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AlexKeyyyy/movies-picker/internal/middleware"
	"github.com/AlexKeyyyy/movies-picker/internal/models"
	"github.com/AlexKeyyyy/movies-picker/internal/service"
	"github.com/go-chi/chi/v5"
)

func TestWatchlistHandler_GetWatchlist_Success(t *testing.T) {
	repo := &mockRepo{
		getWatchlistFn: func(userID int64) ([]models.WatchlistItem, error) {
			if userID != 10 {
				t.Fatalf("expected userID 10, got %d", userID)
			}
			return []models.WatchlistItem{{MovieID: 1, Title: "Dune"}}, nil
		},
	}
	h := NewWatchlistHandler(newTestService(repo, &mockKPClient{}, &mockYTClient{}))

	req := httptest.NewRequest(http.MethodGet, "/users/10/watchlist", nil)
	req = req.WithContext(withUserID(req.Context(), 10))
	rr := httptest.NewRecorder()

	h.GetWatchlist(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	got := decodeJSON[[]models.WatchlistItem](t, rr)
	if len(got) != 1 || got[0].Title != "Dune" {
		t.Fatalf("unexpected response: %+v", got)
	}
}

func TestWatchlistHandler_AddToWatchlist_Success(t *testing.T) {
	repo := &mockRepo{
		addToWatchlistFn: func(item *models.WatchlistItem) error {
			if item.UserID != 10 || item.MovieID != 77 {
				t.Fatalf("unexpected item: %+v", item)
			}
			return nil
		},
	}
	h := NewWatchlistHandler(newTestService(repo, &mockKPClient{}, &mockYTClient{}))

	req := httptest.NewRequest(http.MethodPost, "/users/10/watchlist", jsonBody(t, map[string]int64{
		"movie_id": 77,
	}))
	req = req.WithContext(withUserID(req.Context(), 10))
	rr := httptest.NewRecorder()

	h.AddToWatchlist(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}
	got := decodeJSON[map[string]any](t, rr)
	if got["user_id"].(float64) != 10 || got["movie_id"].(float64) != 77 {
		t.Fatalf("unexpected response: %+v", got)
	}
}

func TestWatchlistHandler_AddToWatchlist_InvalidBody(t *testing.T) {
	h := NewWatchlistHandler(newTestService(&mockRepo{}, &mockKPClient{}, &mockYTClient{}))

	req := httptest.NewRequest(http.MethodPost, "/users/10/watchlist", strings.NewReader("{"))
	req = req.WithContext(withUserID(req.Context(), 10))
	rr := httptest.NewRecorder()

	h.AddToWatchlist(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestWatchlistHandler_RemoveFromWatchlist_Success(t *testing.T) {
	repo := &mockRepo{
		removeFromWatchFn: func(userID, movieID int64) error {
			if userID != 10 || movieID != 77 {
				t.Fatalf("unexpected args: userID=%d movieID=%d", userID, movieID)
			}
			return nil
		},
	}
	h := NewWatchlistHandler(newTestService(repo, &mockKPClient{}, &mockYTClient{}))

	req := httptest.NewRequest(http.MethodDelete, "/users/10/watchlist/77", nil)
	req = req.WithContext(withUserID(req.Context(), 10))
	req = withURLParam(req, "movieID", "77")
	rr := httptest.NewRecorder()

	h.RemoveFromWatchlist(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
}

func TestWatchlistHandler_RemoveFromWatchlist_InvalidID(t *testing.T) {
	h := NewWatchlistHandler(newTestService(&mockRepo{}, &mockKPClient{}, &mockYTClient{}))

	req := httptest.NewRequest(http.MethodDelete, "/users/10/watchlist/bad", nil)
	req = req.WithContext(withUserID(req.Context(), 10))
	req = withURLParam(req, "movieID", "bad")
	rr := httptest.NewRecorder()

	h.RemoveFromWatchlist(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestNewWatchlistHandler(t *testing.T) {
	h := NewWatchlistHandler(service.NewService(&mockRepo{}, nil, nil, "secret"))
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
	if h.svc == nil {
		t.Fatal("expected non-nil service")
	}
}

func TestWatchlistHandler_GetWatchlist_ServiceError(t *testing.T) {
	repo := &mockRepo{
		getWatchlistFn: func(userID int64) ([]models.WatchlistItem, error) {
			return nil, errors.New("db failed")
		},
	}

	h := NewWatchlistHandler(service.NewService(repo, nil, nil, "secret"))

	req := httptest.NewRequest(http.MethodGet, "/users/me/watchlist", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, int64(10)))
	rr := httptest.NewRecorder()

	h.GetWatchlist(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}

func TestWatchlistHandler_RemoveFromWatchlist_ServiceError(t *testing.T) {
	repo := &mockRepo{
		removeFromWatchFn: func(userID, movieID int64) error {
			if userID != 10 || movieID != 55 {
				t.Fatalf("unexpected args: userID=%d movieID=%d", userID, movieID)
			}
			return errors.New("db failed")
		},
	}

	h := NewWatchlistHandler(service.NewService(repo, nil, nil, "secret"))

	req := httptest.NewRequest(http.MethodDelete, "/users/me/watchlist/55", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, int64(10)))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("movieID", "55")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	h.RemoveFromWatchlist(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}
