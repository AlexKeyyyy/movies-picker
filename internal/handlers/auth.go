package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AlexKeyyyy/movies-picker/internal/service"
)

type AuthHandler struct {
	svc service.AuthIface
}

func NewAuthHandler(svc service.AuthIface) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	user, err := h.svc.Register(req.Email, req.Password)
	if err != nil {
		msg := err.Error()
		switch msg {
		case "email is required", "password is required", "invalid email format":
			http.Error(w, msg, http.StatusBadRequest)
		case "user already exists":
			http.Error(w, msg, http.StatusConflict)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
		"created": user.CreatedAt,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	token, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token": token,
		"token_type":   "Bearer",
		"expires_in":   3600,
	})
}
