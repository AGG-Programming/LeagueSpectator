package store

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID     int
	ApiKey string
	Active bool
}

func NewDbPool(ctx context.Context, connStr string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}
	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	log.Println("Connected to database")

	return pool, nil
}

func hashAPIKey(apiKey, pepper string) (string, error) {
	if pepper == "" {
		return "", errors.New("API_KEY_PEPPER is empty")
	}

	mac := hmac.New(sha256.New, []byte(pepper))
	_, _ = mac.Write([]byte(apiKey))
	return hex.EncodeToString(mac.Sum(nil)), nil
}

func GetUserByKey(ctx context.Context, pool *pgxpool.Pool, key string, pepper string) (User, error) {
	keyHash, err := hashAPIKey(key, pepper)
	if err != nil {
		return User{}, err
	}

	query := `SELECT id, key, active FROM api_keys WHERE key = $1`

	var u User
	err = pool.QueryRow(ctx, query, keyHash).Scan(&u.ID, &u.ApiKey, &u.Active)
	if err != nil {
		return User{}, err
	}
	return u, nil
}
