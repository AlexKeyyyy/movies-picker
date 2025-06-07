package models

import "time"

type User struct {
	ID           int64  `db:"user_id" json:"user_id"`
	Email        string `db:"email" json:"email"`
	PasswordHash string `db:"password_hash"`
	CreatedAt    string `db:"created_at" json:"created_at"`
}

type Movie struct {
	ID          int64     `db:"movie_id"    json:"movie_id"`
	Title       string    `db:"title"       json:"title"`
	Year        int       `db:"year"        json:"year"`
	PosterURL   string    `db:"poster_url"  json:"poster_url"`
	Description string    `db:"description" json:"description"`
	LastSync    time.Time `db:"last_sync"   json:"last_sync"`
}

type WatchlistItem struct {
	UserID  int64  `db:"user_id" json:"-"`
	MovieID int64  `db:"movie_id" json:"movie_id"`
	AddedAt string `db:"added_at" json:"added_at"`
	Title   string `db:"title"      json:"title"`
	Poster  string `db:"poster_url" json:"poster_url"`
}

type RatingItem struct {
	UserID  int64  `db:"user_id" json:"-"`
	MovieID int64  `db:"movie_id" json:"movie_id"`
	Rating  int    `db:"rating" json:"rating"`
	RatedAt string `db:"rated_at" json:"rated_at"`
}

type ReviewItem struct {
	VideoID      string `json:"video_id"`
	Title        string `json:"title"`
	ChannelTitle string `json:"channel_title"`
	ThumbnailURL string `json:"thumbnail_url"`
}
