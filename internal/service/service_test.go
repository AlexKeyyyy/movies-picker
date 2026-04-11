package service_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/AlexKeyyyy/movies-picker/internal/models"
	"github.com/AlexKeyyyy/movies-picker/internal/service"
	kp "github.com/AlexKeyyyy/movies-picker/pkg/kinopoisk"
	yt "github.com/AlexKeyyyy/movies-picker/pkg/youtube"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// ─────────────────────────────────────────────────────────────
//  MOCK-ОБЪЕКТЫ (заглушки зависимостей)
// ─────────────────────────────────────────────────────────────

// MockRepo — заглушка репозитория
type MockRepo struct{ mock.Mock }

func (m *MockRepo) CreateUser(u *models.User) error {
	return m.Called(u).Error(0)
}
func (m *MockRepo) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *MockRepo) GetUserByID(userID int64) (*models.User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *MockRepo) UpdateUser(user *models.User) error {
	return m.Called(user).Error(0)
}
func (m *MockRepo) UpsertMovie(mv *models.Movie) error {
	return m.Called(mv).Error(0)
}
func (m *MockRepo) GetMovieByID(id int64) (*models.Movie, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Movie), args.Error(1)
}
func (m *MockRepo) SearchMovies(query string) ([]models.Movie, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Movie), args.Error(1)
}
func (m *MockRepo) ListMovies(offset, limit int) ([]models.Movie, error) {
	args := m.Called(offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Movie), args.Error(1)
}
func (m *MockRepo) ListPopularMovies(limit int) ([]models.Movie, error) {
	args := m.Called(limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Movie), args.Error(1)
}
func (m *MockRepo) AddToWatchlist(item *models.WatchlistItem) error {
	return m.Called(item).Error(0)
}
func (m *MockRepo) GetWatchlist(userID int64) ([]models.WatchlistItem, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.WatchlistItem), args.Error(1)
}
func (m *MockRepo) RemoveFromWatchlist(userID, movieID int64) error {
	return m.Called(userID, movieID).Error(0)
}
func (m *MockRepo) UpsertRating(item *models.RatingItem) error {
	return m.Called(item).Error(0)
}
func (m *MockRepo) GetRatings(userID int64) ([]models.RatingItem, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.RatingItem), args.Error(1)
}
func (m *MockRepo) DeleteRating(userID, movieID int64) error {
	return m.Called(userID, movieID).Error(0)
}

// MockKP — заглушка клиента Кинопоиска
type MockKP struct{ mock.Mock }

func (m *MockKP) SearchByKeyword(keyword string, page int) ([]kp.Film, int, error) {
	args := m.Called(keyword, page)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]kp.Film), args.Int(1), args.Error(2)
}
func (m *MockKP) GetPopularAll(page int) ([]kp.Film, int, error) {
	args := m.Called(page)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]kp.Film), args.Int(1), args.Error(2)
}

// MockYT — заглушка клиента YouTube
type MockYT struct{ mock.Mock }

func (m *MockYT) SearchReviews(keyword string, max int) ([]yt.ReviewResult, error) {
	args := m.Called(keyword, max)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]yt.ReviewResult), args.Error(1)
}

// вспомогательная функция создания сервиса для тестов
func newTestSvc(repo *MockRepo, kpMock *MockKP, ytMock *MockYT) *service.Service {
	return service.NewService(repo, kpMock, ytMock, "test-jwt-secret")
}

// ─────────────────────────────────────────────────────────────
//  ТЕСТЫ 1–2: Register
// ─────────────────────────────────────────────────────────────

