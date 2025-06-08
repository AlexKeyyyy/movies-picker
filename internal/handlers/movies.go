package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/AlexKeyyyy/movies-picker/internal/service"
	"github.com/go-chi/chi/v5"
)

type MoviesHandler struct {
	svc *service.Service
}

func NewMoviesHandler(svc *service.Service) *MoviesHandler {
	return &MoviesHandler{svc: svc}
}

// GET /movies/search?q={query}
func (h *MoviesHandler) SearchMovies(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		http.Error(w, "missing query parameter `q`", http.StatusBadRequest)
		return
	}
	movies, err := h.svc.SearchMovies(q)
	if err != nil {
		http.Error(w, "search failed", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(movies)
}

// GET /movies/{id}
func (h *MoviesHandler) GetMovie(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid movie id", http.StatusBadRequest)
		return
	}
	movie, err := h.svc.GetMovie(id)
	if err != nil {
		http.Error(w, "movie not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(movie)
}

// GET /movies/{id}/reviews
func (h *MoviesHandler) GetMovieReviews(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid movie id", http.StatusBadRequest)
		return
	}
	reviews, err := h.svc.GetMovieReviews(id)
	if err != nil {
		http.Error(w, "cannot fetch reviews", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(reviews)
}
