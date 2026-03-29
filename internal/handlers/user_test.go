package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AlexKeyyyy/movies-picker/internal/middleware"
	"github.com/AlexKeyyyy/movies-picker/internal/models"
	"github.com/AlexKeyyyy/movies-picker/internal/service"
)

func TestUserHandler_GetProfile_Success(t *testing.T) {
	repo := &mockRepo{
		getUserByIDFn: func(userID int64) (*models.User, error) {
			if userID != 123 {
				t.Fatalf("expected userID 123, got %d", userID)
			}
			return &models.User{ID: 123, Email: "user@example.com"}, nil
		},
	}
	h := NewUserHandler(newTestService(repo, &mockKPClient{}, &mockYTClient{}))

	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	req = req.WithContext(withUserID(req.Context(), 123))
	rr := httptest.NewRecorder()

	h.GetProfile(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	got := decodeJSON[models.User](t, rr)
	if got.ID != 123 || got.Email != "user@example.com" {
		t.Fatalf("unexpected response: %+v", got)
	}
}

func TestUserHandler_GetProfile_UnauthorizedNoUserIDInContext(t *testing.T) {
	h := NewUserHandler(service.NewService(&mockRepo{}, nil, nil, "secret"))

	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	rr := httptest.NewRecorder()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic, got none")
		}
	}()

	h.GetProfile(rr, req)
}

func TestUserHandler_UpdateProfile_Success(t *testing.T) {
	repo := &mockRepo{
		getUserByIDFn: func(userID int64) (*models.User, error) {
			return &models.User{ID: userID, Email: "old@example.com", PasswordHash: "old-hash"}, nil
		},
		updateUserFn: func(user *models.User) error {
			if user.Email != "new@example.com" {
				t.Fatalf("expected updated email, got %s", user.Email)
			}
			if user.PasswordHash == "old-hash" {
				t.Fatal("expected password hash to be changed")
			}
			return nil
		},
	}
	h := NewUserHandler(newTestService(repo, &mockKPClient{}, &mockYTClient{}))

	req := httptest.NewRequest(http.MethodPatch, "/users/me", jsonBody(t, map[string]string{
		"email":    "new@example.com",
		"password": "newpass123",
	}))
	req = req.WithContext(withUserID(req.Context(), 123))
	rr := httptest.NewRecorder()

	h.UpdateProfile(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, rr.Code, rr.Body.String())
	}
	got := decodeJSON[models.User](t, rr)
	if got.Email != "new@example.com" {
		t.Fatalf("expected updated email, got %s", got.Email)
	}
}

func TestUserHandler_UpdateProfile_InvalidBody(t *testing.T) {
	h := NewUserHandler(newTestService(&mockRepo{}, &mockKPClient{}, &mockYTClient{}))

	req := httptest.NewRequest(http.MethodPatch, "/users/me", strings.NewReader("{"))
	req = req.WithContext(withUserID(req.Context(), 123))
	rr := httptest.NewRecorder()

	h.UpdateProfile(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestUserHandler_UpdateProfile_ServiceError(t *testing.T) {
	repo := &mockRepo{
		getUserByIDFn: func(userID int64) (*models.User, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewUserHandler(newTestService(repo, &mockKPClient{}, &mockYTClient{}))

	req := httptest.NewRequest(http.MethodPatch, "/users/me", jsonBody(t, map[string]string{
		"email": "new@example.com",
	}))
	req = req.WithContext(withUserID(req.Context(), 123))
	rr := httptest.NewRecorder()

	h.UpdateProfile(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}

func TestNewUserHandler(t *testing.T) {
	h := NewUserHandler(service.NewService(&mockRepo{}, nil, nil, "secret"))
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
	if h.svc == nil {
		t.Fatal("expected non-nil service")
	}
}

func TestUserHandler_GetProfile_NotFound(t *testing.T) {
	repo := &mockRepo{
		getUserByIDFn: func(userID int64) (*models.User, error) {
			return nil, errors.New("not found")
		},
	}

	h := NewUserHandler(service.NewService(repo, nil, nil, "secret"))

	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, int64(7)))
	rr := httptest.NewRecorder()

	h.GetProfile(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestUserHandler_GetProfile_ContentType(t *testing.T) {
	repo := &mockRepo{
		getUserByIDFn: func(userID int64) (*models.User, error) {
			return &models.User{ID: userID, Email: "test@example.com"}, nil
		},
	}

	h := NewUserHandler(service.NewService(repo, nil, nil, "secret"))

	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, int64(7)))
	rr := httptest.NewRecorder()

	h.GetProfile(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", ct)
	}
}

func TestUserHandler_UpdateProfile_ContentType(t *testing.T) {
	repo := &mockRepo{
		getUserByIDFn: func(userID int64) (*models.User, error) {
			return &models.User{ID: userID, Email: "old@example.com"}, nil
		},
		updateUserFn: func(user *models.User) error {
			return nil
		},
	}

	h := NewUserHandler(service.NewService(repo, nil, nil, "secret"))

	body := `{"email":"new@example.com","password":"123456"}`
	req := httptest.NewRequest(http.MethodPatch, "/users/me", strings.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, int64(7)))
	rr := httptest.NewRecorder()

	h.UpdateProfile(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", ct)
	}
}
