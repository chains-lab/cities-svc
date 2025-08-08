package service

import (
	svccities "github.com/chains-lab/cities-dir-proto/gen/go/cities"
	svccitiesadmins "github.com/chains-lab/cities-dir-proto/gen/go/citiesadmins"
	svccountries "github.com/chains-lab/cities-dir-proto/gen/go/countries"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/service/cities"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/service/citiesadmins"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/service/countries"
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
	svccities.CityServiceServer
	svccitiesadmins.CityAdminServiceServer
	svccountries.CountryServiceServer
}

func NewService(cfg config.Config, app *app.App) Service {
	citiesService := cities.NewService(cfg, app)
	citiesAdminsService := citiesadmins.NewService(cfg, app)
	countriesService := countries.NewService(cfg, app)

	return Service{
		CityServiceServer:      citiesService,
		CityAdminServiceServer: citiesAdminsService,
		CountryServiceServer:   countriesService,
	}
}
