package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AGG-Programming/ProxyServer/internal/auth"
	"github.com/AGG-Programming/ProxyServer/internal/proxy"
	"github.com/AGG-Programming/ProxyServer/internal/store"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	plBearer := os.Getenv("UPSTREAM_TOKEN")
	connStr := os.Getenv("DATABASE_URL")
	baseURL := os.Getenv("BASE_URL")
	pepper := os.Getenv("API_KEY_PEPPER")
	if plBearer == "" || connStr == "" || baseURL == "" || pepper == "" {
		log.Fatal("PRIME_LEAGUE_BEARER, DATABASE_URL, BASE_URL or API_KEY_PEPPER is not set")
	}

	dbPool, err := store.NewDbPool(ctx, connStr)
	if err != nil {
		log.Fatal("cannot create db pool: ", err)
	}
	defer dbPool.Close()

	mux := http.NewServeMux()

	protectedProxyHandler := auth.ApiKeyMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxy.ProxyRequest(w, r, plBearer, baseURL)
	}), dbPool, pepper)

	mux.Handle("/api/pl/", protectedProxyHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Println("Listening on port 8080")
		if err = server.ListenAndServe(); err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = server.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}
	log.Println("Server Shutdown Complete")
}
