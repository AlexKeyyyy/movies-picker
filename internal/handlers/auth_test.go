package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AlexKeyyyy/movies-picker/internal/models"
	"github.com/AlexKeyyyy/movies-picker/internal/service"
	"golang.org/x/crypto/bcrypt"
)

func TestNewAuthHandler(t *testing.T) {
	h := NewAuthHandler(newTestService(&mockRepo{}, &mockKPClient{}, &mockYTClient{}))
	if h == nil {
		t.Fatal("expected handler to be created")
	}
	if h.svc == nil {
		t.Fatal("expected service to be set")
	}
}

func TestAuthHandler_Register_Success(t *testing.T) {
	repo := &mockRepo{
		createUserFn: func(u *models.User) error {
			u.ID = 42
			u.CreatedAt = "2026-03-30T10:00:00Z"
			return nil
		},
	}
	h := NewAuthHandler(newTestService(repo, &mockKPClient{}, &mockYTClient{}))

	req := httptest.NewRequest(http.MethodPost, "/register", jsonBody(t, map[string]string{
		"email":    "user@example.com",
		"password": "secret123",
	}))
	rr := httptest.NewRecorder()

	h.Register(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusCreated, rr.Code, rr.Body.String())
	}

	got := decodeJSON[map[string]any](t, rr)
	if got["user_id"].(float64) != 42 {
		t.Fatalf("expected user_id 42, got %v", got["user_id"])
	}
	if got["email"] != "user@example.com" {
		t.Fatalf("expected email user@example.com, got %v", got["email"])
	}
}

func TestAuthHandler_Register_InvalidBody(t *testing.T) {
	h := NewAuthHandler(newTestService(&mockRepo{}, &mockKPClient{}, &mockYTClient{}))

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader("{"))
	rr := httptest.NewRecorder()

	h.Register(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestAuthHandler_Register_ServiceError(t *testing.T) {
	repo := &mockRepo{
		createUserFn: func(u *models.User) error {
			return errors.New("user already exists")
		},
	}
	h := NewAuthHandler(newTestService(repo, &mockKPClient{}, &mockYTClient{}))

	req := httptest.NewRequest(http.MethodPost, "/register", jsonBody(t, map[string]string{
		"email":    "user@example.com",
		"password": "secret123",
	}))
	rr := httptest.NewRecorder()

	h.Register(rr, req)

	if rr.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, rr.Code)
	}
}

func TestAuthHandler_Login_Success(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	repo := &mockRepo{
		getUserByEmailFn: func(email string) (*models.User, error) {
			return &models.User{ID: 7, Email: email, PasswordHash: string(hash)}, nil
		},
	}
	h := NewAuthHandler(newTestService(repo, &mockKPClient{}, &mockYTClient{}))

	req := httptest.NewRequest(http.MethodPost, "/login", jsonBody(t, map[string]string{
		"email":    "user@example.com",
		"password": "secret123",
	}))
	rr := httptest.NewRecorder()

	h.Login(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, rr.Code, rr.Body.String())
	}

	got := decodeJSON[map[string]any](t, rr)
	if got["access_token"] == "" {
		t.Fatal("expected non-empty access_token")
	}
	if got["token_type"] != "Bearer" {
		t.Fatalf("expected token_type Bearer, got %v", got["token_type"])
	}
}

func TestAuthHandler_Login_InvalidBody(t *testing.T) {
	h := NewAuthHandler(newTestService(&mockRepo{}, &mockKPClient{}, &mockYTClient{}))

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("{"))
	rr := httptest.NewRecorder()

	h.Login(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestAuthHandler_Login_WrongCredentials(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("another-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	repo := &mockRepo{
		getUserByEmailFn: func(email string) (*models.User, error) {
			return &models.User{ID: 7, Email: email, PasswordHash: string(hash)}, nil
		},
	}
	h := NewAuthHandler(newTestService(repo, &mockKPClient{}, &mockYTClient{}))

	req := httptest.NewRequest(http.MethodPost, "/login", jsonBody(t, map[string]string{
		"email":    "user@example.com",
		"password": "secret123",
	}))
	rr := httptest.NewRecorder()

	h.Login(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestAuthHandler_Register_StatusCreated(t *testing.T) {
	repo := &mockRepo{
		createUserFn: func(u *models.User) error {
			u.ID = 123
			u.CreatedAt = "2026-03-30"
			return nil
		},
	}

	h := NewAuthHandler(service.NewService(repo, nil, nil, "secret"))

	body := `{"email":"test@example.com","password":"123456"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	rr := httptest.NewRecorder()

	h.Register(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}
}

func TestAuthHandler_Login_ResponseFields(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	repo := &mockRepo{
		getUserByEmailFn: func(email string) (*models.User, error) {
			return &models.User{
				ID:           77,
				Email:        email,
				PasswordHash: string(hash),
			}, nil
		},
	}

	h := NewAuthHandler(service.NewService(repo, nil, nil, "secret"))

	body := `{"email":"test@example.com","password":"123456"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	rr := httptest.NewRecorder()

	h.Login(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["token_type"] != "Bearer" {
		t.Fatalf("expected token_type Bearer, got %v", resp["token_type"])
	}

	if resp["access_token"] == "" {
		t.Fatal("expected non-empty access_token")
	}

	if resp["expires_in"] != float64(3600) {
		t.Fatalf("expected expires_in 3600, got %v", resp["expires_in"])
	}
}