// Тест 1 — Успешная регистрация пользователя
// Техника: класс эквивалентности — валидные входные данные
func TestRegister_Success(t *testing.T) {
	repo := new(MockRepo)
	svc := newTestSvc(repo, nil, nil)

	// Имитируем: БД принимает запрос и возвращает ID = 42
	repo.On("CreateUser", mock.AnythingOfType("*models.User")).
		Run(func(args mock.Arguments) {
			// Симулируем RETURNING user_id из БД
			u := args.Get(0).(*models.User)
			u.ID = 42
		}).
		Return(nil)

	user, err := svc.Register("alice@example.com", "StrongPass1!")

	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "alice@example.com", user.Email)
	assert.Equal(t, int64(42), user.ID)
	// Пароль должен быть захеширован — не хранится открытым текстом
	assert.NotEqual(t, "StrongPass1!", user.PasswordHash, "пароль не должен храниться открытым текстом")
	assert.NotEmpty(t, user.PasswordHash, "хэш пароля должен быть заполнен")
	repo.AssertExpectations(t)
}

// Тест 2 — Ошибка регистрации: дублирующийся email
// Техника: негативный сценарий — ошибка уникальности в БД
func TestRegister_DuplicateEmail(t *testing.T) {
	repo := new(MockRepo)
	svc := newTestSvc(repo, nil, nil)

	repo.On("CreateUser", mock.AnythingOfType("*models.User")).
		Return(errors.New("duplicate key value violates unique constraint"))

	user, err := svc.Register("existing@example.com", "Pass123!")

	assert.Error(t, err)
	assert.Nil(t, user)
	repo.AssertExpectations(t)
}

// ─────────────────────────────────────────────────────────────
//  ТЕСТЫ 3–5: Login
// ─────────────────────────────────────────────────────────────

// Тест 3 — Успешная авторизация
// Техника: класс эквивалентности — верный email + верный пароль
func TestLogin_Success(t *testing.T) {
	repo := new(MockRepo)
	svc := newTestSvc(repo, nil, nil)

	// Генерируем реальный bcrypt-хэш, как это делает сервис при регистрации
	hashedPass, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	require.NoError(t, err)

	repo.On("GetUserByEmail", "bob@example.com").
		Return(&models.User{ID: 1, Email: "bob@example.com", PasswordHash: string(hashedPass)}, nil)

	token, err := svc.Login("bob@example.com", "password123")

	require.NoError(t, err)
	assert.NotEmpty(t, token, "должен быть выдан JWT-токен")
	repo.AssertExpectations(t)
}

// Тест 4 — Пользователь не найден при логине
// Техника: негативный сценарий — несуществующий email
func TestLogin_UserNotFound(t *testing.T) {
	repo := new(MockRepo)
	svc := newTestSvc(repo, nil, nil)

	repo.On("GetUserByEmail", "ghost@example.com").
		Return(nil, errors.New("sql: no rows in result set"))

	token, err := svc.Login("ghost@example.com", "anypassword")

	assert.Error(t, err)
	assert.Empty(t, token)
	repo.AssertExpectations(t)
}

// Тест 5 — Неверный пароль при логине
// Техника: класс эквивалентности — верный email, НЕверный пароль
func TestLogin_WrongPassword(t *testing.T) {
	repo := new(MockRepo)
	svc := newTestSvc(repo, nil, nil)

	hashedPass, _ := bcrypt.GenerateFromPassword([]byte("correct_password"), bcrypt.DefaultCost)
	repo.On("GetUserByEmail", "carol@example.com").
		Return(&models.User{ID: 2, Email: "carol@example.com", PasswordHash: string(hashedPass)}, nil)

	token, err := svc.Login("carol@example.com", "wrong_password")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.EqualError(t, err, "invalid credentials")
	repo.AssertExpectations(t)
}

// ─────────────────────────────────────────────────────────────
//  ТЕСТЫ 6–11: UpdateProfile / GetProfile
// ─────────────────────────────────────────────────────────────

// Тест 6 — Успешное получение профиля
func TestGetProfile_Success(t *testing.T) {
	repo := new(MockRepo)
	svc := newTestSvc(repo, nil, nil)

	expected := &models.User{ID: 5, Email: "dave@example.com"}
	repo.On("GetUserByID", int64(5)).Return(expected, nil)

	user, err := svc.GetProfile(5)

	require.NoError(t, err)
	assert.Equal(t, expected, user)
	repo.AssertExpectations(t)
}

