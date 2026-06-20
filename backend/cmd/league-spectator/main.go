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
	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
)

type Config struct {
	Token      string `toml:"token"`
	TargetTeam int    `toml:"target_team"`
}

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
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal("cannot resolve executable path: ", err)
	}
	exeDir := filepath.Dir(exePath)
	configPath := filepath.Join(exeDir, "config.toml")
	var config Config

	if _, err = toml.DecodeFile(configPath, &config); err != nil {
		log.Printf("Warning: Could not load config.toml (%v). Using empty defaults.", err)
		_ = godotenv.Load()
		config.Token = os.Getenv("TOKEN")
		config.TargetTeam, _ = strconv.Atoi(os.Getenv("TARGET_TEAM"))
	}
	if config.Token == "" || config.TargetTeam == 0 {
		log.Printf("TOKEN or TARGET_TEAM is not set correctly in config.toml. Will not be able to fetch data from Prime League.")
	}

	plClient := pl.NewClient(config.Token)
	ddragonClient := ddragon.NewClient()
	leagueClient := league.NewClient()
	wsHub := websocket.NewHub()
	cache, err := ddragon.NewCache(ddragonClient)
	if err != nil {
		log.Fatal("cannot create cache: ", err)
	}
	proc := processor.NewProcessor(cache)
	handlerClient := handler.NewHandler(leagueClient, wsHub, proc, plClient, config.TargetTeam)

	go wsHub.Run()

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
