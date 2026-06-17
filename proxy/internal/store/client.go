package store

import (
	"context"

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
	defer pool.Close()
	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}

func GetUserByKey(ctx context.Context, pool *pgxpool.Pool, key string) (User, error) {
	query := `SELECT user FROM api_keys WHERE api_key = $1`

	var u User
	err := pool.QueryRow(ctx, query, key).Scan(&u.ID, &u.ApiKey, &u.Active)
	if err != nil {
		return User{}, err
	}
	return u, nil
}
