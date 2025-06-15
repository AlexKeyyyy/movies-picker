package youtube

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	apiKey := "test-api-key"
	client := NewClient(apiKey)

	if client.apiKey != apiKey {
		t.Errorf("Expected apiKey %s, got %s", apiKey, client.apiKey)
	}

	if client.baseURL != "https://www.googleapis.com/youtube/v3" {
		t.Errorf("Unexpected baseURL: %s", client.baseURL)
	}

	if client.httpClient.Timeout != 10*time.Second {
		t.Errorf("Expected timeout 10s, got %v", client.httpClient.Timeout)
	}
}

func TestSearchReviews_Success(t *testing.T) {
	// Создаем тестовый сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем параметры запроса
		if r.URL.Query().Get("q") != "test обзор" {
			t.Errorf("Expected query 'test обзор', got '%s'", r.URL.Query().Get("q"))
		}
		if r.URL.Query().Get("maxResults") != "5" {
			t.Errorf("Expected maxResults 5, got %s", r.URL.Query().Get("maxResults"))
		}
		if r.URL.Query().Get("key") != "test-api-key" {
			t.Error("Missing or invalid API key")
		}

		// Возвращаем тестовый ответ
		response := searchResponse{
			NextPageToken: "next-page-token",
			Items: []struct {
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
			}{
				{
					ID: struct {
						VideoID string `json:"videoId"`
					}{VideoID: "video123"},
					Snippet: struct {
						Title        string `json:"title"`
						ChannelTitle string `json:"channelTitle"`
						Thumbnails   struct {
							High struct {
								URL string `json:"url"`
							} `json:"high"`
						} `json:"thumbnails"`
					}{
						Title:        "Test Review",
						ChannelTitle: "Test Channel",
						Thumbnails: struct {
							High struct {
								URL string `json:"url"`
							} `json:"high"`
						}{
							High: struct {
								URL string `json:"url"`
							}{URL: "http://test.com/thumb.jpg"},
						},
					},
				},
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()

	// Создаем клиент с тестовым URL
	client := &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		apiKey:     "test-api-key",
		baseURL:    ts.URL, // Используем URL тестового сервера
	}

	// Вызываем тестируемый метод
	results, err := client.SearchReviews("test", 5)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Проверяем результаты
	if len(results) != 5 {
		t.Fatalf("Expected 5 result, got %d", len(results))
	}

	expected := ReviewResult{
		VideoID:      "video123",
		VideoURL:     "https://www.youtube.com/watch?v=video123",
		Title:        "Test Review",
		ChannelTitle: "Test Channel",
		ThumbnailURL: "http://test.com/thumb.jpg",
	}

	if results[0] != expected {
		t.Errorf("Expected result %+v, got %+v", expected, results[0])
	}
}

func TestSearchReviews_ErrorHandling(t *testing.T) {
	// Тест для проверки обработки ошибок
	testCases := []struct {
		name        string
		statusCode  int
		response    string
		expectedErr string
	}{
		{
			name:        "API Error",
			statusCode:  http.StatusBadRequest,
			expectedErr: "youtube API status: 400",
		},
		{
			name:        "Invalid JSON",
			statusCode:  http.StatusOK,
			response:    "{invalid}",
			expectedErr: "invalid character",
		},
		{
			name:        "Server Error",
			statusCode:  http.StatusInternalServerError,
			expectedErr: "youtube API status: 500",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				if tc.response != "" {
					w.Write([]byte(tc.response))
				}
			}))
			defer ts.Close()

			client := &Client{
				httpClient: &http.Client{Timeout: 10 * time.Second},
				apiKey:     "test-api-key",
				baseURL:    ts.URL,
			}

			_, err := client.SearchReviews("test", 5)
			if err == nil {
				t.Fatal("Expected error, got nil")
			}

			if tc.expectedErr != "" && err.Error()[:len(tc.expectedErr)] != tc.expectedErr {
				t.Errorf("Expected error '%s', got '%v'", tc.expectedErr, err)
			}
		})
	}
}
