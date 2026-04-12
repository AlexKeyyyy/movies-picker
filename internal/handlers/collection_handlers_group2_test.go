package handlers_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlexKeyyyy/movies-picker/internal/handlers"
	"github.com/AlexKeyyyy/movies-picker/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWatchlistHandler_GetWatchlist_EmptyList(t *testing.T) {
	repo := &stubRepo{
		getWatchlistFn: func(userID int64) ([]models.WatchlistItem, error) {
			assert.Equal(t, int64(1), userID)
			return []models.WatchlistItem{}, nil
		},
	}
	handler := handlers.NewWatchlistHandler(newHandlerService(repo, nil, nil))
	req := requestWithUser(httptest.NewRequest(http.MethodGet, "/users/me/watchlist", nil), 1)
	rec := httptest.NewRecorder()

	handler.GetWatchlist(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	resp := decodeBody[[]models.WatchlistItem](t, rec)
	assert.Empty(t, resp)
}

func TestWatchlistHandler_GetWatchlist_ServiceError(t *testing.T) {
	repo := &stubRepo{
		getWatchlistFn: func(int64) ([]models.WatchlistItem, error) {
			return nil, errors.New("repo failed")
		},
	}
	handler := handlers.NewWatchlistHandler(newHandlerService(repo, nil, nil))
	req := requestWithUser(httptest.NewRequest(http.MethodGet, "/users/me/watchlist", nil), 2)
	rec := httptest.NewRecorder()

	handler.GetWatchlist(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestWatchlistHandler_AddToWatchlist_InvalidJSON(t *testing.T) {
	handler := handlers.NewWatchlistHandler(newHandlerService(nil, nil, nil))
	req := requestWithUser(httptest.NewRequest(http.MethodPost, "/users/me/watchlist", bytes.NewBufferString("{")), 5)
	rec := httptest.NewRecorder()

	handler.AddToWatchlist(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestWatchlistHandler_AddToWatchlist_ServiceError(t *testing.T) {
	repo := &stubRepo{
		addToWatchlistFn: func(item *models.WatchlistItem) error {
			assert.Equal(t, int64(5), item.UserID)
			assert.Equal(t, int64(77), item.MovieID)
			return errors.New("insert failed")
		},
	}
	handler := handlers.NewWatchlistHandler(newHandlerService(repo, nil, nil))
	req := requestWithUser(httptest.NewRequest(http.MethodPost, "/users/me/watchlist", bytes.NewBufferString(`{"movie_id":77}`)), 5)
	rec := httptest.NewRecorder()

	handler.AddToWatchlist(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestWatchlistHandler_AddToWatchlist_Success(t *testing.T) {
	repo := &stubRepo{
		addToWatchlistFn: func(item *models.WatchlistItem) error {
			assert.Equal(t, int64(8), item.UserID)
			assert.Equal(t, int64(15), item.MovieID)
			return nil
		},
	}
	handler := handlers.NewWatchlistHandler(newHandlerService(repo, nil, nil))
	req := requestWithUser(httptest.NewRequest(http.MethodPost, "/users/me/watchlist", bytes.NewBufferString(`{"movie_id":15}`)), 8)
	rec := httptest.NewRecorder()

	handler.AddToWatchlist(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	resp := decodeBody[map[string]int64](t, rec)
	assert.Equal(t, int64(8), resp["user_id"])
	assert.Equal(t, int64(15), resp["movie_id"])
}

func TestWatchlistHandler_RemoveFromWatchlist_InvalidMovieID(t *testing.T) {
	handler := handlers.NewWatchlistHandler(newHandlerService(nil, nil, nil))
	req := requestWithUser(requestWithRouteParam(httptest.NewRequest(http.MethodDelete, "/users/me/watchlist/bad", nil), "movieID", "bad"), 9)
	rec := httptest.NewRecorder()

	handler.RemoveFromWatchlist(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestWatchlistHandler_RemoveFromWatchlist_ServiceError(t *testing.T) {
	repo := &stubRepo{
		removeFromWatchlistFn: func(userID, movieID int64) error {
			assert.Equal(t, int64(9), userID)
			assert.Equal(t, int64(21), movieID)
			return errors.New("delete failed")
		},
	}
	handler := handlers.NewWatchlistHandler(newHandlerService(repo, nil, nil))
	req := requestWithUser(requestWithRouteParam(httptest.NewRequest(http.MethodDelete, "/users/me/watchlist/21", nil), "movieID", "21"), 9)
	rec := httptest.NewRecorder()

	handler.RemoveFromWatchlist(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestWatchlistHandler_RemoveFromWatchlist_Success(t *testing.T) {
	repo := &stubRepo{
		removeFromWatchlistFn: func(userID, movieID int64) error {
			assert.Equal(t, int64(4), userID)
			assert.Equal(t, int64(22), movieID)
			return nil
		},
	}
	handler := handlers.NewWatchlistHandler(newHandlerService(repo, nil, nil))
	req := requestWithUser(requestWithRouteParam(httptest.NewRequest(http.MethodDelete, "/users/me/watchlist/22", nil), "movieID", "22"), 4)
	rec := httptest.NewRecorder()

	handler.RemoveFromWatchlist(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestWatchlistHandler_AddToWatchlist_ZeroMovieIDStillDelegates(t *testing.T) {
	repo := &stubRepo{
		addToWatchlistFn: func(item *models.WatchlistItem) error {
			assert.Equal(t, int64(2), item.UserID)
			assert.Equal(t, int64(0), item.MovieID)
			return nil
		},
	}
	handler := handlers.NewWatchlistHandler(newHandlerService(repo, nil, nil))
	req := requestWithUser(httptest.NewRequest(http.MethodPost, "/users/me/watchlist", bytes.NewBufferString(`{"movie_id":0}`)), 2)
	rec := httptest.NewRecorder()

	handler.AddToWatchlist(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestWatchlistHandler_RemoveFromWatchlist_ZeroMovieIDSuccess(t *testing.T) {
	repo := &stubRepo{
		removeFromWatchlistFn: func(userID, movieID int64) error {
			assert.Equal(t, int64(4), userID)
			assert.Equal(t, int64(0), movieID)
			return nil
		},
	}
	handler := handlers.NewWatchlistHandler(newHandlerService(repo, nil, nil))
	req := requestWithUser(requestWithRouteParam(httptest.NewRequest(http.MethodDelete, "/users/me/watchlist/0", nil), "movieID", "0"), 4)
	rec := httptest.NewRecorder()

	handler.RemoveFromWatchlist(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestRatingsHandler_GetRatings_Success(t *testing.T) {
	repo := &stubRepo{
		getRatingsFn: func(userID int64) ([]models.RatingItem, error) {
			assert.Equal(t, int64(3), userID)
			return []models.RatingItem{{MovieID: 10, Rating: 8}}, nil
		},
	}
	handler := handlers.NewRatingsHandler(newHandlerService(repo, nil, nil))
	req := requestWithUser(httptest.NewRequest(http.MethodGet, "/users/me/ratings", nil), 3)
	rec := httptest.NewRecorder()

	handler.GetRatings(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	resp := decodeBody[[]models.RatingItem](t, rec)
	require.Len(t, resp, 1)
	assert.Equal(t, 8, resp[0].Rating)
}

func TestRatingsHandler_GetRatings_ServiceError(t *testing.T) {
	repo := &stubRepo{
		getRatingsFn: func(int64) ([]models.RatingItem, error) {
			return nil, errors.New("repo failed")
		},
	}
	handler := handlers.NewRatingsHandler(newHandlerService(repo, nil, nil))
	req := requestWithUser(httptest.NewRequest(http.MethodGet, "/users/me/ratings", nil), 3)
	rec := httptest.NewRecorder()

	handler.GetRatings(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestRatingsHandler_GetRatings_EmptyList(t *testing.T) {
	repo := &stubRepo{
		getRatingsFn: func(userID int64) ([]models.RatingItem, error) {
			assert.Equal(t, int64(4), userID)
			return []models.RatingItem{}, nil
		},
	}
	handler := handlers.NewRatingsHandler(newHandlerService(repo, nil, nil))
	req := requestWithUser(httptest.NewRequest(http.MethodGet, "/users/me/ratings", nil), 4)
	rec := httptest.NewRecorder()

	handler.GetRatings(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	resp := decodeBody[[]models.RatingItem](t, rec)
	assert.Empty(t, resp)
}

func TestRatingsHandler_AddOrUpdateRating_InvalidJSON(t *testing.T) {
	handler := handlers.NewRatingsHandler(newHandlerService(nil, nil, nil))
	req := requestWithUser(httptest.NewRequest(http.MethodPost, "/users/me/ratings", bytes.NewBufferString("bad-json")), 1)
	rec := httptest.NewRecorder()

	handler.AddOrUpdateRating(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRatingsHandler_AddOrUpdateRating_ServiceError(t *testing.T) {
	repo := &stubRepo{
		upsertRatingFn: func(item *models.RatingItem) error {
			assert.Equal(t, int64(2), item.UserID)
			assert.Equal(t, int64(11), item.MovieID)
			assert.Equal(t, 4, item.Rating)
			return errors.New("save failed")
		},
	}
	handler := handlers.NewRatingsHandler(newHandlerService(repo, nil, nil))
	req := requestWithUser(httptest.NewRequest(http.MethodPost, "/users/me/ratings", bytes.NewBufferString(`{"movie_id":11,"rating":4}`)), 2)
	rec := httptest.NewRecorder()

	handler.AddOrUpdateRating(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestRatingsHandler_AddOrUpdateRating_Success(t *testing.T) {
	repo := &stubRepo{
		upsertRatingFn: func(item *models.RatingItem) error {
			assert.Equal(t, int64(6), item.UserID)
			assert.Equal(t, int64(99), item.MovieID)
			assert.Equal(t, 10, item.Rating)
			return nil
		},
	}
	handler := handlers.NewRatingsHandler(newHandlerService(repo, nil, nil))
	req := requestWithUser(httptest.NewRequest(http.MethodPost, "/users/me/ratings", bytes.NewBufferString(`{"movie_id":99,"rating":10}`)), 6)
	rec := httptest.NewRecorder()

	handler.AddOrUpdateRating(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	resp := decodeBody[models.RatingItem](t, rec)
	assert.Equal(t, int64(99), resp.MovieID)
	assert.Equal(t, 10, resp.Rating)
}

func TestRatingsHandler_AddOrUpdateRating_ZeroValuesPassThrough(t *testing.T) {
	repo := &stubRepo{
		upsertRatingFn: func(item *models.RatingItem) error {
			assert.Equal(t, int64(6), item.UserID)
			assert.Equal(t, int64(0), item.MovieID)
			assert.Equal(t, 0, item.Rating)
			return nil
		},
	}
	handler := handlers.NewRatingsHandler(newHandlerService(repo, nil, nil))
	req := requestWithUser(httptest.NewRequest(http.MethodPost, "/users/me/ratings", bytes.NewBufferString(`{"movie_id":0,"rating":0}`)), 6)
	rec := httptest.NewRecorder()

	handler.AddOrUpdateRating(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestRatingsHandler_DeleteRating_Success(t *testing.T) {
	repo := &stubRepo{
		deleteRatingFn: func(userID, movieID int64) error {
			assert.Equal(t, int64(7), userID)
			assert.Equal(t, int64(55), movieID)
			return nil
		},
	}
	handler := handlers.NewRatingsHandler(newHandlerService(repo, nil, nil))
	req := requestWithRouteParam(httptest.NewRequest(http.MethodDelete, "/users/7/ratings/55", nil), "userID", "7")
	req = requestWithRouteParam(req, "movieID", "55")
	rec := httptest.NewRecorder()

	handler.DeleteRating(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestRatingsHandler_DeleteRating_ServiceError(t *testing.T) {
	repo := &stubRepo{
		deleteRatingFn: func(userID, movieID int64) error {
			assert.Equal(t, int64(7), userID)
			assert.Equal(t, int64(55), movieID)
			return errors.New("delete failed")
		},
	}
	handler := handlers.NewRatingsHandler(newHandlerService(repo, nil, nil))
	req := requestWithRouteParam(httptest.NewRequest(http.MethodDelete, "/users/7/ratings/55", nil), "userID", "7")
	req = requestWithRouteParam(req, "movieID", "55")
	rec := httptest.NewRecorder()

	handler.DeleteRating(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestRatingsHandler_DeleteRating_InvalidParamsBecomeZero(t *testing.T) {
	repo := &stubRepo{
		deleteRatingFn: func(userID, movieID int64) error {
			assert.Equal(t, int64(0), userID)
			assert.Equal(t, int64(0), movieID)
			return nil
		},
	}
	handler := handlers.NewRatingsHandler(newHandlerService(repo, nil, nil))
	req := requestWithRouteParam(httptest.NewRequest(http.MethodDelete, "/users/bad/ratings/nope", nil), "userID", "bad")
	req = requestWithRouteParam(req, "movieID", "nope")
	rec := httptest.NewRecorder()

	handler.DeleteRating(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestUserHandler_GetProfile_ContentType(t *testing.T) {
	handler := handlers.NewUserHandler(&mockUserService{})
	req := requestWithUserID(http.MethodGet, "/users/me", nil, 1)
	rec := httptest.NewRecorder()

	handler.GetProfile(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
}

func TestUserHandler_UpdateProfile_ContentType(t *testing.T) {
	handler := handlers.NewUserHandler(&mockUserService{})
	req := requestWithUserID(http.MethodPatch, "/users/me", []byte(`{"email":"new@example.com"}`), 1)
	rec := httptest.NewRecorder()

	handler.UpdateProfile(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
}

func TestUserHandler_UpdateProfile_ResponseBody(t *testing.T) {
	handler := handlers.NewUserHandler(&mockUserService{})
	req := requestWithUserID(http.MethodPatch, "/users/me", []byte(`{"email":"new@example.com"}`), 1)
	rec := httptest.NewRecorder()

	handler.UpdateProfile(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	resp := decodeBody[models.User](t, rec)
	assert.Equal(t, int64(1), resp.ID)
	assert.Equal(t, "new@example.com", resp.Email)
}

func TestUserHandler_GetProfile_ResponseBody(t *testing.T) {
	handler := handlers.NewUserHandler(&mockUserService{})
	req := requestWithUserID(http.MethodGet, "/users/me", nil, 1)
	rec := httptest.NewRecorder()

	handler.GetProfile(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	resp := decodeBody[models.User](t, rec)
	assert.Equal(t, int64(1), resp.ID)
	assert.Equal(t, "user@example.com", resp.Email)
}

func TestWatchlistHandler_GetWatchlist_ResponseTitle(t *testing.T) {
	repo := &stubRepo{
		getWatchlistFn: func(userID int64) ([]models.WatchlistItem, error) {
			assert.Equal(t, int64(10), userID)
			return []models.WatchlistItem{{MovieID: 5, Title: "Arrival"}}, nil
		},
	}
	handler := handlers.NewWatchlistHandler(newHandlerService(repo, nil, nil))
	req := requestWithUser(httptest.NewRequest(http.MethodGet, "/users/me/watchlist", nil), 10)
	rec := httptest.NewRecorder()

	handler.GetWatchlist(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	resp := decodeBody[[]models.WatchlistItem](t, rec)
	require.Len(t, resp, 1)
	assert.Equal(t, "Arrival", resp[0].Title)
}
