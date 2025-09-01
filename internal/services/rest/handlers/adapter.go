package handlers

import (
	"net/http"

	"github.com/chains-lab/cities-svc/internal/app"
	"github.com/chains-lab/cities-svc/internal/config"
	"github.com/chains-lab/logium"
)

type Adapter struct {
	app *app.App
	log logium.Logger
	cfg config.Config
}

func NewAdapter(cfg config.Config, log logium.Logger, a *app.App) Adapter {
	return Adapter{
		app: a,
		log: log,
		cfg: cfg,
	}
}

func (s Adapter) Log(r *http.Request) logium.Logger {
	return s.log
}
