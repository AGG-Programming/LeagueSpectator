package ddragon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	httpClient *http.Client
	baseUrl    string
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 2 * time.Second,
		},
		baseUrl: "http://ddragon.leagueoflegends.com",
	}
}

func (c *Client) GetLatestPatchVersion() (string, error) {
	resp, err := c.httpClient.Get(c.baseUrl + "/api/versions.json")
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

func (c *Client) GetChampions(version string) (map[string]string, error) {
	resp, err := c.httpClient.Get(c.baseUrl + "/cdn/" + version + "/data/en_US/champion.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var payload ChampionResponse
	if err = json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	champions := make(map[string]string)
	for _, champ := range payload.Data {
		champions[champ.Name] = fmt.Sprintf("%s/cdn/%s/img/champion/%s", c.baseUrl, version, champ.Image.Full)
	}
	return champions, nil
}

func (c *Client) GetRunes(version string) (map[int]string, error) {
	resp, err := c.httpClient.Get(c.baseUrl + "/cdn/" + version + "/data/en_US/runesReforged.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var payload RuneResponse
	if err = json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	runes := make(map[int]string)
	for _, perk := range payload {
		runes[perk.ID] = fmt.Sprintf("%s/cdn/img/%s", c.baseUrl, perk.Icon)

		for _, keystone := range perk.Slots[0].Runes {
			runes[keystone.ID] = fmt.Sprintf("%s/cdn/img/%s", c.baseUrl, keystone.Icon)
		}
	}
	return runes, nil
}

func (c *Client) GetItems(version string) (map[int]string, error) {
	resp, err := c.httpClient.Get(c.baseUrl + "/cdn/" + version + "/data/en_US/item.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var payload ItemResponse
	if err = json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	items := make(map[int]string)
	for idStr, item := range payload.Data {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, err
		}
		items[id] = fmt.Sprintf("%s/cdn/%s/img/item/%s", c.baseUrl, version, item.Image.Full)
	}
	return items, nil
}

func (c *Client) GetSpells(version string) (map[string]string, error) {
	resp, err := c.httpClient.Get(c.baseUrl + "/cdn/" + version + "/data/en_US/summoner.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var payload SpellResponse
	if err = json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	spells := make(map[string]string)
	for id, spell := range payload.Data {
		spells[id] = fmt.Sprintf("%s/cdn/%s/img/spell/%s", c.baseUrl, version, spell.Image.Full)
	}
	return spells, nil
}

func (c *Client) GetUlts(version string, champs []string) (map[string]string, error) {
	ults := make(map[string]string)

	for _, champ := range champs {
		// 1. Construct the URL for the individual champion's JSON data
		url := fmt.Sprintf("%s/cdn/%s/data/en_US/champion/%s.json", c.baseUrl, version, champ)

		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch data for %s: %w", champ, err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("failed to get valid response for %s, status: %d", champ, resp.StatusCode)
		}

		// 2. Parse the JSON response
		var data ChampionDataResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		resp.Body.Close() // Close the body as soon as we are done decoding
		if err != nil {
			return nil, fmt.Errorf("failed to decode JSON for %s: %w", champ, err)
		}

		// 3. Extract the ultimate image (index 3 in the spells array)
		champData, exists := data.Data[champ]
		if !exists || len(champData.Spells) < 4 {
			return nil, fmt.Errorf("unexpected data structure or missing ultimate for %s", champ)
		}

		ultImageName := champData.Spells[3].Image.Full

		// 4. Save the full image URL to the map
		ults[champ] = fmt.Sprintf("%s/cdn/%s/img/spell/%s", c.baseUrl, version, ultImageName)
	}

	return ults, nil
}
