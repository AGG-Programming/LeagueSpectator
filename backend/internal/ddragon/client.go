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

func (c *Client) GetChampions(version string) (map[string]string, []string, error) {
	resp, err := c.httpClient.Get(c.baseUrl + "/cdn/" + version + "/data/en_US/champion.json")
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	var payload ChampionResponse
	if err = json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, nil, err
	}

	champions := make(map[string]string)
	champIDs := make([]string, 0, len(payload.Data))
	for _, champ := range payload.Data {
		champions[champ.Name] = fmt.Sprintf("%s/cdn/%s/img/champion/%s", c.baseUrl, version, champ.Image.Full)
		champIDs = append(champIDs, champ.ID)
	}
	return champions, champIDs, nil
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
	evolvedSpells := []string{"S5", "S12", "SummonerSmiteAvatarOffensive", "SummonerSmiteAvatarDefensive", "SummonerSmiteAvatarUtility"}

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
	for _, spell := range evolvedSpells {
		spells[spell] = fmt.Sprintf("./assets/images/summs/%s.png", spell)
	}

	return spells, nil
}

func (c *Client) GetUlts(version string, champIDs []string) (map[string]string, error) {
	ults := make(map[string]string)

	for _, champID := range champIDs {
		url := fmt.Sprintf("%s/cdn/%s/data/en_US/champion/%s.json", c.baseUrl, version, champID)

		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch data for %s: %w", champID, err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("failed to get valid response for %s, status: %d", champID, resp.StatusCode)
		}

		var data ChampionDataResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to decode JSON for %s: %w", champID, err)
		}

		champData, exists := data.Data[champID]
		if !exists || len(champData.Spells) < 4 {
			return nil, fmt.Errorf("unexpected data structure or missing ultimate for %s", champID)
		}

		ultImageName := champData.Spells[3].Image.Full

		ults[champData.Name] = fmt.Sprintf("%s/cdn/%s/img/spell/%s", c.baseUrl, version, ultImageName)
	}

	return ults, nil
}
