package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// INT-07: Получение профиля с валидным JWT
func TestINT07_GetProfile_WithValidJWT(t *testing.T) {
    req := authRequest(t, http.MethodGet, baseURL+"/users/me", nil)
    resp := doRequest(t, req)
    defer resp.Body.Close()

    require.Equal(t, http.StatusOK, resp.StatusCode)

    var profile struct {
        UserID int64  `json:"user_id"`
        Email  string `json:"email"`
    }
    require.NoError(t, json.NewDecoder(resp.Body).Decode(&profile))
    assert.Equal(t, userID, profile.UserID,
        "профиль должен принадлежать авторизованному пользователю")
    assert.Equal(t, testEmail, profile.Email)
}

// INT-08: Обновление email профиля
func TestINT08_UpdateProfile_Email(t *testing.T) {
    newEmail := fmt.Sprintf("updated_%d@example.com", time.Now().UnixNano())

    req := authRequest(t, http.MethodPatch, baseURL+"/users/me", map[string]string{
        "email":            newEmail,
        "current_password": testPassword,
    })
    resp := doRequest(t, req)
    defer resp.Body.Close()

    require.Equal(t, http.StatusOK, resp.StatusCode,
        "обновление email должно вернуть 200")

    var updated struct {
        Email string `json:"email"`
    }
    require.NoError(t, json.NewDecoder(resp.Body).Decode(&updated))
    assert.Equal(t, newEmail, updated.Email,
        "в ответе должен быть новый email")

    // Восстанавливаем исходный email чтобы не сломать другие тесты
    restoreReq := authRequest(t, http.MethodPatch, baseURL+"/users/me", map[string]string{
        "email":            testEmail,
        "current_password": testPassword,
    })
    restoreResp := doRequest(t, restoreReq)
    restoreResp.Body.Close()
}

// INT-09: Обновление пароля профиля
func TestINT09_UpdateProfile_Password(t *testing.T) {
    newPassword := "NewIntPass2!"

    req := authRequest(t, http.MethodPatch, baseURL+"/users/me", map[string]string{
        "password":         newPassword,
        "current_password": testPassword,
    })
    resp := doRequest(t, req)
    defer resp.Body.Close()

    require.Equal(t, http.StatusOK, resp.StatusCode,
        "обновление пароля должно вернуть 200")

    // Проверяем что новый пароль работает для логина
    loginBody := map[string]string{"email": testEmail, "password": newPassword}
    loginReq := authRequest(t, http.MethodPost, baseURL+"/auth/login", loginBody)
    // Логин — публичный endpoint, убираем старый токен
    loginReq.Header.Del("Authorization")
    loginResp := doRequest(t, loginReq)
    defer loginResp.Body.Close()
    assert.Equal(t, http.StatusOK, loginResp.StatusCode,
        "логин с новым паролем должен работать")

    // Восстанавливаем исходный пароль
    restoreReq := authRequest(t, http.MethodPatch, baseURL+"/users/me", map[string]string{
        "password":         testPassword,
        "current_password": newPassword,
    })
    // Нужно обновить токен для этого запроса — пересоздаём токен с новым паролем
    var newTokenResp struct{ AccessToken string `json:"access_token"` }
    json.NewDecoder(loginResp.Body).Decode(&newTokenResp)
    restoreReq.Header.Set("Authorization", "Bearer "+newTokenResp.AccessToken)
    restoreResp := doRequest(t, restoreReq)
    restoreResp.Body.Close()
}