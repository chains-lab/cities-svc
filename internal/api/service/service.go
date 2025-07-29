package service

import (
	"context"

	"github.com/chains-lab/cities-dir-svc/internal/api/interceptors"
	"github.com/chains-lab/cities-dir-svc/internal/config"
)

type App interface {
}

type Service struct {
	app App
	cfg config.Config
}

func NewService(cfg config.Config, app *app.App) Service {
	return Service{
		app: app,
		cfg: cfg,
	}
}

func Meta(ctx context.Context) interceptors.MetaData {
	md, ok := ctx.Value(interceptors.MetaCtxKey).(interceptors.MetaData)
	if !ok {
		return interceptors.MetaData{}
	}
	return md
}
