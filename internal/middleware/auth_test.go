package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AlexKeyyyy/movies-picker/internal/middleware"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	secret := "test-secret"
	middlewareFunc := middleware.JWT(secret)

	tests := []struct {
		name           string
		setupRequest   func() *http.Request
		expectedStatus int
		expectedUserID int64
	}{
		{
			name: "Valid token",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"user_id": float64(123),
					"exp":     time.Now().Add(time.Hour).Unix(),
				})
				tokenStr, _ := token.SignedString([]byte(secret))
				req.Header.Set("Authorization", "Bearer "+tokenStr)
				return req
			},
			expectedStatus: http.StatusOK,
			expectedUserID: 123,
		},
		{
			name: "Missing Authorization header",
			setupRequest: func() *http.Request {
				return httptest.NewRequest("GET", "/", nil)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Invalid token format (no Bearer)",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				req.Header.Set("Authorization", "InvalidTokenFormat")
				return req
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Expired token",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"user_id": float64(123),
					"exp":     time.Now().Add(-time.Hour).Unix(),
				})
				tokenStr, _ := token.SignedString([]byte(secret))
				req.Header.Set("Authorization", "Bearer "+tokenStr)
				return req
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Invalid signature",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"user_id": float64(123),
					"exp":     time.Now().Add(time.Hour).Unix(),
				})
				tokenStr, _ := token.SignedString([]byte("wrong-secret"))
				req.Header.Set("Authorization", "Bearer "+tokenStr)
				return req
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем тестовый обработчик, который проверит контекст
			handler := middlewareFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Проверяем, что UserID правильно установлен в контексте
				if tt.expectedUserID > 0 {
					userID := r.Context().Value(middleware.UserIDKey).(int64)
					assert.Equal(t, tt.expectedUserID, userID)
				}
				w.WriteHeader(http.StatusOK)
			}))

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, tt.setupRequest())

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
