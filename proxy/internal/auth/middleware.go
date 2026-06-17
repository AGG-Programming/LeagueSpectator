package auth

import (
	"context"
	"log"
	"net/http"

	"github.com/AGG-Programming/ProxyServer/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ApiKeyMiddleware(next http.Handler, dbPool *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		apiKey := r.Header.Get("X-Api-Key")
		user, err := store.GetUserByKey(r.Context(), dbPool, apiKey)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			log.Println(err.Error())
			return
		}

		ctx := context.WithValue(
			r.Context(),
			"user",
			user,
		)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
