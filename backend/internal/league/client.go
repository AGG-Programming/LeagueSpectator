package league

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
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

func (c *Client) FetchAllGameData() ([]byte, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/allgamedata")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	return io.ReadAll(resp.Body)
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

	// The first element is always the latest live patch
	return versions[0], nil
}
