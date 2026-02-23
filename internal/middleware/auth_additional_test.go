package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AlexKeyyyy/movies-picker/internal/middleware"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func signedToken(t *testing.T, secret string, claims jwt.MapClaims) string {
	t.Helper()
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := tkn.SignedString([]byte(secret))
	require.NoError(t, err)
	return s
}

func makeHandler(secret string) http.Handler {
	mw := middleware.JWT(secret)
	return mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
}

func TestJWTRejectsMissingAuthorization(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	makeHandler("s").ServeHTTP(rr, req)
	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestJWTRejectsNoBearerPrefix(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "token")
	makeHandler("s").ServeHTTP(rr, req)
	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestJWTRejectsEmptyBearerToken(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer ")
	makeHandler("s").ServeHTTP(rr, req)
	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestJWTRejectsExpiredToken(t *testing.T) {
	secret := "secret"
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	token := signedToken(t, secret, jwt.MapClaims{"user_id": float64(1), "exp": time.Now().Add(-time.Minute).Unix()})
	req.Header.Set("Authorization", "Bearer "+token)
	makeHandler(secret).ServeHTTP(rr, req)
	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestJWTRejectsWrongSignature(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	token := signedToken(t, "a", jwt.MapClaims{"user_id": float64(1), "exp": time.Now().Add(time.Minute).Unix()})
	req.Header.Set("Authorization", "Bearer "+token)
	makeHandler("b").ServeHTTP(rr, req)
	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestJWTAcceptsValidToken(t *testing.T) {
	secret := "secret"
	captured := int64(0)
	mw := middleware.JWT(secret)
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = r.Context().Value(middleware.UserIDKey).(int64)
		w.WriteHeader(http.StatusOK)
	}))

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	token := signedToken(t, secret, jwt.MapClaims{"user_id": float64(42), "exp": time.Now().Add(time.Minute).Unix()})
	req.Header.Set("Authorization", "Bearer "+token)
	h.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, int64(42), captured)
}

func TestJWTRejectsMalformedToken(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer malformed")
	makeHandler("s").ServeHTTP(rr, req)
	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestJWTHeaderCaseSensitivePrefix(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "bearer token")
	makeHandler("s").ServeHTTP(rr, req)
	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestJWTRejectsTokenWithoutUserID(t *testing.T) {
	secret := "secret"
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	token := signedToken(t, secret, jwt.MapClaims{"exp": time.Now().Add(time.Minute).Unix()})
	req.Header.Set("Authorization", "Bearer "+token)
	require.Panics(t, func() { makeHandler(secret).ServeHTTP(rr, req) })
}

func TestJWTRejectsTokenWithStringUserID(t *testing.T) {
	secret := "secret"
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	token := signedToken(t, secret, jwt.MapClaims{"user_id": "42", "exp": time.Now().Add(time.Minute).Unix()})
	req.Header.Set("Authorization", "Bearer "+token)
	require.Panics(t, func() { makeHandler(secret).ServeHTTP(rr, req) })
}
