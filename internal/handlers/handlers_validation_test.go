package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AlexKeyyyy/movies-picker/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func withURLParam(req *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func withUserID(req *http.Request, userID int64) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
}

func TestNewHandlersConstructors(t *testing.T) {
	require.NotNil(t, NewAuthHandler(nil))
	require.NotNil(t, NewMoviesHandler(nil))
	require.NotNil(t, NewUserHandler(nil))
	require.NotNil(t, NewWatchlistHandler(nil))
	require.NotNil(t, NewRatingsHandler(nil))
}

func TestAuthRegisterInvalidJSON(t *testing.T) {
	h := NewAuthHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader("{"))
	rr := httptest.NewRecorder()
	h.Register(rr, req)
	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAuthLoginInvalidJSON(t *testing.T) {
	h := NewAuthHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader("{"))
	rr := httptest.NewRecorder()
	h.Login(rr, req)
	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestMoviesSearchMissingQuery(t *testing.T) {
	h := NewMoviesHandler(nil)
	req := httptest.NewRequest(http.MethodGet, "/movies/search", nil)
	rr := httptest.NewRecorder()
	h.SearchMovies(rr, req)
	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestMoviesGetMovieInvalidID(t *testing.T) {
	h := NewMoviesHandler(nil)
	req := withURLParam(httptest.NewRequest(http.MethodGet, "/movies/not-number", nil), "id", "not-number")
	rr := httptest.NewRecorder()
	h.GetMovie(rr, req)
	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestMoviesGetMovieReviewsInvalidID(t *testing.T) {
	h := NewMoviesHandler(nil)
	req := withURLParam(httptest.NewRequest(http.MethodGet, "/movies/not-number/reviews", nil), "id", "not-number")
	rr := httptest.NewRecorder()
	h.GetMovieReviews(rr, req)
	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUserUpdateProfileInvalidJSON(t *testing.T) {
	h := NewUserHandler(nil)
	req := withUserID(httptest.NewRequest(http.MethodPatch, "/users/me", strings.NewReader("{")), 1)
	rr := httptest.NewRecorder()
	h.UpdateProfile(rr, req)
	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestWatchlistAddInvalidJSON(t *testing.T) {
	h := NewWatchlistHandler(nil)
	req := withUserID(httptest.NewRequest(http.MethodPost, "/users/1/watchlist", strings.NewReader("{")), 1)
	rr := httptest.NewRecorder()
	h.AddToWatchlist(rr, req)
	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestWatchlistRemoveInvalidMovieID(t *testing.T) {
	h := NewWatchlistHandler(nil)
	req := withUserID(withURLParam(httptest.NewRequest(http.MethodDelete, "/users/1/watchlist/nope", nil), "movieID", "nope"), 1)
	rr := httptest.NewRecorder()
	h.RemoveFromWatchlist(rr, req)
	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestRatingsAddOrUpdateInvalidJSON(t *testing.T) {
	h := NewRatingsHandler(nil)
	req := withUserID(httptest.NewRequest(http.MethodPost, "/users/1/ratings", strings.NewReader("{")), 1)
	rr := httptest.NewRecorder()
	h.AddOrUpdateRating(rr, req)
	require.Equal(t, http.StatusBadRequest, rr.Code)
}
