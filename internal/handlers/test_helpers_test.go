package handlers

import (
    "bytes"
    "context"
    "encoding/json"
    "errors"
    "net/http/httptest"
    "testing"

    "github.com/AlexKeyyyy/movies-picker/internal/middleware"
    "github.com/AlexKeyyyy/movies-picker/internal/models"
    "github.com/AlexKeyyyy/movies-picker/internal/service"
    kp "github.com/AlexKeyyyy/movies-picker/pkg/kinopoisk"
    yt "github.com/AlexKeyyyy/movies-picker/pkg/youtube"
)

type mockRepo struct {
    createUserFn        func(u *models.User) error
    getUserByEmailFn    func(email string) (*models.User, error)
    getUserByIDFn       func(userID int64) (*models.User, error)
    updateUserFn        func(user *models.User) error
    upsertMovieFn       func(m *models.Movie) error
    getMovieByIDFn      func(id int64) (*models.Movie, error)
    searchMoviesFn      func(query string) ([]models.Movie, error)
    listMoviesFn        func(offset, limit int) ([]models.Movie, error)
    listPopularMoviesFn func(limit int) ([]models.Movie, error)
    addToWatchlistFn    func(item *models.WatchlistItem) error
    getWatchlistFn      func(userID int64) ([]models.WatchlistItem, error)
    removeFromWatchFn   func(userID, movieID int64) error
    upsertRatingFn      func(item *models.RatingItem) error
    getRatingsFn        func(userID int64) ([]models.RatingItem, error)
    deleteRatingFn      func(userID, movieID int64) error
}

func (m *mockRepo) CreateUser(u *models.User) error {
    if m.createUserFn != nil {
        return m.createUserFn(u)
    }
    return nil
}
func (m *mockRepo) GetUserByEmail(email string) (*models.User, error) {
    if m.getUserByEmailFn != nil {
        return m.getUserByEmailFn(email)
    }
    return nil, errors.New("not implemented")
}
func (m *mockRepo) GetUserByID(userID int64) (*models.User, error) {
    if m.getUserByIDFn != nil {
        return m.getUserByIDFn(userID)
    }
    return nil, errors.New("not implemented")
}
func (m *mockRepo) UpdateUser(user *models.User) error {
    if m.updateUserFn != nil {
        return m.updateUserFn(user)
    }
    return nil
}
func (m *mockRepo) UpsertMovie(movie *models.Movie) error {
    if m.upsertMovieFn != nil {
        return m.upsertMovieFn(movie)
    }
    return nil
}
func (m *mockRepo) GetMovieByID(id int64) (*models.Movie, error) {
    if m.getMovieByIDFn != nil {
        return m.getMovieByIDFn(id)
    }
    return nil, errors.New("not implemented")
}
func (m *mockRepo) SearchMovies(query string) ([]models.Movie, error) {
    if m.searchMoviesFn != nil {
        return m.searchMoviesFn(query)
    }
    return nil, nil
}
func (m *mockRepo) ListMovies(offset, limit int) ([]models.Movie, error) {
    if m.listMoviesFn != nil {
        return m.listMoviesFn(offset, limit)
    }
    return nil, nil
}
func (m *mockRepo) ListPopularMovies(limit int) ([]models.Movie, error) {
    if m.listPopularMoviesFn != nil {
        return m.listPopularMoviesFn(limit)
    }
    return nil, nil
}
func (m *mockRepo) AddToWatchlist(item *models.WatchlistItem) error {
    if m.addToWatchlistFn != nil {
        return m.addToWatchlistFn(item)
    }
    return nil
}
func (m *mockRepo) GetWatchlist(userID int64) ([]models.WatchlistItem, error) {
    if m.getWatchlistFn != nil {
        return m.getWatchlistFn(userID)
    }
    return nil, nil
}
func (m *mockRepo) RemoveFromWatchlist(userID, movieID int64) error {
    if m.removeFromWatchFn != nil {
        return m.removeFromWatchFn(userID, movieID)
    }
    return nil
}
func (m *mockRepo) UpsertRating(item *models.RatingItem) error {
    if m.upsertRatingFn != nil {
        return m.upsertRatingFn(item)
    }
    return nil
}
func (m *mockRepo) GetRatings(userID int64) ([]models.RatingItem, error) {
    if m.getRatingsFn != nil {
        return m.getRatingsFn(userID)
    }
    return nil, nil
}
func (m *mockRepo) DeleteRating(userID, movieID int64) error {
    if m.deleteRatingFn != nil {
        return m.deleteRatingFn(userID, movieID)
    }
    return nil
}

type mockKPClient struct {
    searchByKeywordFn func(keyword string, page int) ([]kp.Film, int, error)
    getPopularAllFn   func(page int) ([]kp.Film, int, error)
}

func (m *mockKPClient) SearchByKeyword(keyword string, page int) ([]kp.Film, int, error) {
    if m.searchByKeywordFn != nil {
        return m.searchByKeywordFn(keyword, page)
    }
    return nil, 0, nil
}
func (m *mockKPClient) GetPopularAll(page int) ([]kp.Film, int, error) {
    if m.getPopularAllFn != nil {
        return m.getPopularAllFn(page)
    }
    return nil, 0, nil
}

type mockYTClient struct {
    searchReviewsFn func(keyword string, maxResultsPerPage int) ([]yt.ReviewResult, error)
}

func (m *mockYTClient) SearchReviews(keyword string, maxResultsPerPage int) ([]yt.ReviewResult, error) {
    if m.searchReviewsFn != nil {
        return m.searchReviewsFn(keyword, maxResultsPerPage)
    }
    return nil, nil
}

func newTestService(repo service.RepoIface, kp service.KPIface, yt service.YTIface) *service.Service {
    return service.NewService(repo, kp, yt, "test-secret")
}

func jsonBody(t *testing.T, v any) *bytes.Reader {
    t.Helper()
    b, err := json.Marshal(v)
    if err != nil {
        t.Fatalf("marshal body: %v", err)
    }
    return bytes.NewReader(b)
}

func decodeJSON[T any](t *testing.T, rr *httptest.ResponseRecorder) T {
    t.Helper()
    var out T
    if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
        t.Fatalf("unmarshal response: %v; body=%s", err, rr.Body.String())
    }
    return out
}

func withUserID(ctx context.Context, userID int64) context.Context {
    return context.WithValue(ctx, middleware.UserIDKey, userID)
}
