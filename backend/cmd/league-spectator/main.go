package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/AGG-Programming/LeagueSpectator/internal/ddragon"
	"github.com/AGG-Programming/LeagueSpectator/internal/handler"
	"github.com/AGG-Programming/LeagueSpectator/internal/league"
	"github.com/AGG-Programming/LeagueSpectator/internal/processor"
	"github.com/AGG-Programming/LeagueSpectator/internal/websocket"
)

func startAnalyzer(exeDir string) *exec.Cmd {
	analyzerPath := filepath.Join(exeDir, "displayAnalyzer/dist/displayAnalyzer")
	cmd := exec.Command(analyzerPath)

	if err := cmd.Start(); err != nil {
		log.Printf("error starting analyzer: %v", err)
		return nil
	}
	log.Fatal("displayAnalyzer started")
	return cmd
}

func main() {
	ddragonClient := ddragon.NewClient()
	leagueClient := league.NewClient()
	wsHub := websocket.NewHub()
	cache, err := ddragon.NewCache(ddragonClient)
	if err != nil {
		log.Fatal("cannot create cache: ", err)
	}
	proc := processor.NewProcessor(cache)
	handlerClient := handler.NewHandler(leagueClient, wsHub, proc)

	go wsHub.Run()

	exePath, err := os.Executable()
	if err != nil {
		log.Fatal("cannot resolve executable path: ", err)
	}
	exeDir := filepath.Dir(exePath)
	startAnalyzer(exeDir)
	frontendPath := filepath.Join(exeDir, "frontend")

	frontendDir := http.Dir(frontendPath)
	http.Handle("/", http.FileServer(frontendDir))

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(wsHub, w, r)
	})

	handlerClient.Handle()

	log.Println("Listening on port 8080")
	if err = http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
