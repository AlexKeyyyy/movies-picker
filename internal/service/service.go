package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/AlexKeyyyy/movies-picker/internal/models"
	"github.com/AlexKeyyyy/movies-picker/internal/repository"
	"github.com/AlexKeyyyy/movies-picker/pkg/kinopoisk"
	"github.com/AlexKeyyyy/movies-picker/pkg/youtube"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo      *repository.Repo
	kpClient  *kinopoisk.Client
	ytClient  *youtube.Client
	jwtSecret string
}

func NewService(repo *repository.Repo, kp *kinopoisk.Client, yt *youtube.Client, jwtSecret string) *Service {
	return &Service{repo: repo, kpClient: kp, ytClient: yt, jwtSecret: jwtSecret}
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

// GetProfile возвращает профиль текущего пользователя
func (s *Service) GetProfile(userID int64) (*models.User, error) {
	return s.repo.GetUserByID(userID)
}

// UpdateProfile обновляет профиль пользователя (email и/или пароль)
func (s *Service) UpdateProfile(userID int64, newEmail, newPassword string) (*models.User, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	if newEmail != "" && newEmail != user.Email {
		user.Email = newEmail
	}

	if newPassword != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = string(hashed)
	}

	if err := s.repo.UpdateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

// --- Movies ---
func (s *Service) SearchMovies(query string) ([]models.Movie, error) {
	// 1) Сначала пытаемся найти в БД
	movies, err := s.repo.SearchMovies(query)
	if len(movies) > 0 {
		return movies, nil
	}

	// 2) Иначе — ищем по API
	films, totalPages, err := s.kpClient.SearchByKeyword(query, 1)
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
		films, _, err := s.kpClient.SearchByKeyword(query, page)
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

// ListMovies отдаёт фильмы по страницам
func (s *Service) ListMovies(page, size int) ([]models.Movie, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	offset := (page - 1) * size
	return s.repo.ListMovies(offset, size)
}

// ListPopular возвращает топ-N популярных фильмов
func (s *Service) ListPopular(limit int) ([]models.Movie, error) {
	if limit < 1 {
		limit = 10
	}
	return s.repo.ListPopularMovies(limit)
}

// --- Reviews ---

// GetMovieReviews возвращает список обзоров для фильма по его ID
func (s *Service) GetMovieReviews(id int64) ([]models.ReviewItem, error) {
	m, err := s.repo.GetMovieByID(id)
	if err != nil {
		return nil, fmt.Errorf("movie not found: %w", err)
	}

	reviews, err := s.ytClient.SearchReviews(m.Title, 10)
	if err != nil {
		return nil, fmt.Errorf("youtube search failed: %w", err)
	}

	var out []models.ReviewItem
	for _, r := range reviews {
		out = append(out, models.ReviewItem{
			VideoID:      r.VideoID,
			Title:        r.Title,
			ChannelTitle: r.ChannelTitle,
			ThumbnailURL: r.ThumbnailURL,
		})
	}

	return out, nil
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

// DeleteRating удаляет оценку
func (s *Service) DeleteRating(userID, movieID int64) error {
	return s.repo.DeleteRating(userID, movieID)
}
