package app

import (
	"database/sql"

	"github.com/chains-lab/cities-svc/internal/app/entities"
	"github.com/chains-lab/cities-svc/internal/config"
)

type App struct {
	country entities.Country
	city    entities.City
	gov     entities.Gov

	db *sql.DB
}

func NewApp(cfg config.Config) (App, error) {
	pg, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		return App{}, err
	}

	return App{
		country: entities.NewCountry(pg),
		city:    entities.NewCitySvc(pg),
		gov:     entities.NewGov(pg),

		db: pg,
	}, nil
}
