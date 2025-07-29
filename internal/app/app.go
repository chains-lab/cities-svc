package app

import "github.com/chains-lab/cities-dir-svc/internal/config"

type App struct {
}

func NewApp(cfg config.Config) (App, error) {
	return App{}, nil
}
