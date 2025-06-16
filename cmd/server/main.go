package main

import (
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/AlexKeyyyy/movies-picker/config"
	"github.com/AlexKeyyyy/movies-picker/internal/handlers"
	"github.com/AlexKeyyyy/movies-picker/internal/middleware"
	"github.com/AlexKeyyyy/movies-picker/internal/repository"
	"github.com/AlexKeyyyy/movies-picker/internal/service"
	"github.com/AlexKeyyyy/movies-picker/pkg/kinopoisk"
	"github.com/AlexKeyyyy/movies-picker/pkg/youtube"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
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
	userH := handlers.NewUserHandler(svc)
	moviesH := handlers.NewMoviesHandler(svc)
	watchH := handlers.NewWatchlistHandler(svc)
	rateH := handlers.NewRatingsHandler(svc)

	r := chi.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Применяем CORS к маршрутизатору `r`
	r.Use(c.Handler)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// --- Public endpoints ---
	r.Post("/auth/register", authH.Register)
	r.Post("/auth/login", authH.Login)

	r.Get("/movies", moviesH.ListMovies)                   // список фильмов с пагинацией
	r.Get("/movies/search", moviesH.SearchMovies)          // поиск
	r.Get("/movies/{id}", moviesH.GetMovie)                // детали
	r.Get("/movies/{id}/reviews", moviesH.GetMovieReviews) // обзоры
	r.Get("/movies/popular", moviesH.ListPopular)          // топ-N популярных

	// --- Protected endpoints (JWT required) ---
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWT(cfg.JWTSecret))

		// профиль
		r.Get("/users/me", userH.GetProfile)
		r.Patch("/users/me", userH.UpdateProfile)

		// «Смотреть позже»
		r.Get("/users/{userID}/watchlist", watchH.GetWatchlist)
		r.Post("/users/{userID}/watchlist", watchH.AddToWatchlist)
		r.Delete("/users/{userID}/watchlist/{movieID}", watchH.RemoveFromWatchlist)

		// рейтинги
		r.Get("/users/{userID}/ratings", rateH.GetRatings)
		r.Post("/users/{userID}/ratings", rateH.AddOrUpdateRating)
		r.Delete("/users/{userID}/ratings/{movieID}", rateH.DeleteRating)
	})

	// --- OpenAPI спецификация ---
	fsSpec := http.StripPrefix("/docs/spec/", http.FileServer(http.Dir("./docs")))
	r.Handle("/docs/spec/*", fsSpec)

	r.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:"+cfg.Port+"/docs/spec/openapi.yml"),
	))

	log.Printf("Server running on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))

}
