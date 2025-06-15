package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/AlexKeyyyy/movies-picker/internal/middleware"
	"github.com/AlexKeyyyy/movies-picker/internal/models"
	"github.com/AlexKeyyyy/movies-picker/internal/service"
	"github.com/go-chi/chi/v5"
)

type RatingsHandler struct {
	svc *service.Service
}

func NewRatingsHandler(svc *service.Service) *RatingsHandler {
	return &RatingsHandler{svc: svc}
}

// GET /users/{userID}/ratings
func (h *RatingsHandler) GetRatings(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(middleware.UserIDKey).(int64)
	list, err := h.svc.GetRatings(uid)
	if err != nil {
		http.Error(w, "failed to get ratings", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(list)
}

// POST /users/{userID}/ratings
func (h *RatingsHandler) AddOrUpdateRating(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(middleware.UserIDKey).(int64)
	var req struct {
		MovieID int64 `json:"movie_id"`
		Rating  int   `json:"rating"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	item := &models.RatingItem{
		UserID:  uid,
		MovieID: req.MovieID,
		Rating:  req.Rating,
	}
	if err := h.svc.UpsertRating(item); err != nil {
		http.Error(w, "cannot set rating", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

// handlers/ratings.go (добавить DeleteRating)
func (h *RatingsHandler) DeleteRating(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	movieID, _ := strconv.ParseInt(chi.URLParam(r, "movieID"), 10, 64)
	if err := h.svc.DeleteRating(userID, movieID); err != nil {
		http.Error(w, "failed to delete rating", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