// Тест 7 — Профиль не найден
// Техника: негативный сценарий
func TestGetProfile_NotFound(t *testing.T) {
	repo := new(MockRepo)
	svc := newTestSvc(repo, nil, nil)

	repo.On("GetUserByID", int64(999)).Return(nil, errors.New("user not found"))

	user, err := svc.GetProfile(999)

	assert.Error(t, err)
	assert.Nil(t, user)
	repo.AssertExpectations(t)
}

// Тест 8 — Обновление профиля: смена email
// Техника: класс эквивалентности — изменяется только email
func TestUpdateProfile_ChangeEmail(t *testing.T) {
	repo := new(MockRepo)
	svc := newTestSvc(repo, nil, nil)

	original := &models.User{ID: 1, Email: "old@example.com", PasswordHash: "some_hash"}
	repo.On("GetUserByID", int64(1)).Return(original, nil)
	repo.On("UpdateUser", mock.AnythingOfType("*models.User")).Return(nil)

	user, err := svc.UpdateProfile(1, "new@example.com", "")

	require.NoError(t, err)
	assert.Equal(t, "new@example.com", user.Email)
	assert.Equal(t, "some_hash", user.PasswordHash, "хэш не должен меняться при смене только email")
	repo.AssertExpectations(t)
}

// Тест 9 — Обновление профиля: смена пароля
// Техника: класс эквивалентности — изменяется только пароль
func TestUpdateProfile_ChangePassword(t *testing.T) {
	repo := new(MockRepo)
	svc := newTestSvc(repo, nil, nil)

	original := &models.User{ID: 1, Email: "user@example.com", PasswordHash: "old_hash"}
	repo.On("GetUserByID", int64(1)).Return(original, nil)
	repo.On("UpdateUser", mock.AnythingOfType("*models.User")).Return(nil)

	user, err := svc.UpdateProfile(1, "", "NewPass123!")

	require.NoError(t, err)
	assert.Equal(t, "user@example.com", user.Email, "email не должен меняться")
	assert.NotEqual(t, "old_hash", user.PasswordHash, "хэш пароля должен обновиться")
	// Проверяем, что новый хэш — реальный bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("NewPass123!"))
	assert.NoError(t, err, "новый хэш должен соответствовать новому паролю")
	repo.AssertExpectations(t)
}

// Тест 10 — Обновление профиля: пустые поля (нет изменений)
// Техника: граничное условие — оба параметра пустые строки
func TestUpdateProfile_NoChanges(t *testing.T) {
	repo := new(MockRepo)
	svc := newTestSvc(repo, nil, nil)

	original := &models.User{ID: 1, Email: "same@example.com", PasswordHash: "stable_hash"}
	repo.On("GetUserByID", int64(1)).Return(original, nil)
	repo.On("UpdateUser", mock.AnythingOfType("*models.User")).Return(nil)

	user, err := svc.UpdateProfile(1, "", "")

	require.NoError(t, err)
	assert.Equal(t, "same@example.com", user.Email)
	assert.Equal(t, "stable_hash", user.PasswordHash, "ничего не должно измениться")
	repo.AssertExpectations(t)
}

// Тест 11 — Обновление профиля: пользователь не найден
// Техника: негативный сценарий
func TestUpdateProfile_SameEmailKeepsValue(t *testing.T) {
	repo := new(MockRepo)
	svc := newTestSvc(repo, nil, nil)

	original := &models.User{ID: 1, Email: "same@example.com", PasswordHash: "stable_hash"}
	repo.On("GetUserByID", int64(1)).Return(original, nil)
	repo.On("UpdateUser", mock.AnythingOfType("*models.User")).Return(nil)

	user, err := svc.UpdateProfile(1, "same@example.com", "")

	require.NoError(t, err)
	assert.Equal(t, "same@example.com", user.Email)
	assert.Equal(t, "stable_hash", user.PasswordHash)
	repo.AssertExpectations(t)
}

func TestUpdateProfile_UserNotFound(t *testing.T) {
	repo := new(MockRepo)
	svc := newTestSvc(repo, nil, nil)

	repo.On("GetUserByID", int64(404)).Return(nil, errors.New("not found"))

	user, err := svc.UpdateProfile(404, "x@x.com", "")

	assert.Error(t, err)
	assert.Nil(t, user)
	repo.AssertExpectations(t)
}

