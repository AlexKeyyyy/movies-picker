package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AlexKeyyyy/movies-picker/internal/middleware"
	"github.com/AlexKeyyyy/movies-picker/internal/service"
)

type UserHandler struct {
	svc *service.Service
}

func NewUserHandler(svc *service.Service) *UserHandler {
	return &UserHandler{svc: svc}
}

// GET /users/me
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)
	user, err := h.svc.GetProfile(userID)
	if err != nil {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// PATCH /users/me
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	user, err := h.svc.UpdateProfile(userID, req.Email, req.Password)
	if err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
