package main

import (
	"log"
	"net/http"
	"time"

	"github.com/AGG-Programming/LeagueSpectator/internal/league"
	"github.com/AGG-Programming/LeagueSpectator/internal/websocket"
)

func main() {
	leagueClient := league.NewClient()
	wsHub := websocket.NewHub()

	go wsHub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(wsHub, w, r)
	})

	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			data, err := leagueClient.FetchAllGameData()
			if err != nil {
				continue
			}

			wsHub.Broadcast <- data
		}
	}()

	log.Println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
