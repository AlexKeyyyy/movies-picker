package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

var (
	baseURL string
	token   string
	userID  int64
)

func TestMain(m *testing.M) {
	time.Sleep(8 * time.Second)

	baseURL = os.Getenv("API_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	email := fmt.Sprintf("int_%d@example.com", time.Now().UnixNano())
	credentials, _ := json.Marshal(map[string]string{"email": email, "password": "pass123"})

	resp, err := http.Post(baseURL+"/auth/register", "application/json", bytes.NewReader(credentials))
	if err != nil {
		panic("register error: " + err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		panic(fmt.Sprintf("register status %d", resp.StatusCode))
	}

	var regResp struct {
		UserID int64 `json:"user_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&regResp); err != nil {
		panic("decode register: " + err.Error())
	}
	userID = regResp.UserID

	resp2, err := http.Post(baseURL+"/auth/login", "application/json", bytes.NewReader(credentials))
	if err != nil {
		panic("login error: " + err.Error())
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("login status %d", resp2.StatusCode))
	}

	var loginResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp2.Body).Decode(&loginResp); err != nil {
		panic("decode login: " + err.Error())
	}
	token = loginResp.AccessToken

	os.Exit(m.Run())
}
