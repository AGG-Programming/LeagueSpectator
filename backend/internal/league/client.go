package league

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	DDragonURL string
}

func NewClient() *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &Client{
		httpClient: &http.Client{
			Transport: tr,
			Timeout:   2 * time.Second,
		},
		baseURL:    "https://127.0.0.1:2999/liveclientdata",
		DDragonURL: "http://ddragon.leagueoflegends.com",
	}
}

func (c *Client) FetchAllGameData() (GameResponse, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/allgamedata")
	if err != nil {
		return GameResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GameResponse{}, err
	}
	var gameData GameResponse
	if err = json.NewDecoder(resp.Body).Decode(&gameData); err != nil {
		return GameResponse{}, err
	}

	return gameData, nil
}

func (c *Client) GetLatestPatchVersion() (string, error) {
	resp, err := c.httpClient.Get(c.DDragonURL + "/api/versions.json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var versions []string
	if err = json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return "", err
	}

	if len(versions) == 0 {
		return "", fmt.Errorf("no versions found")
	}

	return versions[0], nil
}
