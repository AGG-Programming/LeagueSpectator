package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/AGG-Programming/LeagueSpectator/internal/ddragon"
	"github.com/AGG-Programming/LeagueSpectator/internal/handler"
	"github.com/AGG-Programming/LeagueSpectator/internal/league"
	"github.com/AGG-Programming/LeagueSpectator/internal/pl"
	"github.com/AGG-Programming/LeagueSpectator/internal/processor"
	"github.com/AGG-Programming/LeagueSpectator/internal/websocket"
	"github.com/joho/godotenv"
)

func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:63342")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Api-Key")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func main() {
	_ = godotenv.Load()

	token := os.Getenv("TOKEN")
	targetTeam := os.Getenv("TARGET_TEAM")
	if token == "" || targetTeam == "" {
		log.Printf("TOKEN or TARGET_TEAM is not set. Will not be able to fetch data from Prime League.")
	}

	plClient := pl.NewClient(token)
	ddragonClient := ddragon.NewClient()
	leagueClient := league.NewClient()
	wsHub := websocket.NewHub()
	cache, err := ddragon.NewCache(ddragonClient)
	if err != nil {
		log.Fatal("cannot create cache: ", err)
	}
	proc := processor.NewProcessor(cache)

	var targetTeamID int
	if targetTeam == "" {
		targetTeamID = 0
	} else {
		targetTeamID, err = strconv.Atoi(targetTeam)
		if err != nil {
			log.Fatal("cannot parse TARGET_TEAM. Must be a number.")
		}
	}

	handlerClient := handler.NewHandler(leagueClient, wsHub, proc, plClient, targetTeamID)

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
	http.HandleFunc("/pl", withCORS(func(w http.ResponseWriter, r *http.Request) {
		handlerClient.HandlePl(w, r)
	}))

	handlerClient.Handle()

	log.Println("Listening on port 8080")
	if err = http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
