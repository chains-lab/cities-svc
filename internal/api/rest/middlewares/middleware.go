package middlewares

import (
	"context"

	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/cities-svc/internal/domain/services/admin"
	"github.com/chains-lab/logium"
)

type domain struct {
	admin CityAdminSvc
}

type Middleware struct {
	log    logium.Logger
	domain domain
}

func New(log logium.Logger, gov CityAdminSvc) Middleware {
	return Middleware{
		log: log,
		domain: domain{
			admin: gov,
		},
	}
}

type CityAdminSvc interface {
	Get(ctx context.Context, filters admin.GetFilters) (models.CityAdmin, error)
}
