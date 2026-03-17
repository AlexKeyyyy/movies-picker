package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

// ─── Kinopoisk types ──────────────────
type Film struct {
    KinopoiskID     int64   `json:"kinopoiskId"`
    NameRu          string  `json:"nameRu"`
    NameEn          string  `json:"nameEn"`
    Year            string  `json:"year"`
    PosterURL       string  `json:"posterUrl"`
    Description     string  `json:"description"`
    RatingKinopoisk float64 `json:"ratingKinopoisk"`
}

type KPResponse struct {
    Total      int    `json:"total"`
    TotalPages int    `json:"totalPages"`
    Items      []Film `json:"items"`
}

// ─── YouTube ────────────────────────────────────────────
type YTResponse struct {
    Items []YTItem `json:"items"`
}

type YTItem struct {
    ID struct {
        VideoID string `json:"videoId"`
    } `json:"id"`
    Snippet struct {
        Title        string `json:"title"`
        ChannelTitle string `json:"channelTitle"`
        Thumbnails   struct {
            High struct {
                URL string `json:"url"`
            } `json:"high"`
        } `json:"thumbnails"`
    } `json:"snippet"`
}

// ─── Mock данные ──────────────────────────────────────────────
var mockFilms = []Film{
    {
        KinopoiskID:     301,
        NameRu:          "Интерстеллар",
        Year:            "2014",
        PosterURL:       "http://mock/poster1.jpg",
        Description:     "Фантастика о межзвёздных путешествиях",
        RatingKinopoisk: 8.6,
    },
    {
        KinopoiskID:     302,
        NameRu:          "Начало",
        Year:            "2010",
        PosterURL:       "http://mock/poster2.jpg",
        Description:     "Фильм о снах внутри снов",
        RatingKinopoisk: 8.7,
    },
    {
        KinopoiskID:     303,
        NameRu:          "Тёмный рыцарь",
        Year:            "2008",
        PosterURL:       "http://mock/poster3.jpg",
        Description:     "История Бэтмена и Джокера",
        RatingKinopoisk: 9.0,
    },
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(v)
}

func main() {
    port := os.Getenv("MOCK_PORT")
    if port == "" {
        port = "9090"
    }

    mux := http.NewServeMux()

    // ── Kinopoisk: популярные коллекции ──────────────────────
    // Вызывается при GET /api/v2.2/films/collections
    mux.HandleFunc("/api/v2.2/films/collections", func(w http.ResponseWriter, r *http.Request) {
        log.Printf("[MOCK-KP] collections request: %s", r.URL.String())
        writeJSON(w, http.StatusOK, KPResponse{
            Total:      len(mockFilms),
            TotalPages: 1,
            Items:      mockFilms,
        })
    })

    // ── Kinopoisk: поиск по ключевому слову ──────────────────
    // Вызывается при GET /api/v2.2/films?keyword=...
    mux.HandleFunc("/api/v2.2/films", func(w http.ResponseWriter, r *http.Request) {
        keyword := strings.ToLower(r.URL.Query().Get("keyword"))
        log.Printf("[MOCK-KP] search request: keyword=%s", keyword)

        var results []Film
        for _, f := range mockFilms {
            if keyword == "" || strings.Contains(strings.ToLower(f.NameRu), keyword) {
                results = append(results, f)
            }
        }
        writeJSON(w, http.StatusOK, KPResponse{
            Total:      len(results),
            TotalPages: 1,
            Items:      results,
        })
    })

    // ── YouTube: поиск обзоров ────────────────────────────────
    // Вызывается при GET /youtube/v3/search?q=...
    mux.HandleFunc("/youtube/v3/search", func(w http.ResponseWriter, r *http.Request) {
        q := r.URL.Query().Get("q")
        log.Printf("[MOCK-YT] search request: q=%s", q)

        resp := YTResponse{}
        // Создаём 3 фиктивных видео
        for i := 1; i <= 3; i++ {
            var item YTItem
            item.ID.VideoID = fmt.Sprintf("mock_video_%d", i)
            item.Snippet.Title = fmt.Sprintf("Обзор: %s — часть %d", q, i)
            item.Snippet.ChannelTitle = fmt.Sprintf("MockChannel%d", i)
            item.Snippet.Thumbnails.High.URL = fmt.Sprintf("http://mock/thumb%d.jpg", i)
            resp.Items = append(resp.Items, item)
        }
        writeJSON(w, http.StatusOK, resp)
    })

    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("ok"))
    })

    log.Printf("[MOCK] Mock external services server starting on :%s", port)
    if err := http.ListenAndServe(":"+port, mux); err != nil {
        log.Fatalf("Mock server failed: %v", err)
    }
}