package database

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool cria um pool de conexões com o PostgreSQL.
func NewPool(ctx context.Context) (*pgxpool.Pool, error) {
	databaseURL := os.Getenv("DATABASE_URL")

	// Se DATABASE_URL não estiver configurada,
	// monta a conexão usando as variáveis do PostgreSQL.
	if databaseURL == "" {
		user := os.Getenv("POSTGRES_USER")
		password := os.Getenv("POSTGRES_PASSWORD")
		dbName := os.Getenv("POSTGRES_DB")

		if user == "" || password == "" || dbName == "" {
			return nil, errors.New(
				"DATABASE_URL ou as variáveis POSTGRES_USER, POSTGRES_PASSWORD e POSTGRES_DB precisam estar configuradas",
			)
		}

		u := &url.URL{
			Scheme: "postgres",
			User:   url.UserPassword(user, password),
			Host:   "localhost:5432",
			Path:   dbName,
		}

		databaseURL = u.String()
	}

	// Interpreta e valida a configuração da conexão.
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf(
			"configuração do PostgreSQL inválida: %w",
			err,
		)
	}

	// Limita o número de conexões simultâneas.
	config.MaxConns = 5
	config.MinConns = 0

	// Cria o pool de conexões.
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf(
			"não foi possível criar o pool: %w",
			err,
		)
	}

	// Testa a conexão com um limite de cinco segundos.
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf(
			"não foi possível conectar ao PostgreSQL: %w",
			err,
		)
	}

	return pool, nil
}
