package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/AlexKeyyyy/movies-picker/internal/handlers"
	"github.com/AlexKeyyyy/movies-picker/internal/models"
	"github.com/stretchr/testify/assert"
)

type mockAuthService struct{}

func (m *mockAuthService) Register(email, password string) (*models.User, error) {
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if password == "" {
		return nil, fmt.Errorf("password is required")
	}

	re := regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
	if !re.MatchString(email) {
		return nil, fmt.Errorf("invalid email format")
	}

	if email == "exists@example.com" || email == "error@example.com" {
		return nil, fmt.Errorf("user already exists")
	}

	return &models.User{ID: 1, Email: email}, nil
}

func (m *mockAuthService) Login(email, password string) (string, error) {
	if email == "" || password == "" {
		return "", errors.New("empty fields")
	}
	if email != "test@example.com" || password != "pass123" {
		return "", errors.New("invalid credentials")
	}
	return "mock-token", nil
}

func TestAuthHandler_Register(t *testing.T) {
	handler := handlers.NewAuthHandler(&mockAuthService{})

	t.Run("success", func(t *testing.T) {
		body := `{"email":"test@example.com","password":"pass123"}`
		req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		handler.Register(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var data map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&data)
		assert.Equal(t, "test@example.com", data["email"])
		assert.Equal(t, float64(1), data["user_id"])
	})

	t.Run("invalid payload", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/register", bytes.NewBufferString("invalid-json"))
		w := httptest.NewRecorder()
		handler.Register(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("service error", func(t *testing.T) {
		body := `{"email":"error@example.com","password":"pass123"}`
		req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(body))
		w := httptest.NewRecorder()
		handler.Register(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("user already exists", func(t *testing.T) {
		body := `{"email":"exists@example.com","password":"pass123"}`
		req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(body))
		w := httptest.NewRecorder()
		handler.Register(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("empty email", func(t *testing.T) {
		body := `{"email":"","password":"pass123"}`
		req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(body))
		w := httptest.NewRecorder()
		handler.Register(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("empty password", func(t *testing.T) {
		body := `{"email":"test@example.com","password":""}`
		req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(body))
		w := httptest.NewRecorder()
		handler.Register(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("invalid_email_format", func(t *testing.T) {
		body := `{"email":"not-an-email","password":"pass123"}`
		req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(body))
		w := httptest.NewRecorder()
		handler.Register(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestAuthHandler_Login(t *testing.T) {
	handler := handlers.NewAuthHandler(&mockAuthService{})

	t.Run("success", func(t *testing.T) {
		body := `{"email":"test@example.com","password":"pass123"}`
		req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(body))
		w := httptest.NewRecorder()
		handler.Login(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var data map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&data)
		assert.Equal(t, "mock-token", data["access_token"])
	})

	t.Run("invalid payload", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/login", bytes.NewBufferString("invalid-json"))
		w := httptest.NewRecorder()
		handler.Login(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		body := `{"email":"fail@example.com","password":"wrong"}`
		req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(body))
		w := httptest.NewRecorder()
		handler.Login(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("empty email", func(t *testing.T) {
		body := `{"email":"","password":"pass123"}`
		req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(body))
		w := httptest.NewRecorder()
		handler.Login(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("empty password", func(t *testing.T) {
		body := `{"email":"test@example.com","password":""}`
		req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(body))
		w := httptest.NewRecorder()
		handler.Login(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
