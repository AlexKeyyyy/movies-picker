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
    baseURL       string
    mockServerURL string
    token         string
    userID        int64
    testEmail     string
    testPassword  string
)

func TestMain(m *testing.M) {
    
    time.Sleep(5 * time.Second)

    baseURL = os.Getenv("API_URL")
    if baseURL == "" {
        baseURL = "http://localhost:8080"
    }

    mockServerURL = os.Getenv("MOCK_SERVER_URL")
    if mockServerURL == "" {
        mockServerURL = "http://localhost:9090"
    }

    testEmail = fmt.Sprintf("int_%d@example.com", time.Now().UnixNano())
    testPassword = "IntPass1!"

    credentials := map[string]string{"email": testEmail, "password": testPassword}
    body, _ := json.Marshal(credentials)

    resp, err := http.Post(baseURL+"/auth/register", "application/json", bytes.NewReader(body))
    if err != nil {
        panic("register error: " + err.Error())
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusCreated {
        panic(fmt.Sprintf("register expected 201, got %d", resp.StatusCode))
    }

    var regResp struct {
        UserID int64  `json:"user_id"`
        Email  string `json:"email"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&regResp); err != nil {
        panic("decode register: " + err.Error())
    }
    userID = regResp.UserID

    resp2, err := http.Post(baseURL+"/auth/login", "application/json", bytes.NewReader(body))
    if err != nil {
        panic("login error: " + err.Error())
    }
    defer resp2.Body.Close()
    if resp2.StatusCode != http.StatusOK {
        panic(fmt.Sprintf("login expected 200, got %d", resp2.StatusCode))
    }

    var loginResp struct {
        AccessToken string `json:"access_token"`
    }
    if err := json.NewDecoder(resp2.Body).Decode(&loginResp); err != nil {
        panic("decode login: " + err.Error())
    }
    token = loginResp.AccessToken
    if token == "" {
        panic("login returned empty token")
    }

    os.Exit(m.Run())
}

func authRequest(t *testing.T, method, url string, body interface{}) *http.Request {
    t.Helper()
    var payload []byte
    if body != nil {
        var err error
        payload, err = json.Marshal(body)
        if err != nil {
            t.Fatalf("marshal body: %v", err)
        }
    }
    req, err := http.NewRequest(method, url, bytes.NewReader(payload))
    if err != nil {
        t.Fatalf("new request: %v", err)
    }
    req.Header.Set("Authorization", "Bearer "+token)
    if body != nil {
        req.Header.Set("Content-Type", "application/json")
    }
    return req
}

func doRequest(t *testing.T, req *http.Request) *http.Response {
    t.Helper()
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        t.Fatalf("do request: %v", err)
    }
    return resp
}

func getFirstMovieID(t *testing.T) int {
    t.Helper()
    resp, err := http.Get(fmt.Sprintf("%s/movies?page=1&size=1", baseURL))
    if err != nil {
        t.Fatalf("get movies: %v", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        t.Fatalf("list movies expected 200, got %d", resp.StatusCode)
    }
    var arr []struct {
        MovieID int `json:"movie_id"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&arr); err != nil {
        t.Fatalf("decode movies: %v", err)
    }
    if len(arr) == 0 {
        t.Skip("no movies in DB — skipping test that requires a movie")
    }
    return arr[0].MovieID
}