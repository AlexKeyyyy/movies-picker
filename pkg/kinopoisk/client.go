package kinopoisk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

// GetPopularAll получает одну страницу популярного списка
func (c *Client) GetPopularAll(page int) ([]Film, int, error) {
	url := fmt.Sprintf("%s/films/collections?type=TOP_POPULAR_ALL&page=%d", c.baseURL, page)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("X-API-KEY", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	dec := json.NewDecoder(resp.Body)
	dec.UseNumber()

	var cr CollectionsResponse
	if err := dec.Decode(&cr); err != nil {
		return nil, 0, err
	}
	return cr.Items, cr.TotalPages, nil
}

// SearchByKeyword получает одну страницу по ключевому слову
func (c *Client) SearchByKeyword(keyword string, page int) ([]Film, int, error) {
	url := fmt.Sprintf("%s/films?type=ALL&ratingFrom=0&ratingTo=10&yearFrom=1000&yearTo=3000&keyword=%s&page=%d",
		c.baseURL, url.QueryEscape(keyword), page)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("X-API-KEY", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	dec := json.NewDecoder(resp.Body)
	dec.UseNumber()

	var cr CollectionsResponse
	if err := dec.Decode(&cr); err != nil {
		return nil, 0, err
	}
	return cr.Items, cr.TotalPages, nil
}
