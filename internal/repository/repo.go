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
