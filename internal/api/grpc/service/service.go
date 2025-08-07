package service

import (
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/service/cities"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/service/citiesadmins"
	countries2 "github.com/chains-lab/cities-dir-svc/internal/api/grpc/service/countries"
	"github.com/chains-lab/cities-dir-svc/internal/app"
	"github.com/chains-lab/cities-dir-svc/internal/config"
)

type Cities interface {
}

type CitiesAdmins interface {
}

type Countries interface {
}

type Service struct {
	Cities
	CitiesAdmins
	Countries
}

func NewService(cfg config.Config, app *app.App) Service {
	citiesService := cities.NewService(cfg, app)
	citiesAdminsService := citiesadmins.NewService(cfg, app)
	countriesService := countries2.NewService(cfg, app)

	return Service{
		Cities:       citiesService,
		CitiesAdmins: citiesAdminsService,
		Countries:    countriesService,
	}
}
