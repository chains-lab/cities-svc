package middlewares

import (
	"context"

	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/cities-svc/internal/domain/services/citymod"
	"github.com/chains-lab/logium"
)

type domain struct {
	gov govSvc
}

type Middleware struct {
	log    logium.Logger
	domain domain
}

func New(log logium.Logger, gov govSvc) Middleware {
	return Middleware{
		log: log,
		domain: domain{
			gov: gov,
		},
	}
}

type govSvc interface {
	Get(ctx context.Context, filters citymod.GetFilters) (models.CityModer, error)
}
