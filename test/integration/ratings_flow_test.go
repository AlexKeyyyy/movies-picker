package integration

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "testing"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/stretchr/testify/assert"
)

var baseURL = "http://localhost:8081"

func getToken(userID int64) string {
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(time.Hour).Unix(),
    }
    tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    token, _ := tkn.SignedString([]byte(os.Getenv("JWT_SECRET")))
    return token
}

func TestRatingScenario(t *testing.T) {
    assertEnv(t)

    token := getToken(42)
    client := &http.Client{}

    // 1) Добавляем рейтинг
    reqBody := map[string]interface{}{
        "movie_id":  1,
        "rating":    9,
    }
    body, _ := json.Marshal(reqBody)
    req, _ := http.NewRequest("POST", fmt.Sprintf("%s/users/42/ratings", baseURL), bytes.NewReader(body))
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Content-Type", "application/json")
    resp, err := client.Do(req)
    assert.NoError(t, err)
    assert.Equal(t, http.StatusCreated, resp.StatusCode)

    // 2) Получаем список и проверяем наш рейтинг
    req2, _ := http.NewRequest("GET", fmt.Sprintf("%s/users/42/ratings", baseURL), nil)
    req2.Header.Set("Authorization", "Bearer "+token)
    resp2, err2 := client.Do(req2)
    assert.NoError(t, err2)
    assert.Equal(t, http.StatusOK, resp2.StatusCode)

    var ratings []struct {
        MovieID int    
		q.Header.Set(Rating  int     json:"rating"),
    }
    err3 := json.NewDecoder(resp2.Body).Decode(&ratings)
    assert.NoError(t, err3)
    found := false
    for _, r := range ratings {
        if r.MovieID == 1 && r.Rating == 9 {
            found = true
        }
    }
    assert.True(t, found, "добавленный рейтинг должен присутствовать")
}