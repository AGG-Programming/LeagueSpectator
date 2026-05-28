package handler

import (
	"encoding/json"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/AGG-Programming/LeagueSpectator/internal/league"
	"github.com/AGG-Programming/LeagueSpectator/internal/models"
	"github.com/AGG-Programming/LeagueSpectator/internal/websocket"
)

type Handler struct {
	LeagueClient *league.Client
	WsHub        *websocket.Hub
}

func NewHandler(leagueClient *league.Client, wsHub *websocket.Hub) *Handler {
	return &Handler{
		LeagueClient: leagueClient,
		WsHub:        wsHub,
	}
}

func (h *Handler) Handle() {
	go func() {
		slowTicker := 5 * time.Second
		fastTicker := 500 * time.Millisecond

		ticker := time.NewTicker(slowTicker)
		defer ticker.Stop()

		inGame := false

		for range ticker.C {
			log.Println("Fetching game data...")
			data, err := h.LeagueClient.FetchAllGameData()
			if err != nil {
				if inGame {
					log.Println("Game finished or disconnected. Backing off to slow polling.")
					inGame = false
					ticker.Reset(slowTicker)
				}
				continue
			}
			if !inGame {
				log.Println("Game detected! Switching to fast 500ms streaming loop.")
				inGame = true
				ticker.Reset(fastTicker)
			}
			latest, err := h.LeagueClient.GetLatestPatchVersion()
			if err != nil {
				continue
			}

			var dynamicData models.DynamicGameData
			err = dynamicData.UnmarshalJSON(data)
			if err != nil {
				continue
			}

			completeData, err := h.SetIconURLs(&dynamicData, latest)

			msg, err := json.Marshal(completeData)
			if err != nil {
				log.Println("Error marshalling dynamic data: ", err)
				continue
			}

			h.WsHub.Broadcast <- msg
		}
	}()
}

func (h *Handler) SetIconURLs(dynamicData *models.DynamicGameData, latest string) (*models.DynamicGameData, error) {
	baseURL := "https://ddragon.leagueoflegends.com/cdn/"
	re := regexp.MustCompile(`\d+$`)

	for _, objective := range dynamicData.BlueTeam.Objectives {
		objective.Icon = "./assets/" + objective.Key + ".png"
	}
	for _, objective := range dynamicData.RedTeam.Objectives {
		objective.Icon = "./assets/" + objective.Key + ".png"
	}

	for _, bPlayer := range dynamicData.BlueTeam.Players {
		bPlayer.Icon = baseURL + latest + "/img/champion/" + bPlayer.ChampionName + ".png"

		for _, item := range bPlayer.Items {
			item.Icon = baseURL + latest + "/img/item/" + item.Id + ".png"
		}
		for _, spell := range bPlayer.Spells {
			name := strings.Split(spell.Extended, "_")[2]
			spell.Icon = baseURL + latest + "/img/spell/" + name + ".png"
		}

		primary := re.FindString(bPlayer.Runes.Primary.Extended)
		secondary := re.FindString(bPlayer.Runes.Primary.Extended)

		bPlayer.Runes.Keystone.Icon = baseURL + "img/perk-images/Styles/" + bPlayer.Runes.Primary.DisplayName + "/" + bPlayer.Runes.Keystone.DisplayName + "/" + bPlayer.Runes.Keystone.DisplayName + ".png"
		bPlayer.Runes.Primary.Icon = baseURL + "img/perk-images/Styles/" + primary + "_" + bPlayer.Runes.Primary.DisplayName + "/" + bPlayer.Runes.Primary.DisplayName + ".png"
		bPlayer.Runes.Secondary.Icon = baseURL + "img/perk-images/Styles/" + secondary + "_" + bPlayer.Runes.Primary.DisplayName + "/" + bPlayer.Runes.Primary.DisplayName + ".png"
	}
	for _, rPlayer := range dynamicData.RedTeam.Players {
		rPlayer.Icon = baseURL + latest + "/img/champion/" + rPlayer.ChampionName + ".png"

		for _, item := range rPlayer.Items {
			item.Icon = baseURL + latest + "/img/item/" + item.Id + ".png"
		}
		for _, spell := range rPlayer.Spells {
			name := strings.Split(spell.Extended, "_")[2]
			spell.Icon = baseURL + latest + "/img/spell/" + name + ".png"
		}

		primary := re.FindString(rPlayer.Runes.Primary.Extended)
		secondary := re.FindString(rPlayer.Runes.Primary.Extended)

		rPlayer.Runes.Keystone.Icon = baseURL + "img/perk-images/Styles/" + rPlayer.Runes.Primary.DisplayName + "/" + rPlayer.Runes.Keystone.DisplayName + "/" + rPlayer.Runes.Keystone.DisplayName + ".png"
		rPlayer.Runes.Primary.Icon = baseURL + "img/perk-images/Styles/" + primary + "_" + rPlayer.Runes.Primary.DisplayName + "/" + rPlayer.Runes.Primary.DisplayName + ".png"
		rPlayer.Runes.Secondary.Icon = baseURL + "img/perk-images/Styles/" + secondary + "_" + rPlayer.Runes.Primary.DisplayName + "/" + rPlayer.Runes.Primary.DisplayName + ".png"
	}

	return dynamicData, nil
}
