package repository

import (
	"errors"
	"testing"
	"time"

	"github.com/AlexKeyyyy/movies-picker/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repo := &Repo{db: sqlxDB}

	tests := []struct {
		name    string
		user    *models.User
		mock    func()
		wantErr bool
	}{
		{
			name: "Success",
			user: &models.User{
				Email:        "test@example.com",
				PasswordHash: "hashed_password",
			},
			mock: func() {
				mock.ExpectQuery(`INSERT INTO users \(email, password_hash\) VALUES \(\$1, \$2\) RETURNING user_id`).
					WithArgs("test@example.com", "hashed_password").
					WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))
			},
			wantErr: false,
		},
		{
			name: "Empty Email",
			user: &models.User{
				Email:        "",
				PasswordHash: "hashed_password",
			},
			mock: func() {
				mock.ExpectQuery(`INSERT INTO users \(email, password_hash\) VALUES \(\$1, \$2\) RETURNING user_id`).
					WithArgs("", "hashed_password").
					WillReturnError(errors.New("empty email"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := repo.CreateUser(tt.user)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, int64(1), tt.user.ID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repo := &Repo{db: sqlxDB}

	now := time.Now().Format(time.RFC3339)

	tests := []struct {
		name    string
		email   string
		mock    func()
		want    *models.User
		wantErr bool
	}{
		{
			name:  "Success",
			email: "test@example.com",
			mock: func() {
				rows := sqlmock.NewRows([]string{"user_id", "email", "password_hash", "created_at"}).
					AddRow(1, "test@example.com", "hashed_password", now)
				mock.ExpectQuery("SELECT \\* FROM users WHERE email=\\$1").
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			want: &models.User{
				ID:           1,
				Email:        "test@example.com",
				PasswordHash: "hashed_password",
				CreatedAt:    now,
			},
			wantErr: false,
		},
		{
			name:  "Not Found",
			email: "notfound@example.com",
			mock: func() {
				mock.ExpectQuery("SELECT \\* FROM users WHERE email=\\$1").
					WithArgs("notfound@example.com").
					WillReturnError(errors.New("not found"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := repo.GetUserByEmail(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetUserByID(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repo := &Repo{db: sqlxDB}

	now := time.Now().Format(time.RFC3339)

	tests := []struct {
		name    string
		userID  int64
		mock    func()
		want    *models.User
		wantErr bool
	}{
		{
			name:   "Success",
			userID: 1,
			mock: func() {
				rows := sqlmock.NewRows([]string{"user_id", "email", "created_at"}).
					AddRow(1, "test@example.com", now)
				mock.ExpectQuery("SELECT user_id, email, created_at FROM users WHERE user_id = \\$1").
					WithArgs(1).
					WillReturnRows(rows)
			},
			want: &models.User{
				ID:        1,
				Email:     "test@example.com",
				CreatedAt: now,
			},
			wantErr: false,
		},
		{
			name:   "Not Found",
			userID: 999,
			mock: func() {
				mock.ExpectQuery("SELECT user_id, email, created_at FROM users WHERE user_id = \\$1").
					WithArgs(999).
					WillReturnError(errors.New("not found"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := repo.GetUserByID(tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdateUser(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repo := &Repo{db: sqlxDB}

	tests := []struct {
		name    string
		user    *models.User
		mock    func()
		wantErr bool
	}{
		{
			name: "Success",
			user: &models.User{
				ID:           1,
				Email:        "new@example.com",
				PasswordHash: "new_hash",
			},
			mock: func() {
				mock.ExpectExec("UPDATE users SET email = \\$1, password_hash = \\$2 WHERE user_id = \\$3").
					WithArgs("new@example.com", "new_hash", int64(1)).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "User Not Found",
			user: &models.User{
				ID:           999,
				Email:        "notfound@example.com",
				PasswordHash: "hash",
			},
			mock: func() {
				mock.ExpectExec("UPDATE users SET email = \\$1, password_hash = \\$2 WHERE user_id = \\$3").
					WithArgs("notfound@example.com", "hash", int64(999)).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := repo.UpdateUser(tt.user)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpsertMovie(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repo := &Repo{db: sqlxDB}

	movie := &models.Movie{
		ID:              1,
		Title:           "Test Movie",
		Year:            2023,
		PosterURL:       "http://example.com/poster.jpg",
		Description:     "Test description",
		RatingKinopoisk: 7.5,
	}

	tests := []struct {
		name    string
		movie   *models.Movie
		mock    func()
		wantErr bool
	}{
		{
			name:  "Success",
			movie: movie,
			mock: func() {
				mock.ExpectExec(`INSERT INTO movies.*`).
					WithArgs(
						movie.ID,
						movie.Title,
						movie.Year,
						movie.PosterURL,
						movie.Description,
						movie.RatingKinopoisk,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name:  "Empty Title",
			movie: &models.Movie{Title: ""},
			mock: func() {
				mock.ExpectExec(`INSERT INTO movies.*`).
					WillReturnError(errors.New("empty title"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := repo.UpsertMovie(tt.movie)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetMovieByID(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repo := &Repo{db: sqlxDB}

	now := time.Now()
	movie := &models.Movie{
		ID:              1,
		Title:           "Test Movie",
		Year:            2023,
		PosterURL:       "http://example.com/poster.jpg",
		Description:     "Test description",
		RatingKinopoisk: 7.5,
		LastSync:        now,
	}

	tests := []struct {
		name    string
		movieID int64
		mock    func()
		want    *models.Movie
		wantErr bool
	}{
		{
			name:    "Success",
			movieID: 1,
			mock: func() {
				rows := sqlmock.NewRows([]string{
					"movie_id", "title", "year", "poster_url",
					"description", "rating_kinopoisk", "last_sync",
				}).AddRow(
					movie.ID, movie.Title, movie.Year, movie.PosterURL,
					movie.Description, movie.RatingKinopoisk, movie.LastSync,
				)
				mock.ExpectQuery("SELECT \\* FROM movies WHERE movie_id=\\$1").
					WithArgs(1).
					WillReturnRows(rows)
			},
			want:    movie,
			wantErr: false,
		},
		{
			name:    "Not Found",
			movieID: 999,
			mock: func() {
				mock.ExpectQuery("SELECT \\* FROM movies WHERE movie_id=\\$1").
					WithArgs(999).
					WillReturnError(errors.New("not found"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := repo.GetMovieByID(tt.movieID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSearchMovies(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repo := &Repo{db: sqlxDB}

	movies := []models.Movie{
		{
			ID:              1,
			Title:           "Test Movie 1",
			Year:            2023,
			PosterURL:       "http://example.com/poster1.jpg",
			Description:     "Test description 1",
			RatingKinopoisk: 7.5,
		},
		{
			ID:              2,
			Title:           "Test Movie 2",
			Year:            2022,
			PosterURL:       "http://example.com/poster2.jpg",
			Description:     "Test description 2",
			RatingKinopoisk: 8.0,
		},
	}

	tests := []struct {
		name    string
		query   string
		mock    func()
		want    []models.Movie
		wantErr bool
	}{
		{
			name:  "Success",
			query: "test",
			mock: func() {
				rows := sqlmock.NewRows([]string{
					"movie_id", "title", "year", "poster_url",
					"description", "rating_kinopoisk", "last_sync",
				})
				for _, m := range movies {
					rows.AddRow(
						m.ID, m.Title, m.Year, m.PosterURL,
						m.Description, m.RatingKinopoisk, m.LastSync,
					)
				}
				mock.ExpectQuery("SELECT \\* FROM movies WHERE title ILIKE \\$1").
					WithArgs("%test%").
					WillReturnRows(rows)
			},
			want:    movies,
			wantErr: false,
		},
		{
			name:  "No Results",
			query: "nonexistent",
			mock: func() {
				rows := sqlmock.NewRows([]string{
					"movie_id", "title", "year", "poster_url",
					"description", "rating_kinopoisk", "last_sync",
				})
				mock.ExpectQuery("SELECT \\* FROM movies WHERE title ILIKE \\$1").
					WithArgs("%nonexistent%").
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := repo.SearchMovies(tt.query)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestListMovies(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repo := &Repo{db: sqlxDB}

	movies := []models.Movie{
		{
			ID:              1,
			Title:           "Test Movie 1",
			Year:            2023,
			PosterURL:       "http://example.com/poster1.jpg",
			Description:     "Test description 1",
			RatingKinopoisk: 7.5,
		},
		{
			ID:              2,
			Title:           "Test Movie 2",
			Year:            2022,
			PosterURL:       "http://example.com/poster2.jpg",
			Description:     "Test description 2",
			RatingKinopoisk: 8.0,
		},
	}

	tests := []struct {
		name    string
		offset  int
		limit   int
		mock    func()
		want    []models.Movie
		wantErr bool
	}{
		{
			name:   "Success",
			offset: 0,
			limit:  10,
			mock: func() {
				rows := sqlmock.NewRows([]string{
					"movie_id", "title", "year", "poster_url",
					"description", "rating_kinopoisk",
				})
				for _, m := range movies {
					rows.AddRow(
						m.ID, m.Title, m.Year, m.PosterURL,
						m.Description, m.RatingKinopoisk,
					)
				}
				mock.ExpectQuery(`SELECT movie_id, title, year, poster_url, description, rating_kinopoisk`).
					WithArgs(10, 0).
					WillReturnRows(rows)
			},
			want:    movies,
			wantErr: false,
		},
		{
			name:   "Empty Result",
			offset: 100,
			limit:  10,
			mock: func() {
				rows := sqlmock.NewRows([]string{
					"movie_id", "title", "year", "poster_url",
					"description", "rating_kinopoisk",
				})
				mock.ExpectQuery(`SELECT movie_id, title, year, poster_url, description, rating_kinopoisk`).
					WithArgs(10, 100).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := repo.ListMovies(tt.offset, tt.limit)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestListPopularMovies(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repo := &Repo{db: sqlxDB}

	movies := []models.Movie{
		{
			ID:              1,
			Title:           "Top Movie",
			Year:            2023,
			PosterURL:       "http://example.com/poster1.jpg",
			RatingKinopoisk: 9.5,
		},
		{
			ID:              2,
			Title:           "Second Movie",
			Year:            2022,
			PosterURL:       "http://example.com/poster2.jpg",
			RatingKinopoisk: 8.0,
		},
	}

	tests := []struct {
		name    string
		limit   int
		mock    func()
		want    []models.Movie
		wantErr bool
	}{
		{
			name:  "Success",
			limit: 2,
			mock: func() {
				rows := sqlmock.NewRows([]string{
					"movie_id", "title", "year", "poster_url", "rating_kinopoisk",
				})
				for _, m := range movies {
					rows.AddRow(
						m.ID, m.Title, m.Year, m.PosterURL, m.RatingKinopoisk,
					)
				}
				mock.ExpectQuery(`SELECT movie_id, title, year, poster_url, rating_kinopoisk`).
					WithArgs(2).
					WillReturnRows(rows)
			},
			want:    movies,
			wantErr: false,
		},
		{
			name:  "Zero Limit",
			limit: 0,
			mock: func() {
				rows := sqlmock.NewRows([]string{
					"movie_id", "title", "year", "poster_url", "rating_kinopoisk",
				})
				mock.ExpectQuery(`SELECT movie_id, title, year, poster_url, rating_kinopoisk`).
					WithArgs(0).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := repo.ListPopularMovies(tt.limit)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAddToWatchlist(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repo := &Repo{db: sqlxDB}

	tests := []struct {
		name    string
		item    *models.WatchlistItem
		mock    func()
		wantErr bool
	}{
		{
			name: "Success",
			item: &models.WatchlistItem{
				UserID:  1,
				MovieID: 1,
			},
			mock: func() {
				mock.ExpectExec(`INSERT INTO watchlist \(user_id, movie_id\) VALUES \(\$1,\$2\) ON CONFLICT DO NOTHING`).
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "Duplicate",
			item: &models.WatchlistItem{
				UserID:  1,
				MovieID: 1,
			},
			mock: func() {
				mock.ExpectExec(`INSERT INTO watchlist \(user_id, movie_id\) VALUES \(\$1,\$2\) ON CONFLICT DO NOTHING`).
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(0, 0)) // No rows affected
			},
			wantErr: false, // No error expected for duplicate
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := repo.AddToWatchlist(tt.item)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetWatchlist(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repo := &Repo{db: sqlxDB}

	now := time.Now().Format(time.RFC3339)
	items := []models.WatchlistItem{
		{
			MovieID: 1,
			AddedAt: now,
			Title:   "Movie 1",
			Poster:  "http://example.com/poster1.jpg",
		},
		{
			MovieID: 2,
			AddedAt: now,
			Title:   "Movie 2",
			Poster:  "http://example.com/poster2.jpg",
		},
	}

	tests := []struct {
		name    string
		userID  int64
		mock    func()
		want    []models.WatchlistItem
		wantErr bool
	}{
		{
			name:   "Success",
			userID: 1,
			mock: func() {
				rows := sqlmock.NewRows([]string{"movie_id", "added_at", "title", "poster_url"})
				for _, item := range items {
					rows.AddRow(item.MovieID, item.AddedAt, item.Title, item.Poster)
				}
				mock.ExpectQuery(`SELECT w.movie_id, w.added_at, m.title, m.poster_url`).
					WithArgs(1).
					WillReturnRows(rows)
			},
			want:    items,
			wantErr: false,
		},
		{
			name:   "Empty Watchlist",
			userID: 2,
			mock: func() {
				rows := sqlmock.NewRows([]string{"movie_id", "added_at", "title", "poster_url"})
				mock.ExpectQuery(`SELECT w.movie_id, w.added_at, m.title, m.poster_url`).
					WithArgs(2).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := repo.GetWatchlist(tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRemoveFromWatchlist(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repo := &Repo{db: sqlxDB}

	tests := []struct {
		name    string
		userID  int64
		movieID int64
		mock    func()
		wantErr bool
	}{
		{
			name:    "Success",
			userID:  1,
			movieID: 1,
			mock: func() {
				mock.ExpectExec(`DELETE FROM watchlist WHERE user_id=\$1 AND movie_id=\$2`).
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name:    "Not Found",
			userID:  1,
			movieID: 999,
			mock: func() {
				mock.ExpectExec(`DELETE FROM watchlist WHERE user_id=\$1 AND movie_id=\$2`).
					WithArgs(1, 999).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := repo.RemoveFromWatchlist(tt.userID, tt.movieID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpsertRating(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repo := &Repo{db: sqlxDB}

	tests := []struct {
		name    string
		item    *models.RatingItem
		mock    func()
		wantErr bool
	}{
		{
			name: "Success - New Rating",
			item: &models.RatingItem{
				UserID:  1,
				MovieID: 1,
				Rating:  8,
			},
			mock: func() {
				mock.ExpectExec(`INSERT INTO ratings \(user_id, movie_id, rating\) VALUES \(\$1,\$2,\$3\)`).
					WithArgs(1, 1, 8).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "Success - Update Rating",
			item: &models.RatingItem{
				UserID:  1,
				MovieID: 1,
				Rating:  9,
			},
			mock: func() {
				mock.ExpectExec(`INSERT INTO ratings \(user_id, movie_id, rating\) VALUES \(\$1,\$2,\$3\)`).
					WithArgs(1, 1, 9).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "Invalid Rating",
			item: &models.RatingItem{
				UserID:  1,
				MovieID: 1,
				Rating:  11, // Invalid rating
			},
			mock: func() {
				mock.ExpectExec(`INSERT INTO ratings \(user_id, movie_id, rating\) VALUES \(\$1,\$2,\$3\)`).
					WithArgs(1, 1, 11).
					WillReturnError(errors.New("invalid rating"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := repo.UpsertRating(tt.item)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetRatings(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repo := &Repo{db: sqlxDB}

	now := time.Now().Format(time.RFC3339)
	items := []models.RatingItem{
		{
			MovieID: 1,
			Rating:  8,
			RatedAt: now,
		},
		{
			MovieID: 2,
			Rating:  7,
			RatedAt: now,
		},
	}

	tests := []struct {
		name    string
		userID  int64
		mock    func()
		want    []models.RatingItem
		wantErr bool
	}{
		{
			name:   "Success",
			userID: 1,
			mock: func() {
				rows := sqlmock.NewRows([]string{"movie_id", "rating", "rated_at"})
				for _, item := range items {
					rows.AddRow(item.MovieID, item.Rating, item.RatedAt)
				}
				mock.ExpectQuery(`SELECT movie_id, rating, rated_at FROM ratings WHERE user_id=\$1`).
					WithArgs(1).
					WillReturnRows(rows)
			},
			want:    items,
			wantErr: false,
		},
		{
			name:   "No Ratings",
			userID: 2,
			mock: func() {
				rows := sqlmock.NewRows([]string{"movie_id", "rating", "rated_at"})
				mock.ExpectQuery(`SELECT movie_id, rating, rated_at FROM ratings WHERE user_id=\$1`).
					WithArgs(2).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := repo.GetRatings(tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDeleteRating(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repo := &Repo{db: sqlxDB}

	tests := []struct {
		name    string
		userID  int64
		movieID int64
		mock    func()
		wantErr bool
	}{
		{
			name:    "Success",
			userID:  1,
			movieID: 1,
			mock: func() {
				mock.ExpectExec(`DELETE FROM ratings WHERE user_id = \$1 AND movie_id = \$2`).
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name:    "Not Found",
			userID:  1,
			movieID: 999,
			mock: func() {
				mock.ExpectExec(`DELETE FROM ratings WHERE user_id = \$1 AND movie_id = \$2`).
					WithArgs(1, 999).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := repo.DeleteRating(tt.userID, tt.movieID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