func TestUpdateProfile_UpdateUserError(t *testing.T) {
	repo := new(MockRepo)
	svc := newTestSvc(repo, nil, nil)

	original := &models.User{ID: 1, Email: "user@example.com", PasswordHash: "hash"}
	repo.On("GetUserByID", int64(1)).Return(original, nil)
	repo.On("UpdateUser", mock.AnythingOfType("*models.User")).Return(errors.New("update failed"))

	user, err := svc.UpdateProfile(1, "new@example.com", "")

	assert.Error(t, err)
	assert.Nil(t, user)
	repo.AssertExpectations(t)
}

// ─────────────────────────────────────────────────────────────
//  ТЕСТЫ 12–16: ListMovies (граничные условия пагинации)
// ─────────────────────────────────────────────────────────────

// Тест 12 — Стандартный запрос с валидными параметрами
// Техника: класс эквивалентности — нормальные значения
func TestListMovies_ValidPage(t *testing.T) {
	repo := new(MockRepo)
	svc := newTestSvc(repo, nil, nil)

	movies := []models.Movie{{ID: 1, Title: "Movie A"}, {ID: 2, Title: "Movie B"}}
	// page=2, size=10 → offset = (2-1)*10 = 10
	repo.On("ListMovies", 10, 10).Return(movies, nil)

	result, err := svc.ListMovies(2, 10)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

// Тест 13 — page=0 должен автоматически стать page=1
// Техника: граничное условие — нижняя граница номера страницы
func TestListMovies_PageZeroBecomesOne(t *testing.T) {
	repo := new(MockRepo)
	svc := newTestSvc(repo, nil, nil)

	// page=0 → исправляется до 1 → offset = (1-1)*20 = 0
	repo.On("ListMovies", 0, 20).Return([]models.Movie{}, nil)

	result, err := svc.ListMovies(0, 20)

	require.NoError(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t) // убеждаемся, что вызван именно с offset=0
}

// Тест 14 — size=0 должен автоматически стать size=20
// Техника: граничное условие — нижняя граница размера страницы
func TestListMovies_SizeZeroBecomesDefault(t *testing.T) {
	repo := new(MockRepo)
	svc := newTestSvc(repo, nil, nil)

	// size=0 → исправляется до 20; page=1 → offset=0
	repo.On("ListMovies", 0, 20).Return([]models.Movie{}, nil)

	_, err := svc.ListMovies(1, 0)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

// Тест 15 — size=200 должен автоматически стать size=20
// Техника: граничное условие — верхняя граница (> 100)
func TestListMovies_SizeExceedsMaxBecomesDefault(t *testing.T) {
	repo := new(MockRepo)
	svc := newTestSvc(repo, nil, nil)

	// size=200 > 100 → исправляется до 20
	repo.On("ListMovies", 0, 20).Return([]models.Movie{}, nil)

	_, err := svc.ListMovies(1, 200)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestListMovies_SizeExceedsMinBecomesDefault(t *testing.T) {
	repo := new(MockRepo)
	svc := newTestSvc(repo, nil, nil)

	repo.On("ListMovies", 0, 20).Return([]models.Movie{}, nil)

	_, err := svc.ListMovies(1, 0)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestListMovies_Normal1(t *testing.T) {
	repo := new(MockRepo)
	svc := newTestSvc(repo, nil, nil)

	repo.On("ListMovies", 0, 1).Return([]models.Movie{}, nil)

	_, err := svc.ListMovies(1, 1)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestListMovies_Normal36(t *testing.T) {
	repo := new(MockRepo)
	svc := newTestSvc(repo, nil, nil)

	repo.On("ListMovies", 0, 36).Return([]models.Movie{}, nil)

	_, err := svc.ListMovies(1, 36)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestListMovies_Normal100(t *testing.T) {
	repo := new(MockRepo)
	svc := newTestSvc(repo, nil, nil)

	repo.On("ListMovies", 0, 100).Return([]models.Movie{}, nil)

	_, err := svc.ListMovies(1, 100)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

// Тест 16 — Ошибка репозитория при получении списка
// Техника: негативный сценарий — сбой БД
func TestListMovies_RepoError(t *testing.T) {
	repo := new(MockRepo)
	svc := newTestSvc(repo, nil, nil)

	repo.On("ListMovies", mock.Anything, mock.Anything).
		Return(nil, errors.New("db connection refused"))

	result, err := svc.ListMovies(1, 10)

	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ─────────────────────────────────────────────────────────────
//  ТЕСТЫ 17–18: ListPopular
// ─────────────────────────────────────────────────────────────

// Тест 17 — Успешное получение популярных фильмов
func TestListPopular_Success(t *testing.T) {
	repo := new(MockRepo)
	svc := newTestSvc(repo, nil, nil)

	movies := []models.Movie{
		{ID: 10, Title: "Inception", RatingKinopoisk: 9.5},
		{ID: 11, Title: "Interstellar", RatingKinopoisk: 8.8},
	}
	repo.On("ListPopularMovies", 5).Return(movies, nil)

	result, err := svc.ListPopular(5)

	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, float64(9.5), result[0].RatingKinopoisk)
	repo.AssertExpectations(t)
}

// Тест 18 — limit=0 автоматически становится 10
// Техника: граничное условие — нулевой лимит
func TestListPopular_ZeroLimitBecomesDefault(t *testing.T) {
	repo := new(MockRepo)
	svc := newTestSvc(repo, nil, nil)

	// limit=0 → исправляется до 10
	repo.On("ListPopularMovies", 10).Return([]models.Movie{}, nil)

	_, err := svc.ListPopular(0)

	require.NoError(t, err)
	// AssertExpectations проверит, что был вызов именно с 10, а не 0
	repo.AssertExpectations(t)
}

// ─────────────────────────────────────────────────────────────
//  ТЕСТЫ 19–20: GetMovie
// ─────────────────────────────────────────────────────────────

// Тест 19 — Успешное получение фильма по ID
func TestGetMovie_Success(t *testing.T) {
	repo := new(MockRepo)
	svc := newTestSvc(repo, nil, nil)

	expected := &models.Movie{ID: 301, Title: "The Matrix", Year: 1999}
	repo.On("GetMovieByID", int64(301)).Return(expected, nil)

	movie, err := svc.GetMovie(301)

	require.NoError(t, err)
	assert.Equal(t, expected, movie)
	repo.AssertExpectations(t)
}

// Тест 20 — Фильм не найден (несуществующий ID)
// Техника: негативный сценарий
func TestGetMovie_NotFound(t *testing.T) {
	repo := new(MockRepo)
	svc := newTestSvc(repo, nil, nil)

	repo.On("GetMovieByID", int64(99999)).Return(nil, errors.New("sql: no rows in result set"))

	movie, err := svc.GetMovie(99999)

	assert.Error(t, err)
	assert.Nil(t, movie)
	repo.AssertExpectations(t)
}

// ─────────────────────────────────────────────────────────────
//  ТЕСТЫ 21–23: GetMovieReviews
// ─────────────────────────────────────────────────────────────

// Тест 21 — Успешное получение обзоров (связка repo + youtube)
// Техника: интеграционный UoW — два мока взаимодействуют
func TestGetMovieReviews_Success(t *testing.T) {
	repo := new(MockRepo)
	ytMock := new(MockYT)
	svc := newTestSvc(repo, nil, ytMock)

	movie := &models.Movie{ID: 1, Title: "Fight Club"}
	reviews := []yt.ReviewResult{
		{VideoID: "abc123", Title: "Fight Club — Разбор", VideoURL: "https://youtube.com/watch?v=abc123"},
		{VideoID: "def456", Title: "Анализ финала", VideoURL: "https://youtube.com/watch?v=def456"},
	}
	repo.On("GetMovieByID", int64(1)).Return(movie, nil)
	ytMock.On("SearchReviews", "Fight Club", 10).Return(reviews, nil)

	result, err := svc.GetMovieReviews(1)

	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "abc123", result[0].VideoID)
	assert.Equal(t, "Fight Club — Разбор", result[0].Title)
	repo.AssertExpectations(t)
	ytMock.AssertExpectations(t)
}

// Тест 22 — Фильм не найден при запросе обзоров
// Техника: негативный сценарий — первый шаг цепочки падает
func TestGetMovieReviews_MovieNotFound(t *testing.T) {
	repo := new(MockRepo)
	svc := newTestSvc(repo, nil, nil)

	repo.On("GetMovieByID", int64(777)).Return(nil, errors.New("movie not found"))

	result, err := svc.GetMovieReviews(777)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "movie not found")
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// Тест 23 — Ошибка YouTube API при корректном фильме
// Техника: негативный сценарий — второй шаг цепочки падает
func TestGetMovieReviews_YouTubeError(t *testing.T) {
	repo := new(MockRepo)
	ytMock := new(MockYT)
	svc := newTestSvc(repo, nil, ytMock)

	movie := &models.Movie{ID: 2, Title: "Memento"}
	repo.On("GetMovieByID", int64(2)).Return(movie, nil)
	ytMock.On("SearchReviews", "Memento", 10).
		Return(nil, errors.New("youtube: quota exceeded"))

	result, err := svc.GetMovieReviews(2)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "youtube search failed")
	assert.Nil(t, result)
	repo.AssertExpectations(t)
	ytMock.AssertExpectations(t)
}

// ─────────────────────────────────────────────────────────────
//  ТЕСТ 24: Watchlist
// ─────────────────────────────────────────────────────────────

// Тест 24 — Успешное добавление фильма в список «Смотреть позже»
func TestAddToWatchlist_Success(t *testing.T) {
	repo := new(MockRepo)
	svc := newTestSvc(repo, nil, nil)

	repo.On("AddToWatchlist", mock.AnythingOfType("*models.WatchlistItem")).Return(nil)

	err := svc.AddToWatchlist(1, 100)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

// ─────────────────────────────────────────────────────────────
//  ТЕСТ 25: Ratings
// ─────────────────────────────────────────────────────────────

// Тест 25 — Полный цикл рейтинга: выставить → получить → удалить
// Техника: попарное тестирование — покрываем три операции одной связанной цепочкой
func TestSearchMovies_ReturnsRepoResults(t *testing.T) {
	repo := new(MockRepo)
	kpMock := new(MockKP)
	svc := newTestSvc(repo, kpMock, nil)

	expected := []models.Movie{
		{ID: 1, Title: "Repo Movie"},
	}

	repo.On("SearchMovies", "repo").Return(expected, nil)

	result, err := svc.SearchMovies("repo")

	require.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
	kpMock.AssertNotCalled(t, "SearchByKeyword", mock.Anything, mock.Anything)
}

func TestSearchMovies_FallsBackToKinopoiskAndAggregatesPages(t *testing.T) {
	repo := new(MockRepo)
	kpMock := new(MockKP)
	svc := newTestSvc(repo, kpMock, nil)

	pageOne := []kp.Film{
		{
			KinopoiskID: 101,
			NameRu:      "Movie One",
			Year:        json.Number("2020"),
			Description: "First page film",
			PosterURL:   "https://img/1.jpg",
		},
	}
	pageTwo := []kp.Film{
		{
			KinopoiskID: 202,
			NameRu:      "Movie Two",
			Year:        json.Number("2021"),
			Description: "Second page film",
			PosterURL:   "https://img/2.jpg",
		},
	}

	repo.On("SearchMovies", "matrix").Return([]models.Movie{}, nil)
	kpMock.On("SearchByKeyword", "matrix", 1).Return(pageOne, 2, nil)
	kpMock.On("SearchByKeyword", "matrix", 2).Return(pageTwo, 2, nil)
	repo.On("UpsertMovie", mock.MatchedBy(func(m *models.Movie) bool {
		return m.ID == 101 && m.Title == "Movie One" && m.Year == 2020
	})).Return(nil)
	repo.On("UpsertMovie", mock.MatchedBy(func(m *models.Movie) bool {
		return m.ID == 202 && m.Title == "Movie Two" && m.Year == 2021
	})).Return(nil)

	result, err := svc.SearchMovies("matrix")

	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, int64(101), result[0].ID)
	assert.Equal(t, int64(202), result[1].ID)
	assert.Equal(t, "First page film", result[0].Description)
	assert.Equal(t, "https://img/2.jpg", result[1].PosterURL)
	repo.AssertExpectations(t)
	kpMock.AssertExpectations(t)
}

func TestSearchMovies_KinopoiskFirstPageError(t *testing.T) {
	repo := new(MockRepo)
	kpMock := new(MockKP)
	svc := newTestSvc(repo, kpMock, nil)

	repo.On("SearchMovies", "broken").Return([]models.Movie{}, nil)
	kpMock.On("SearchByKeyword", "broken", 1).Return(nil, 0, errors.New("kinopoisk unavailable"))

	result, err := svc.SearchMovies("broken")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "kinopoisk unavailable")
	repo.AssertExpectations(t)
	kpMock.AssertExpectations(t)
}

func TestSearchMovies_SkipsBrokenSecondPage(t *testing.T) {
	repo := new(MockRepo)
	kpMock := new(MockKP)
	svc := newTestSvc(repo, kpMock, nil)

	pageOne := []kp.Film{
		{
			KinopoiskID: 303,
			NameRu:      "Page One",
			Year:        json.Number("2022"),
		},
	}

	repo.On("SearchMovies", "partial").Return([]models.Movie{}, nil)
	kpMock.On("SearchByKeyword", "partial", 1).Return(pageOne, 2, nil)
	kpMock.On("SearchByKeyword", "partial", 2).Return(nil, 2, errors.New("page 2 failed"))
	repo.On("UpsertMovie", mock.MatchedBy(func(m *models.Movie) bool {
		return m.ID == 303 && m.Title == "Page One" && m.Year == 2022
	})).Return(nil)

	result, err := svc.SearchMovies("partial")

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, int64(303), result[0].ID)
	repo.AssertExpectations(t)
	kpMock.AssertExpectations(t)
}

func TestGetWatchlist_Success(t *testing.T) {
	repo := new(MockRepo)
	svc := newTestSvc(repo, nil, nil)

	expected := []models.WatchlistItem{
		{MovieID: 10, Title: "Interstellar"},
		{MovieID: 11, Title: "Arrival"},
	}
	repo.On("GetWatchlist", int64(5)).Return(expected, nil)

	result, err := svc.GetWatchlist(5)

	require.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestRemoveFromWatchlist_Success(t *testing.T) {
	repo := new(MockRepo)
	svc := newTestSvc(repo, nil, nil)

	repo.On("RemoveFromWatchlist", int64(5), int64(10)).Return(nil)

	err := svc.RemoveFromWatchlist(5, 10)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestRatings_FullCycle(t *testing.T) {
	repo := new(MockRepo)
	svc := newTestSvc(repo, nil, nil)

	item := &models.RatingItem{UserID: 7, MovieID: 50, Rating: 9}

	// Шаг 1: добавить рейтинг
	repo.On("UpsertRating", item).Return(nil)
	err := svc.UpsertRating(item)
	require.NoError(t, err, "добавление рейтинга должно пройти успешно")

	// Шаг 2: получить список рейтингов
	repo.On("GetRatings", int64(7)).
		Return([]models.RatingItem{{UserID: 7, MovieID: 50, Rating: 9}}, nil)
	ratings, err := svc.GetRatings(7)
	require.NoError(t, err)
	require.Len(t, ratings, 1)
	assert.Equal(t, 9, ratings[0].Rating)

	// Шаг 3: удалить рейтинг
	repo.On("DeleteRating", int64(7), int64(50)).Return(nil)
	err = svc.DeleteRating(7, 50)
	require.NoError(t, err, "удаление рейтинга должно пройти успешно")

	repo.AssertExpectations(t)
}
