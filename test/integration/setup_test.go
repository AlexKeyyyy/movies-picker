package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	baseURL string
	token   string
	userID  int64
)

func TestMain(m *testing.M) {
	// 1) Дадим время Docker Compose запустить postgres и api (с импортом фильмов)
	time.Sleep(8 * time.Second)

	// 2) Определяем baseURL из окружения
	baseURL = os.Getenv("API_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	// 3) Регистрируем нового пользователя
	regBody, _ := json.Marshal(map[string]string{"email": "int@int.com", "password": "pass123"})
	resp, err := http.Post(baseURL+"/auth/register", "application/json", bytes.NewReader(regBody))
	if err != nil {
		panic("register error: " + err.Error())
	}
	defer resp.Body.Close()
	var regResp struct {
		UserID int64  `json:"user_id"`
		Email  string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&regResp); err != nil {
		panic("decode register: " + err.Error())
	}
	assert.Equal(nil, http.StatusCreated, resp.StatusCode)
	userID = regResp.UserID

	// 4) Логинимся, чтобы получить токен
	resp2, err := http.Post(baseURL+"/auth/login", "application/json", bytes.NewReader(regBody))
	if err != nil {
		panic("login error: " + err.Error())
	}
	defer resp2.Body.Close()
	var loginResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp2.Body).Decode(&loginResp); err != nil {
		panic("decode login: " + err.Error())
	}
	assert.Equal(nil, http.StatusOK, resp2.StatusCode)
	token = loginResp.AccessToken

	// 5) Запускаем остальные тесты
	os.Exit(m.Run())
}
