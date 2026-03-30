package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlexKeyyyy/movies-picker/internal/middleware"
	"github.com/AlexKeyyyy/movies-picker/internal/models"
	"github.com/AlexKeyyyy/movies-picker/internal/service"
	kp "github.com/AlexKeyyyy/movies-picker/pkg/kinopoisk"
	yt "github.com/AlexKeyyyy/movies-picker/pkg/youtube"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

type stubRepo struct {
	createUserFn          func(*models.User) error
	getUserByEmailFn      func(string) (*models.User, error)
	getUserByIDFn         func(int64) (*models.User, error)
	updateUserFn          func(*models.User) error
	upsertMovieFn         func(*models.Movie) error
	getMovieByIDFn        func(int64) (*models.Movie, error)
	searchMoviesFn        func(string) ([]models.Movie, error)
	listMoviesFn          func(int, int) ([]models.Movie, error)
	listPopularMoviesFn   func(int) ([]models.Movie, error)
	addToWatchlistFn      func(*models.WatchlistItem) error
	getWatchlistFn        func(int64) ([]models.WatchlistItem, error)
	removeFromWatchlistFn func(int64, int64) error
	upsertRatingFn        func(*models.RatingItem) error
	getRatingsFn          func(int64) ([]models.RatingItem, error)
	deleteRatingFn        func(int64, int64) error
}

func (s *stubRepo) CreateUser(u *models.User) error {
	if s.createUserFn != nil {
		return s.createUserFn(u)
	}
	return nil
}

func (s *stubRepo) GetUserByEmail(email string) (*models.User, error) {
	if s.getUserByEmailFn != nil {
		return s.getUserByEmailFn(email)
	}
	return nil, errors.New("not implemented")
}

func (s *stubRepo) GetUserByID(userID int64) (*models.User, error) {
	if s.getUserByIDFn != nil {
		return s.getUserByIDFn(userID)
	}
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateUser(user *models.User) error {
	if s.updateUserFn != nil {
		return s.updateUserFn(user)
	}
	return nil
}

func (s *stubRepo) UpsertMovie(movie *models.Movie) error {
	if s.upsertMovieFn != nil {
		return s.upsertMovieFn(movie)
	}
	return nil
}

func (s *stubRepo) GetMovieByID(id int64) (*models.Movie, error) {
	if s.getMovieByIDFn != nil {
		return s.getMovieByIDFn(id)
	}
	return nil, errors.New("not implemented")
}

func (s *stubRepo) SearchMovies(query string) ([]models.Movie, error) {
	if s.searchMoviesFn != nil {
		return s.searchMoviesFn(query)
	}
	return nil, nil
}

func (s *stubRepo) ListMovies(offset, limit int) ([]models.Movie, error) {
	if s.listMoviesFn != nil {
		return s.listMoviesFn(offset, limit)
	}
	return nil, nil
}

func (s *stubRepo) ListPopularMovies(limit int) ([]models.Movie, error) {
	if s.listPopularMoviesFn != nil {
		return s.listPopularMoviesFn(limit)
	}
	return nil, nil
}

func (s *stubRepo) AddToWatchlist(item *models.WatchlistItem) error {
	if s.addToWatchlistFn != nil {
		return s.addToWatchlistFn(item)
	}
	return nil
}

func (s *stubRepo) GetWatchlist(userID int64) ([]models.WatchlistItem, error) {
	if s.getWatchlistFn != nil {
		return s.getWatchlistFn(userID)
	}
	return nil, nil
}

func (s *stubRepo) RemoveFromWatchlist(userID, movieID int64) error {
	if s.removeFromWatchlistFn != nil {
		return s.removeFromWatchlistFn(userID, movieID)
	}
	return nil
}

func (s *stubRepo) UpsertRating(item *models.RatingItem) error {
	if s.upsertRatingFn != nil {
		return s.upsertRatingFn(item)
	}
	return nil
}

func (s *stubRepo) GetRatings(userID int64) ([]models.RatingItem, error) {
	if s.getRatingsFn != nil {
		return s.getRatingsFn(userID)
	}
	return nil, nil
}

func (s *stubRepo) DeleteRating(userID, movieID int64) error {
	if s.deleteRatingFn != nil {
		return s.deleteRatingFn(userID, movieID)
	}
	return nil
}

type stubKP struct {
	searchByKeywordFn func(string, int) ([]kp.Film, int, error)
	getPopularAllFn   func(int) ([]kp.Film, int, error)
}

func (s *stubKP) SearchByKeyword(keyword string, page int) ([]kp.Film, int, error) {
	if s.searchByKeywordFn != nil {
		return s.searchByKeywordFn(keyword, page)
	}
	return nil, 0, errors.New("not implemented")
}

func (s *stubKP) GetPopularAll(page int) ([]kp.Film, int, error) {
	if s.getPopularAllFn != nil {
		return s.getPopularAllFn(page)
	}
	return nil, 0, errors.New("not implemented")
}

type stubYT struct {
	searchReviewsFn func(string, int) ([]yt.ReviewResult, error)
}

func (s *stubYT) SearchReviews(keyword string, maxResultsPerPage int) ([]yt.ReviewResult, error) {
	if s.searchReviewsFn != nil {
		return s.searchReviewsFn(keyword, maxResultsPerPage)
	}
	return nil, errors.New("not implemented")
}

func newHandlerService(repo *stubRepo, kpClient *stubKP, ytClient *stubYT) *service.Service {
	if repo == nil {
		repo = &stubRepo{}
	}
	if kpClient == nil {
		kpClient = &stubKP{}
	}
	if ytClient == nil {
		ytClient = &stubYT{}
	}
	return service.NewService(repo, kpClient, ytClient, "test-secret")
}

func requestWithRouteParam(req *http.Request, key, value string) *http.Request {
	routeCtx, _ := req.Context().Value(chi.RouteCtxKey).(*chi.Context)
	if routeCtx == nil {
		routeCtx = chi.NewRouteContext()
	}
	routeCtx.URLParams.Add(key, value)
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx)
	return req.WithContext(ctx)
}

func requestWithUser(req *http.Request, userID int64) *http.Request {
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
	return req.WithContext(ctx)
}

func decodeBody[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	t.Helper()
	var v T
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&v))
	return v
}

type flexibleAuthService struct {
	registerFn func(string, string) (*models.User, error)
	loginFn    func(string, string) (string, error)
}

func (s *flexibleAuthService) Register(email, password string) (*models.User, error) {
	return s.registerFn(email, password)
}

func (s *flexibleAuthService) Login(email, password string) (string, error) {
	return s.loginFn(email, password)
}
