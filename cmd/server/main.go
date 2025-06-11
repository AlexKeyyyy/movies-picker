package main

import (
	"log"
	"net/http"

	"github.com/AlexKeyyyy/movies-picker/config"
	"github.com/AlexKeyyyy/movies-picker/internal/handlers"
	"github.com/AlexKeyyyy/movies-picker/internal/middleware"
	"github.com/AlexKeyyyy/movies-picker/internal/repository"
	"github.com/AlexKeyyyy/movies-picker/internal/service"
	"github.com/AlexKeyyyy/movies-picker/pkg/kinopoisk"
	"github.com/AlexKeyyyy/movies-picker/pkg/youtube"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.Load()
	repo, err := repository.NewRepo(cfg.DBUrl)
	if err != nil {
		log.Fatal(err)
	}
	kpClient := kinopoisk.NewClient(cfg.KinopoiskApiKey)
	ytClient := youtube.NewClient(cfg.YouTubeApiKey)
	svc := service.NewService(repo, kpClient, ytClient, cfg.JWTSecret)

	authH := handlers.NewAuthHandler(svc)
	moviesH := handlers.NewMoviesHandler(svc)
	watchH := handlers.NewWatchlistHandler(svc)
	rateH := handlers.NewRatingsHandler(svc)

	r := chi.NewRouter()
	// public
	r.Post("/auth/register", authH.Register)
	r.Post("/auth/login", authH.Login)
	r.Get("/movies", moviesH.ListMovies)
	r.Get("/movies/search", moviesH.SearchMovies)
	r.Get("/movies/{id}", moviesH.GetMovie)
	r.Get("/movies/{id}/reviews", moviesH.GetMovieReviews)

	// protected
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWT(cfg.JWTSecret))
		r.Get("/users/{userID}/watchlist", watchH.GetWatchlist)
		r.Post("/users/{userID}/watchlist", watchH.AddToWatchlist)
		r.Delete("/users/{userID}/watchlist/{movieID}", watchH.RemoveFromWatchlist)
		r.Get("/users/{userID}/ratings", rateH.GetRatings)
		r.Post("/users/{userID}/ratings", rateH.AddOrUpdateRating)
	})

	log.Printf("Server running on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
