package app

import (
	"context"

	"github.com/chains-lab/cities-svc/internal/app/entities"
	"github.com/chains-lab/cities-svc/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	country entities.Country
	city    entities.City
	gov     entities.Gov

	db *pgxpool.Pool
}

func NewApp(ctx context.Context, cfg config.Config, db *pgxpool.Pool) App {
	pool, err := pgxpool.New(ctx, cfg.Database.SQL.URL)
	if err != nil {
		return nil, err
	}

	return App{}
}
