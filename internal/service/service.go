package service

import (
	"errors"
	"time"

	"github.com/AlexKeyyyy/movies-picker/internal/models"
	"github.com/AlexKeyyyy/movies-picker/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo      *repository.Repo
	jwtSecret string
}

func NewService(repo *repository.Repo, jwtSecret string) *Service {
	return &Service{repo: repo, jwtSecret: jwtSecret}
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
	return s.repo.SearchMovies(query)
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
