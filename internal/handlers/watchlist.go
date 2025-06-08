package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/AlexKeyyyy/movies-picker/internal/middleware"
	"github.com/AlexKeyyyy/movies-picker/internal/service"
	"github.com/go-chi/chi/v5"
)

type WatchlistHandler struct {
	svc *service.Service
}

func NewWatchlistHandler(svc *service.Service) *WatchlistHandler {
	return &WatchlistHandler{svc: svc}
}

// GET /users/{userID}/watchlist
func (h *WatchlistHandler) GetWatchlist(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(middleware.UserIDKey).(int64)
	log.Println(uid)
	list, err := h.svc.GetWatchlist(uid)
	if err != nil {
		http.Error(w, "failed to get watchlist", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(list)
}

// POST /users/{userID}/watchlist
func (h *WatchlistHandler) AddToWatchlist(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(middleware.UserIDKey).(int64)
	var req struct {
		MovieID int64 `json:"movie_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	if err := h.svc.AddToWatchlist(uid, req.MovieID); err != nil {
		http.Error(w, "cannot add to watchlist", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":  uid,
		"movie_id": req.MovieID,
	})
}

// DELETE /users/{userID}/watchlist/{movieID}
func (h *WatchlistHandler) RemoveFromWatchlist(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(middleware.UserIDKey).(int64)
	midStr := chi.URLParam(r, "movieID")
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid movie id", http.StatusBadRequest)
		return
	}
	if err := h.svc.RemoveFromWatchlist(uid, mid); err != nil {
		http.Error(w, "cannot remove from watchlist", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
