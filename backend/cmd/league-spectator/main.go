package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/AGG-Programming/LeagueSpectator/internal/league"
	"github.com/AGG-Programming/LeagueSpectator/internal/models"
	"github.com/AGG-Programming/LeagueSpectator/internal/websocket"
)

func main() {
	leagueClient := league.NewClient()
	wsHub := websocket.NewHub()

	go wsHub.Run()

	exePath, err := os.Executable()
	if err != nil {
		log.Fatal("cannot resolve executable path: ", err)
	}
	exeDir := filepath.Dir(exePath)
	frontendPath := filepath.Join(exeDir, "frontend")

	frontendDir := http.Dir(frontendPath)
	http.Handle("/", http.FileServer(frontendDir))

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(wsHub, w, r)
	})

	go func() {
		slowTicker := 5 * time.Second
		fastTicker := 500 * time.Millisecond

		ticker := time.NewTicker(slowTicker)
		defer ticker.Stop()

		inGame := false

		for range ticker.C {
			log.Println("Fetching game data...")
			data, err := leagueClient.FetchAllGameData()
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
			var dynamicData models.DynamicGameData
			err = dynamicData.UnmarshalJSON(data)
			if err != nil {
				continue
			}

			msg, err := json.Marshal(dynamicData)
			if err != nil {
				log.Println("Error marshalling dynamic data: ", err)
				continue
			}

			wsHub.Broadcast <- msg
		}
	}()

	log.Println("Listening on port 8080")
	if err = http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
