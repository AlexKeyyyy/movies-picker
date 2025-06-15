package kinopoisk

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

	if client.baseURL != "https://kinopoiskapiunofficial.tech/api/v2.2" {
		t.Errorf("Unexpected baseURL: %s", client.baseURL)
	}

	if client.httpClient.Timeout != 10*time.Second {
		t.Errorf("Expected timeout 10s, got %v", client.httpClient.Timeout)
	}
}

func TestGetPopularAll_Success(t *testing.T) {
	// Создаем тестовый сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем заголовки
		if r.Header.Get("X-API-KEY") != "test-api-key" {
			t.Error("Missing or invalid API key header")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("Missing or invalid Content-Type header")
		}

		// Проверяем URL
		expectedPath := "/films/collections"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// Возвращаем тестовый ответ
		response := CollectionsResponse{
			Total:      100,
			TotalPages: 5,
			Items: []Film{
				{
					KinopoiskID:     123,
					NameRu:          "Тестовый фильм",
					Year:            json.Number("2020"),
					RatingKinopoisk: 8.5,
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
	films, totalPages, err := client.GetPopularAll(1)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Проверяем результаты
	if totalPages != 5 {
		t.Errorf("Expected totalPages 5, got %d", totalPages)
	}

	if len(films) != 1 {
		t.Fatalf("Expected 1 film, got %d", len(films))
	}

	if films[0].KinopoiskID != 123 {
		t.Errorf("Expected film ID 123, got %d", films[0].KinopoiskID)
	}
}

func TestGetPopularAll_ErrorHandling(t *testing.T) {
	// Тест для проверки обработки ошибок
	testCases := []struct {
		name         string
		statusCode   int
		expectedErr  string
		mockResponse string
	}{
		{
			name:        "Server Error",
			statusCode:  http.StatusInternalServerError,
			expectedErr: "unexpected status code: 500",
		},
		{
			name:        "Bad Request",
			statusCode:  http.StatusBadRequest,
			expectedErr: "unexpected status code: 400",
		},
		{
			name:         "Invalid JSON",
			statusCode:   http.StatusOK,
			expectedErr:  "invalid character",
			mockResponse: "{invalid json}",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				if tc.mockResponse != "" {
					w.Write([]byte(tc.mockResponse))
				}
			}))
			defer ts.Close()

			client := &Client{
				httpClient: &http.Client{Timeout: 10 * time.Second},
				apiKey:     "test-api-key",
				baseURL:    ts.URL,
			}

			_, _, err := client.GetPopularAll(1)
			if err == nil {
				t.Fatal("Expected error, got nil")
			}

			if tc.expectedErr != "" && err.Error()[:len(tc.expectedErr)] != tc.expectedErr {
				t.Errorf("Expected error '%s', got '%v'", tc.expectedErr, err)
			}
		})
	}
}

func TestSearchByKeyword_Success(t *testing.T) {
	// Создаем тестовый сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем параметры запроса
		keyword := r.URL.Query().Get("keyword")
		if keyword != "test" {
			t.Errorf("Expected keyword 'test', got '%s'", keyword)
		}

		page := r.URL.Query().Get("page")
		if page != "2" {
			t.Errorf("Expected page 2, got %s", page)
		}

		// Возвращаем тестовый ответ
		response := CollectionsResponse{
			Total:      50,
			TotalPages: 3,
			Items: []Film{
				{
					KinopoiskID:     456,
					NameRu:          "Поисковый фильм",
					Year:            json.Number("2019"),
					RatingKinopoisk: 7.8,
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
		baseURL:    ts.URL,
	}

	// Вызываем тестируемый метод
	films, totalPages, err := client.SearchByKeyword("test", 2)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Проверяем результаты
	if totalPages != 3 {
		t.Errorf("Expected totalPages 3, got %d", totalPages)
	}

	if len(films) != 1 {
		t.Fatalf("Expected 1 film, got %d", len(films))
	}

	if films[0].KinopoiskID != 456 {
		t.Errorf("Expected film ID 456, got %d", films[0].KinopoiskID)
	}
}

func TestSearchByKeyword_ErrorHandling(t *testing.T) {
	// Тест для проверки обработки ошибок
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	client := &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		apiKey:     "test-api-key",
		baseURL:    ts.URL,
	}

	_, _, err := client.SearchByKeyword("test", 1)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedErr := "unexpected status code: 404"
	if err.Error() != expectedErr {
		t.Errorf("Expected error '%s', got '%v'", expectedErr, err)
	}
}
