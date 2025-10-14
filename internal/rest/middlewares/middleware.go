package middlewares

import (
	"context"
	"net/http"

	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/cities-svc/internal/domain/services/admin"
	"github.com/chains-lab/logium"
	"github.com/chains-lab/restkit/mdlv"
)

type domain struct {
	admin CityAdminSvc
}

type Service struct {
	log    logium.Logger
	domain domain
}

func New(log logium.Logger, gov CityAdminSvc) Service {
	return Service{
		log: log,
		domain: domain{
			admin: gov,
		},
	}
}

type CityAdminSvc interface {
	Get(ctx context.Context, filters admin.GetFilters) (models.CityAdmin, error)
}

func (s Service) ServiceGrant(serviceName, skService string) func(http.Handler) http.Handler {
	return mdlv.ServiceGrant(serviceName, skService)
}

func (s Service) Auth(userCtxKey interface{}, skUser string) func(http.Handler) http.Handler {
	return mdlv.Auth(userCtxKey, skUser)
}

func (s Service) RoleGrant(userCtxKey interface{}, allowedRoles map[string]bool) func(http.Handler) http.Handler {
	return mdlv.RoleGrant(userCtxKey, allowedRoles)
}
