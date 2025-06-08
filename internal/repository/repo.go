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

// --- Movie ---
func (r *Repo) UpsertMovie(m *models.Movie) error {
	// INSERT … ON CONFLICT (movie_id) DO UPDATE …
	_, err := r.db.NamedExec(`
        INSERT INTO movies (movie_id, title, year, poster_url, description, last_sync)
        VALUES (:movie_id,:title,:year,:poster_url,:description,NOW())
        ON CONFLICT (movie_id) DO UPDATE
          SET title=:title,year=:year,poster_url=:poster_url,description=:description,last_sync=NOW()
    `, m)
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
