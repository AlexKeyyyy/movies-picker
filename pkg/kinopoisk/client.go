package kinopoisk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

func NewClient(apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		apiKey:     apiKey,
		baseURL:    "https://kinopoiskapiunofficial.tech/api/v2.2",
	}
}

type Film struct {
	KinopoiskID     int64       `json:"kinopoiskId"`
	NameRu          string      `json:"nameRu"`
	NameEn          string      `json:"nameEn"`
	Year            json.Number `json:"year"`
	PosterURL       string      `json:"posterUrl"`
	Description     string      `json:"description"`
	RatingKinopoisk float64     `json:"ratingKinopoisk"`
}

type CollectionsResponse struct {
	Total      int    `json:"total"`
	TotalPages int    `json:"totalPages"`
	Items      []Film `json:"items"`
}

func (c *Client) GetTop250Films(page int) ([]Film, error) {
	url := fmt.Sprintf("%s/films/collections?type=TOP_250_MOVIES&page=%d", c.baseURL, page)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-API-KEY", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	dec := json.NewDecoder(resp.Body)
	dec.UseNumber()

	var result CollectionsResponse
	if err := dec.Decode(&result); err != nil {
		return nil, err
	}
	return result.Items, nil
}
