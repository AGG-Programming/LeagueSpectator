package handler

import (
	"encoding/json"
	"log"
	"time"

	"github.com/AGG-Programming/LeagueSpectator/internal/league"
	"github.com/AGG-Programming/LeagueSpectator/internal/processor"
	"github.com/AGG-Programming/LeagueSpectator/internal/websocket"
)

type Handler struct {
	LeagueClient *league.Client
	WsHub        *websocket.Hub
	Processor    *processor.Processor
}

func NewHandler(leagueClient *league.Client, wsHub *websocket.Hub, processor *processor.Processor) *Handler {
	return &Handler{
		LeagueClient: leagueClient,
		WsHub:        wsHub,
		Processor:    processor,
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
			data, err := h.LeagueClient.FetchAllGameData()
			if err != nil {
				if inGame {
					log.Println("Game finished or disconnected. Backing off to slow polling.")
					inGame = false
					ticker.Reset(slowTicker)
				}
				log.Println("Error fetching game data: ", err)
				continue
			}
			if !inGame {
				log.Println("Game detected! Switching to fast 500ms streaming loop.")
				inGame = true
				ticker.Reset(fastTicker)
			}

			dynamicData, err := h.Processor.Transformer(data)

			msg, err := json.Marshal(dynamicData)
			if err != nil {
				log.Println("Error marshalling dynamic data: ", err)
				continue
			}

			h.WsHub.Broadcast <- msg
		}
	}()
}
