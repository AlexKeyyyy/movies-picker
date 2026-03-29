package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AlexKeyyyy/movies-picker/internal/models"
	"github.com/AlexKeyyyy/movies-picker/internal/service"
	"github.com/go-chi/chi/v5"
)

func TestRatingsHandler_GetRatings_Success(t *testing.T) {
	repo := &mockRepo{
		getRatingsFn: func(userID int64) ([]models.RatingItem, error) {
			if userID != 5 {
				t.Fatalf("expected userID 5, got %d", userID)
			}
			return []models.RatingItem{{MovieID: 10, Rating: 8}}, nil
		},
	}
	h := NewRatingsHandler(newTestService(repo, &mockKPClient{}, &mockYTClient{}))

	req := httptest.NewRequest(http.MethodGet, "/users/5/ratings", nil)
	req = req.WithContext(withUserID(req.Context(), 5))
	rr := httptest.NewRecorder()

	h.GetRatings(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	got := decodeJSON[[]models.RatingItem](t, rr)
	if len(got) != 1 || got[0].Rating != 8 {
		t.Fatalf("unexpected response: %+v", got)
	}
}

func TestRatingsHandler_AddOrUpdateRating_Success(t *testing.T) {
	repo := &mockRepo{
		upsertRatingFn: func(item *models.RatingItem) error {
			if item.UserID != 5 || item.MovieID != 10 || item.Rating != 9 {
				t.Fatalf("unexpected item: %+v", item)
			}
			return nil
		},
	}
	h := NewRatingsHandler(newTestService(repo, &mockKPClient{}, &mockYTClient{}))

	req := httptest.NewRequest(http.MethodPost, "/users/5/ratings", jsonBody(t, map[string]any{
		"movie_id": 10,
		"rating":   9,
	}))
	req = req.WithContext(withUserID(req.Context(), 5))
	rr := httptest.NewRecorder()

	h.AddOrUpdateRating(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}
	got := decodeJSON[models.RatingItem](t, rr)
	if got.MovieID != 10 || got.Rating != 9 {
		t.Fatalf("unexpected response: %+v", got)
	}
}

func TestRatingsHandler_AddOrUpdateRating_InvalidBody(t *testing.T) {
	h := NewRatingsHandler(newTestService(&mockRepo{}, &mockKPClient{}, &mockYTClient{}))

	req := httptest.NewRequest(http.MethodPost, "/users/5/ratings", strings.NewReader("{"))
	req = req.WithContext(withUserID(req.Context(), 5))
	rr := httptest.NewRecorder()

	h.AddOrUpdateRating(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestRatingsHandler_DeleteRating_Success(t *testing.T) {
	repo := &mockRepo{
		deleteRatingFn: func(userID, movieID int64) error {
			if userID != 7 || movieID != 10 {
				t.Fatalf("unexpected args: userID=%d movieID=%d", userID, movieID)
			}
			return nil
		},
	}

	h := NewRatingsHandler(service.NewService(repo, nil, nil, "secret"))

	req := httptest.NewRequest(http.MethodDelete, "/users/7/ratings/10", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userID", "7")
	rctx.URLParams.Add("movieID", "10")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	h.DeleteRating(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
}

func TestRatingsHandler_DeleteRating_InvalidID(t *testing.T) {
	repo := &mockRepo{
		deleteRatingFn: func(userID, movieID int64) error {
			if userID != 0 || movieID != 0 {
				t.Fatalf("unexpected args: userID=%d movieID=%d", userID, movieID)
			}
			return nil
		},
	}

	h := NewRatingsHandler(service.NewService(repo, nil, nil, "secret"))

	req := httptest.NewRequest(http.MethodDelete, "/users/abc/ratings/xyz", nil)
	req = withURLParam(req, "userID", "abc")
	req = withURLParam(req, "movieID", "xyz")

	rr := httptest.NewRecorder()
	h.DeleteRating(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
}

func TestNewRatingsHandler(t *testing.T) {
	h := NewRatingsHandler(service.NewService(&mockRepo{}, nil, nil, "secret"))
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
	if h.svc == nil {
		t.Fatal("expected non-nil service")
	}
}
