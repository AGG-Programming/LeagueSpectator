package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/AGG-Programming/LeagueSpectator/internal/league"
	"github.com/AGG-Programming/LeagueSpectator/internal/pl"
	"github.com/AGG-Programming/LeagueSpectator/internal/websocket"
	"github.com/AGG-Programming/LeagueSpectator/pkg/models"
)

type PrimeLeague interface {
	GetLeagueData(ctx context.Context) (*pl.PrimeLeagueResponse, error)
}

type Processor interface {
	Transformer(data league.GameResponse) (models.DynamicGameData, error)
	TransformPL(data pl.PrimeLeagueResponse, targetID int) (*models.PrimeLeague, error)
}

type Handler struct {
	LeagueClient *league.Client
	WsHub        *websocket.Hub
	Processor    Processor
	PlClient     PrimeLeague
	TargetTeam   int
}

func NewHandler(leagueClient *league.Client, wsHub *websocket.Hub, processor Processor, plClient PrimeLeague, targetTeam int) *Handler {
	return &Handler{
		LeagueClient: leagueClient,
		WsHub:        wsHub,
		Processor:    processor,
		PlClient:     plClient,
		TargetTeam:   targetTeam,
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

func (h *Handler) HandlePl(w http.ResponseWriter, r *http.Request) {
	data, err := h.PlClient.GetLeagueData(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := h.Processor.TransformPL(*data, h.TargetTeam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
