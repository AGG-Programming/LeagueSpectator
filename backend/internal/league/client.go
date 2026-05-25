package league

import (
	"crypto/tls"
	"io"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
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
		baseURL: "https://127.0.0.1:2999/liveclientdata",
	}
}

func (c *Client) FetchAllGameData() ([]byte, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/allgamedata")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
