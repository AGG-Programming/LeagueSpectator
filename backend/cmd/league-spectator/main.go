package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

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

type analyzerStart struct {
	cmd *exec.Cmd
	via string
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

func existingFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func uniqueDirs(paths ...string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		if strings.TrimSpace(p) == "" {
			continue
		}
		clean := filepath.Clean(p)
		if _, ok := seen[clean]; ok {
			continue
		}
		seen[clean] = struct{}{}
		out = append(out, clean)
	}
	return out
}

func startDisplayAnalyzer(exeDir string) (*analyzerStart, error) {
	cwd, _ := os.Getwd()
	searchDirs := uniqueDirs(
		exeDir,
		cwd,
		filepath.Dir(cwd),
		filepath.Join(exeDir, ".."),
	)

	binaryNames := []string{"displayAnalyzer.exe", "displayAnalyzer"}
	for _, dir := range searchDirs {
		for _, name := range binaryNames {
			direct := filepath.Join(dir, name)
			nested := filepath.Join(dir, "displayAnalyzer", "dist", name)
			for _, candidate := range []string{direct, nested} {
				if !existingFile(candidate) {
					continue
				}
				cmd := exec.Command(candidate)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Start(); err != nil {
					log.Printf("Analyzer candidate %q failed to start: %v", candidate, err)
					continue
				}
				go func() {
					if err := cmd.Wait(); err != nil {
						log.Printf("Display analyzer exited with error: %v", err)
						return
					}
					log.Printf("Display analyzer stopped.")
				}()
				return &analyzerStart{
					cmd: cmd,
					via: candidate,
				}, nil
			}
		}
	}

	pyScriptCandidates := make([]string, 0, len(searchDirs))
	for _, dir := range searchDirs {
		pyScriptCandidates = append(pyScriptCandidates, filepath.Join(dir, "displayAnalyzer", "main.py"))
	}

	pythonCommands := [][]string{
		{"python3"},
		{"python"},
	}
	if runtime.GOOS == "windows" {
		pythonCommands = append([][]string{{"py", "-3"}}, pythonCommands...)
	}

	for _, script := range pyScriptCandidates {
		if !existingFile(script) {
			continue
		}
		for _, baseCmd := range pythonCommands {
			if _, err := exec.LookPath(baseCmd[0]); err != nil {
				continue
			}
			args := append(baseCmd[1:], script)
			cmd := exec.Command(baseCmd[0], args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Dir = filepath.Dir(script)
			if err := cmd.Start(); err != nil {
				return nil, fmt.Errorf("failed starting python analyzer via %q %q: %w", baseCmd[0], script, err)
			}
			go func() {
				if err := cmd.Wait(); err != nil {
					log.Printf("Display analyzer exited with error: %v", err)
					return
				}
				log.Printf("Display analyzer stopped.")
			}()
			return &analyzerStart{
				cmd: cmd,
				via: fmt.Sprintf("%s %s", baseCmd[0], script),
			}, nil
		}
	}

	return nil, errors.New("could not find analyzer binary or python script")
}

func stopDisplayAnalyzer(started *analyzerStart) {
	if started == nil || started.cmd == nil || started.cmd.Process == nil {
		return
	}
	if runtime.GOOS == "windows" {
		if err := started.cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
			log.Printf("Could not stop display analyzer: %v", err)
		}
		return
	}

	if err := started.cmd.Process.Signal(os.Interrupt); err != nil && !errors.Is(err, os.ErrProcessDone) {
		if killErr := started.cmd.Process.Kill(); killErr != nil && !errors.Is(killErr, os.ErrProcessDone) {
			log.Printf("Could not stop display analyzer: %v", killErr)
		}
		return
	}

	time.Sleep(750 * time.Millisecond)
	if err := started.cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
		log.Printf("Could not force-stop display analyzer: %v", err)
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

	analyzer, err := startDisplayAnalyzer(exeDir)
	if err != nil {
		log.Printf("Warning: could not auto-start display analyzer: %v", err)
	} else {
		log.Printf("Started display analyzer via %s", analyzer.via)
		defer stopDisplayAnalyzer(analyzer)
	}

	plClient := pl.NewClient(config.Token)
	ddragonClient := ddragon.NewClient()
	leagueClient := league.NewClient()
	wsHub := websocket.NewHub()
	cache, err := ddragon.NewCache(ddragonClient)
	if err != nil {
		log.Printf("cannot create cache: %v", err)
		return
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
	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.ListenAndServe()
	}()

	shutdownSignals := make(chan os.Signal, 1)
	signal.Notify(shutdownSignals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(shutdownSignals)

	select {
	case sig := <-shutdownSignals:
		log.Printf("Received signal %v, shutting down...", sig)
	case srvErr := <-serverErr:
		if !errors.Is(srvErr, http.ErrServerClosed) {
			log.Printf("ListenAndServe: %v", srvErr)
			return
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = server.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("Server shutdown error: %v", err)
	}
}
