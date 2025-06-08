package service

import (
	"errors"
	"time"

	"github.com/AlexKeyyyy/movies-picker/internal/models"
	"github.com/AlexKeyyyy/movies-picker/internal/repository"
	"github.com/AlexKeyyyy/movies-picker/pkg/kinopoisk"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo      *repository.Repo
	client    *kinopoisk.Client
	jwtSecret string
}

func NewService(repo *repository.Repo, client *kinopoisk.Client, jwtSecret string) *Service {
	return &Service{repo: repo, client: client, jwtSecret: jwtSecret}
}

// --- Auth ---
func (s *Service) Register(email, password string) (*models.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &models.User{Email: email, PasswordHash: string(hashed)}
	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) Login(email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// --- Movies ---
func (s *Service) SearchMovies(query string) ([]models.Movie, error) {
	// 1) Сначала пытаемся найти в БД
	movies, err := s.repo.SearchMovies(query)
	if len(movies) > 0 {
		return movies, nil
	}

	// 2) Иначе — ищем по API
	films, totalPages, err := s.client.SearchByKeyword(query, 1)
	if err != nil {
		return nil, err
	}

	var result []models.Movie
	for _, f := range films {
		m := s.mapFilmToModel(f)
		_ = s.repo.UpsertMovie(&m)
		result = append(result, m)
	}

	for page := 2; page <= totalPages; page++ {
		films, _, err := s.client.SearchByKeyword(query, page)
		if err != nil {
			continue
		}
		for _, f := range films {
			m := s.mapFilmToModel(f)
			_ = s.repo.UpsertMovie(&m)
			result = append(result, m)
		}
	}
	return result, nil
}

func (s *Service) mapFilmToModel(f kinopoisk.Film) models.Movie {
	yearInt, _ := f.Year.Int64()
	return models.Movie{
		ID:          f.KinopoiskID,
		Title:       f.NameRu,
		Year:        int(yearInt),
		Description: f.Description,
		PosterURL:   f.PosterURL,
	}
}

func (s *Service) GetMovie(id int64) (*models.Movie, error) {
	return s.repo.GetMovieByID(id)
}

// --- Reviews ---
func (s *Service) GetMovieReviews(id int64) ([]models.ReviewItem, error) {
	// пока пустой список
	// TODO:
	// 1) получить фильм: m, err := s.repo.GetMovieByID(id)
	// 2) вызвать клиент YouTube: reviews, err := youtubeClient.SearchReviews(m.Title)
	return []models.ReviewItem{}, nil
}

// --- Watchlist ---
func (s *Service) AddToWatchlist(userID, movieID int64) error {
	item := &models.WatchlistItem{
		UserID:  userID,
		MovieID: movieID,
	}
	return s.repo.AddToWatchlist(item)
}

func (s *Service) GetWatchlist(userID int64) ([]models.WatchlistItem, error) {
	return s.repo.GetWatchlist(userID)
}

func (s *Service) RemoveFromWatchlist(userID, movieID int64) error {
	return s.repo.RemoveFromWatchlist(userID, movieID)
}

// --- Ratings ---
func (s *Service) UpsertRating(item *models.RatingItem) error {
	return s.repo.UpsertRating(item)
}

func (s *Service) GetRatings(userID int64) ([]models.RatingItem, error) {
	return s.repo.GetRatings(userID)
}
