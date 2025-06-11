package repository

import (
	"github.com/AlexKeyyyy/movies-picker/internal/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Repo struct {
	db *sqlx.DB
}

func NewRepo(dbURL string) (*Repo, error) {
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	return &Repo{db: db}, nil
}

// --- User ---
func (r *Repo) CreateUser(u *models.User) error {
	return r.db.Get(&u.ID,
		`INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING user_id`,
		u.Email, u.PasswordHash)
}

func (r *Repo) GetUserByEmail(email string) (*models.User, error) {
	var u models.User
	err := r.db.Get(&u, "SELECT * FROM users WHERE email=$1", email)
	return &u, err
}

// GetUserByID возвращает пользователя по ID
func (r *Repo) GetUserByID(userID int64) (*models.User, error) {
	var user models.User
	err := r.db.Get(&user, "SELECT user_id, email, created_at FROM users WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser обновляет email и/или пароль пользователя
// UpdateUser обновляет email и/или пароль пользователя
func (r *Repo) UpdateUser(user *models.User) error {
	_, err := r.db.Exec(
		"UPDATE users SET email = $1, password_hash = $2 WHERE user_id = $3",
		user.Email, user.PasswordHash, user.ID,
	)
	return err
}

// --- Movie ---
func (r *Repo) UpsertMovie(m *models.Movie) error {
	// INSERT … ON CONFLICT (movie_id) DO UPDATE …
	_, err := r.db.NamedExec(`
      INSERT INTO movies
        (movie_id, title, year, poster_url, description, rating_kinopoisk, last_sync)
      VALUES
        (:movie_id, :title, :year, :poster_url, :description, :rating_kinopoisk, NOW())
      ON CONFLICT (movie_id) DO UPDATE SET
        title            = EXCLUDED.title,
        year             = EXCLUDED.year,
        poster_url       = EXCLUDED.poster_url,
        description      = EXCLUDED.description,
        rating_kinopoisk = EXCLUDED.rating_kinopoisk,
        last_sync        = NOW()`,
		m,
	)
	return err
}

func (r *Repo) GetMovieByID(id int64) (*models.Movie, error) {
	var m models.Movie
	err := r.db.Get(&m, "SELECT * FROM movies WHERE movie_id=$1", id)
	return &m, err
}

func (r *Repo) SearchMovies(query string) ([]models.Movie, error) {
	// полнотекстовый поиск или ILIKE
	var movies []models.Movie
	err := r.db.Select(&movies, "SELECT * FROM movies WHERE title ILIKE $1", "%"+query+"%")
	return movies, err
}

// ListMovies возвращает список фильмов с пагинацией
func (r *Repo) ListMovies(offset, limit int) ([]models.Movie, error) {
	var movies []models.Movie
	// выбираем только поля нужные для списка
	query := `
      SELECT movie_id, title, year, poster_url
      FROM movies
      ORDER BY title
      LIMIT $1 OFFSET $2`
	if err := r.db.Select(&movies, query, limit, offset); err != nil {
		return nil, err
	}
	return movies, nil
}

// ListPopularMovies возвращает топ-N фильмов по рейтингу
func (r *Repo) ListPopularMovies(limit int) ([]models.Movie, error) {
	var movies []models.Movie
	err := r.db.Select(&movies,
		`SELECT movie_id, title, year, poster_url, rating_kinopoisk
         FROM movies ORDER BY rating_kinopoisk DESC NULLS LAST LIMIT $1`,
		limit,
	)
	return movies, err
}

// --- Watchlist ---
func (r *Repo) AddToWatchlist(item *models.WatchlistItem) error {
	_, err := r.db.Exec(`
        INSERT INTO watchlist (user_id, movie_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`,
		item.UserID, item.MovieID)
	return err
}

func (r *Repo) GetWatchlist(userID int64) ([]models.WatchlistItem, error) {
	var list []models.WatchlistItem
	err := r.db.Select(&list, `
        SELECT w.movie_id, w.added_at, m.title, m.poster_url
        FROM watchlist w JOIN movies m ON w.movie_id = m.movie_id
        WHERE w.user_id = $1`, userID)
	return list, err
}

func (r *Repo) RemoveFromWatchlist(userID, movieID int64) error {
	_, err := r.db.Exec(
		"DELETE FROM watchlist WHERE user_id=$1 AND movie_id=$2",
		userID, movieID)
	return err
}

// --- Ratings ---
func (r *Repo) UpsertRating(item *models.RatingItem) error {
	_, err := r.db.Exec(`
        INSERT INTO ratings (user_id, movie_id, rating) VALUES ($1,$2,$3)
        ON CONFLICT (user_id,movie_id) DO UPDATE SET rating = $3, rated_at = NOW()`,
		item.UserID, item.MovieID, item.Rating)
	return err
}

func (r *Repo) GetRatings(userID int64) ([]models.RatingItem, error) {
	var list []models.RatingItem
	err := r.db.Select(&list,
		"SELECT movie_id, rating, rated_at FROM ratings WHERE user_id=$1", userID)
	return list, err
}

// DeleteRating удаляет оценку пользователя для фильма
func (r *Repo) DeleteRating(userID, movieID int64) error {
	_, err := r.db.Exec(
		"DELETE FROM ratings WHERE user_id = $1 AND movie_id = $2",
		userID, movieID,
	)
	return err
}
