package youtube

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
		baseURL:    "https://www.googleapis.com/youtube/v3",
	}
}

type searchResponse struct {
	NextPageToken string `json:"nextPageToken"`
	Items         []struct {
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
	} `json:"items"`
}

// ReviewResult представляет один обзор из YouTube
type ReviewResult struct {
	VideoID      string `json:"video_id"`
	VideoURL     string `json:"video_url"` // ссылка на видео
	Title        string `json:"title"`
	ChannelTitle string `json:"channel_title"`
	ThumbnailURL string `json:"thumbnail_url"`
}

// SearchReviews ищет видеообзоры по названию фильма
func (c *Client) SearchReviews(keyword string, maxResultsPerPage int) ([]ReviewResult, error) {
	var result []ReviewResult
	pageToken := ""
	count := 0

	for {
		params := url.Values{}
		params.Set("part", "snippet")
		params.Set("q", fmt.Sprintf("%s обзор", keyword))
		params.Set("type", "video")
		params.Set("maxResults", fmt.Sprintf("%d", maxResultsPerPage))
		if pageToken != "" {
			params.Set("pageToken", pageToken)
		}

		u := fmt.Sprintf("%s/search?%s&key=%s", c.baseURL, params.Encode(), c.apiKey)
		resp, err := c.httpClient.Get(u)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("youtube API status: %d", resp.StatusCode)
		}

		var sr searchResponse
		if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
			return nil, err
		}

		for _, item := range sr.Items {
			vid := item.ID.VideoID
			result = append(result, ReviewResult{
				VideoID:      vid,
				VideoURL:     fmt.Sprintf("https://www.youtube.com/watch?v=%s", vid),
				Title:        item.Snippet.Title,
				ChannelTitle: item.Snippet.ChannelTitle,
				ThumbnailURL: item.Snippet.Thumbnails.High.URL,
			})
		}

		count += len(sr.Items)
		// остановимся, если страниц больше нет или набрали слишком много
		if sr.NextPageToken == "" || count >= maxResultsPerPage {
			break
		}
		pageToken = sr.NextPageToken
	}

	return result, nil
}
