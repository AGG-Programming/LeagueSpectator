package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/AGG-Programming/LeagueSpectator/internal/handler"
	"github.com/AGG-Programming/LeagueSpectator/internal/league"
	"github.com/AGG-Programming/LeagueSpectator/internal/websocket"
)

func main() {
	leagueClient := league.NewClient()
	wsHub := websocket.NewHub()
	handlerClient := handler.NewHandler(leagueClient, wsHub)

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

	handlerClient.Handle()

	log.Println("Listening on port 8080")
	if err = http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
