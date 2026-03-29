package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlexKeyyyy/movies-picker/internal/handlers"
	"github.com/AlexKeyyyy/movies-picker/internal/middleware"
	"github.com/AlexKeyyyy/movies-picker/internal/models"
	"github.com/stretchr/testify/assert"
)

type mockUserService struct{}

func (m *mockUserService) GetProfile(userID int64) (*models.User, error) {
	if userID == 1 || userID == 3 {
		return &models.User{ID: userID, Email: "user@example.com"}, nil
	}
	return nil, errors.New("not found")
}

func (m *mockUserService) UpdateProfile(userID int64, email, password string) (*models.User, error) {
	if userID != 1 {
		return nil, errors.New("not found")
	}
	if email == "" && password == "" {
		return &models.User{ID: userID, Email: "user@example.com"}, nil
	}
	if email == "" {
		return &models.User{ID: userID, Email: "user@example.com"}, nil
	}
	if password == "" {
		return &models.User{ID: userID, Email: email}, nil
	}
	return &models.User{ID: userID, Email: email}, nil
}

func requestWithUserID(method, path string, body []byte, userID int64) *http.Request {
	req := httptest.NewRequest(method, path, bytes.NewBuffer(body))
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
	return req.WithContext(ctx)
}

func TestUserHandler_GetProfile(t *testing.T) {
	handler := handlers.NewUserHandler(&mockUserService{})

	t.Run("success - user 1", func(t *testing.T) {
		req := requestWithUserID("GET", "/users/me", nil, 1)
		w := httptest.NewRecorder()
		handler.GetProfile(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var user models.User
		json.NewDecoder(resp.Body).Decode(&user)
		assert.Equal(t, int64(1), user.ID)
	})

	t.Run("success - user 3", func(t *testing.T) {
		req := requestWithUserID("GET", "/users/me", nil, 3)
		w := httptest.NewRecorder()
		handler.GetProfile(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("profile not found", func(t *testing.T) {
		req := requestWithUserID("GET", "/users/me", nil, 2)
		w := httptest.NewRecorder()
		handler.GetProfile(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("userID zero", func(t *testing.T) {
		req := requestWithUserID("GET", "/users/me", nil, 0)
		w := httptest.NewRecorder()
		handler.GetProfile(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestUserHandler_UpdateProfile(t *testing.T) {
	handler := handlers.NewUserHandler(&mockUserService{})

	t.Run("update email only", func(t *testing.T) {
		body := `{"email":"new@example.com"}`
		req := requestWithUserID("PATCH", "/users/me", []byte(body), 1)
		w := httptest.NewRecorder()
		handler.UpdateProfile(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("update password only", func(t *testing.T) {
		body := `{"password":"newpass"}`
		req := requestWithUserID("PATCH", "/users/me", []byte(body), 1)
		w := httptest.NewRecorder()
		handler.UpdateProfile(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("update email + password", func(t *testing.T) {
		body := `{"email":"combo@example.com","password":"combo123"}`
		req := requestWithUserID("PATCH", "/users/me", []byte(body), 1)
		w := httptest.NewRecorder()
		handler.UpdateProfile(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("invalid JSON payload", func(t *testing.T) {
		req := requestWithUserID("PATCH", "/users/me", []byte("not-json"), 1)
		w := httptest.NewRecorder()
		handler.UpdateProfile(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("update user not found", func(t *testing.T) {
		body := `{"email":"fail@example.com"}`
		req := requestWithUserID("PATCH", "/users/me", []byte(body), 2)
		w := httptest.NewRecorder()
		handler.UpdateProfile(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("empty body payload", func(t *testing.T) {
		req := requestWithUserID("PATCH", "/users/me", []byte(`{}`), 1)
		w := httptest.NewRecorder()
		handler.UpdateProfile(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("partial invalid fields ignored - empty email", func(t *testing.T) {
		body := `{"email":"","password":"newpass"}`
		req := requestWithUserID("PATCH", "/users/me", []byte(body), 1)
		w := httptest.NewRecorder()
		handler.UpdateProfile(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("partial invalid fields ignored - empty password", func(t *testing.T) {
		body := `{"email":"new@example.com","password":""}`
		req := requestWithUserID("PATCH", "/users/me", []byte(body), 1)
		w := httptest.NewRecorder()
		handler.UpdateProfile(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("empty_email_and_password", func(t *testing.T) {
		body := `{"email":"","password":""}`
		req := requestWithUserID("PATCH", "/users/me", []byte(body), 1)
		w := httptest.NewRecorder()
		handler.UpdateProfile(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var user models.User
		json.NewDecoder(resp.Body).Decode(&user)
		assert.Equal(t, int64(1), user.ID)
		assert.Equal(t, "user@example.com", user.Email)
	})
}
