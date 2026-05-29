package ddragon

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func (c *Client) GetItems(version string) (map[string]string, error) {
	resp, err := c.httpClient.Get(c.baseUrl + "/cdn/" + version + "/data/en_US/item.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var payload ItemResponse
	if err = json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	items := make(map[string]string)
	for id, item := range payload.Data {
		items[id] = fmt.Sprintf("%s/cdn/%s/img/item/%s.png", c.baseUrl, version, item.Image.Full)
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
		spells[id] = fmt.Sprintf("%s/cdn/%s/img/spell/%s.png", c.baseUrl, version, spell.Image.Full)
	}
	return spells, nil
}
