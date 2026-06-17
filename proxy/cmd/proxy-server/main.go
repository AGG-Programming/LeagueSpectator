package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/AGG-Programming/ProxyServer/internal/auth"
	"github.com/AGG-Programming/ProxyServer/internal/proxy"
	"github.com/AGG-Programming/ProxyServer/internal/store"
)

func main() {
	ctx := context.Background()

	plBearer := os.Getenv("UPSTREAM_TOKEN")
	connStr := os.Getenv("DATABASE_URL")
	baseURL := os.Getenv("BASE_URL")
	if plBearer == "" || connStr == "" {
		log.Fatal("PRIME_LEAGUE_BEARER or  is not set")
	}

	dbPool, err := store.NewDbPool(ctx, connStr)

	mux := http.NewServeMux()

	protectedProxyHandler := auth.ApiKeyMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxy.ProxyRequest(w, r, plBearer, baseURL)
	}), dbPool, ctx)

	mux.Handle("/api/pl/", protectedProxyHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Listening on port 8080")
	if err = server.ListenAndServe(); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
