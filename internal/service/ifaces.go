package service

import (
	"github.com/AlexKeyyyy/movies-picker/internal/models"
	kp "github.com/AlexKeyyyy/movies-picker/pkg/kinopoisk"
	yt "github.com/AlexKeyyyy/movies-picker/pkg/youtube"
)

// RepoIface — контракт для слоя БД; *repository.Repo уже удовлетворяет этому интерфейсу
type RepoIface interface {
	CreateUser(u *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(userID int64) (*models.User, error)
	UpdateUser(user *models.User) error
	UpsertMovie(m *models.Movie) error
	GetMovieByID(id int64) (*models.Movie, error)
	SearchMovies(query string) ([]models.Movie, error)
	ListMovies(offset, limit int) ([]models.Movie, error)
	ListPopularMovies(limit int) ([]models.Movie, error)
	AddToWatchlist(item *models.WatchlistItem) error
	GetWatchlist(userID int64) ([]models.WatchlistItem, error)
	RemoveFromWatchlist(userID, movieID int64) error
	UpsertRating(item *models.RatingItem) error
	GetRatings(userID int64) ([]models.RatingItem, error)
	DeleteRating(userID, movieID int64) error
}

// KPIface — контракт для клиента Кинопоиска
type KPIface interface {
	SearchByKeyword(keyword string, page int) ([]kp.Film, int, error)
	GetPopularAll(page int) ([]kp.Film, int, error)
}

// YTIface — контракт для клиента YouTube
type YTIface interface {
	SearchReviews(keyword string, maxResultsPerPage int) ([]yt.ReviewResult, error)
}

type AuthIface interface {
	Register(email, password string) (*models.User, error)
	Login(email, password string) (string, error)
}

type UserIface interface {
	GetProfile(userID int64) (*models.User, error)
	UpdateProfile(userID int64, email, password string) (*models.User, error)
}
