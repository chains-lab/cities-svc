package countryadmin

import (
	svc "github.com/chains-lab/cities-dir-proto/gen/go/svc/countryadmin"
	"github.com/chains-lab/cities-dir-svc/internal/app"
	"github.com/chains-lab/cities-dir-svc/internal/config"
)

type Service struct {
	app *app.App
	cfg config.Config

	svc.UnimplementedCountryAdminServiceServer
}

func NewService(cfg config.Config, app *app.App) Service {
	return Service{
		app: app,
		cfg: cfg,
	}
}
