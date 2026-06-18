package pl

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	httpClient *http.Client
	baseUrl    string
	token      string
}

func NewClient(token string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 2 * time.Second,
		},
		baseUrl: "https://api.alwaysgoodgames.de/api/pl",
		token:   token,
	}
}

func (c *Client) GetLeagueData(ctx context.Context) (*PrimeLeagueResponse, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token is empty")
	}
	params := url.Values{}
	//TODO: get from config
	params.Add("season_id", "3220")
	params.Add("stage_id", "509")
	params.Add("group_id", "6883")
	reqUrl := fmt.Sprintf("%s/league_season_group_get?%s", c.baseUrl, params.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Api-Key", c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var payload PrimeLeagueResponse
	if err = json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return &payload, nil
}
