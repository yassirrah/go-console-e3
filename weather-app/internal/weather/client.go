package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	APIKey string
	HTTP   *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		APIKey: apiKey,
		HTTP:   &http.Client{Timeout: 8 * time.Second},
	}
}

func (c *Client) GetByCity(city, units string) (Response, error) {
	var out Response
	if units == "" {
		units = "metric"
	}
	q := url.QueryEscape(city)
	u := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=%s",
		q, c.APIKey, units,
	)

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return out, err
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return out, fmt.Errorf("upstream error: %s", resp.Status)
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return out, err
	}
	return out, nil
}
