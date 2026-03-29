package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// INT-01: Успешная регистрация нового пользователя
func TestINT01_Register_Success(t *testing.T) {
    uniqueEmail := fmt.Sprintf("new_%d@example.com", time.Now().UnixNano())
    body, _ := json.Marshal(map[string]string{
        "email":    uniqueEmail,
        "password": "NewPass1!",
    })

    resp, err := http.Post(baseURL+"/auth/register", "application/json", bytes.NewReader(body))
    require.NoError(t, err)
    defer resp.Body.Close()

    assert.Equal(t, http.StatusCreated, resp.StatusCode)

    var result struct {
        UserID int64  `json:"user_id"`
        Email  string `json:"email"`
    }
    require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
    assert.NotZero(t, result.UserID, "user_id должен быть ненулевым")
    assert.Equal(t, uniqueEmail, result.Email)
}

// INT-02: Успешный логин зарегистрированного пользователя
func TestINT02_Login_Success(t *testing.T) {
    body, _ := json.Marshal(map[string]string{
        "email":    testEmail,
        "password": testPassword,
    })

    resp, err := http.Post(baseURL+"/auth/login", "application/json", bytes.NewReader(body))
    require.NoError(t, err)
    defer resp.Body.Close()

    assert.Equal(t, http.StatusOK, resp.StatusCode)

    var result struct {
        AccessToken string `json:"access_token"`
        TokenType   string `json:"token_type"`
        ExpiresIn   int    `json:"expires_in"`
    }
    require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
    assert.NotEmpty(t, result.AccessToken, "access_token должен быть непустым")
    assert.Equal(t, "Bearer", result.TokenType)
    assert.Greater(t, result.ExpiresIn, 0)
}

// INT-03: Попытка регистрации с уже существующим email (негативный)
func TestINT03_Register_DuplicateEmail_Negative(t *testing.T) {
    body, _ := json.Marshal(map[string]string{
        "email":    testEmail,
        "password": "AnotherPass1!",
    })

    resp, err := http.Post(baseURL+"/auth/register", "application/json", bytes.NewReader(body))
    require.NoError(t, err)
    defer resp.Body.Close()

    // Должен вернуть 409 Conflict — email уже занят
    assert.Equal(t, http.StatusConflict, resp.StatusCode,
        "повторная регистрация с тем же email должна вернуть 409")
}

// INT-04: Логин с несуществующим email (негативный)
func TestINT04_Login_NonExistentEmail_Negative(t *testing.T) {
    body, _ := json.Marshal(map[string]string{
        "email":    "ghost_nobody@example.com",
        "password": "SomePass1!",
    })

    resp, err := http.Post(baseURL+"/auth/login", "application/json", bytes.NewReader(body))
    require.NoError(t, err)
    defer resp.Body.Close()

    assert.Equal(t, http.StatusUnauthorized, resp.StatusCode,
        "логин с несуществующим email должен вернуть 401")
}

// INT-05: Логин с неверным паролем (негативный)
func TestINT05_Login_WrongPassword_Negative(t *testing.T) {
    body, _ := json.Marshal(map[string]string{
        "email":    testEmail,
        "password": "completely_wrong_password",
    })

    resp, err := http.Post(baseURL+"/auth/login", "application/json", bytes.NewReader(body))
    require.NoError(t, err)
    defer resp.Body.Close()

    assert.Equal(t, http.StatusUnauthorized, resp.StatusCode,
        "логин с неверным паролем должен вернуть 401")
}

// INT-06: Обращение к защищённому ресурсу без JWT (негативный)
func TestINT06_ProtectedEndpoint_NoToken_Negative(t *testing.T) {
    req, _ := http.NewRequest(http.MethodGet, baseURL+"/users/me", nil)
    resp, err := http.DefaultClient.Do(req)
    require.NoError(t, err)
    defer resp.Body.Close()

    assert.Equal(t, http.StatusUnauthorized, resp.StatusCode,
        "доступ без токена должен вернуть 401")
}