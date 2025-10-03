package rest

import (
	"context"
	"net/http"

	"github.com/chains-lab/cities-svc/internal"
	"github.com/chains-lab/cities-svc/internal/api/rest/meta"
	"github.com/chains-lab/cities-svc/internal/domain/enum"
	"github.com/chains-lab/gatekit/mdlv"
	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/logium"
	"github.com/go-chi/chi/v5"
)

type Handlers interface {
	CreateCountry(w http.ResponseWriter, r *http.Request)
	ListCountries(w http.ResponseWriter, r *http.Request)
	GetCountry(w http.ResponseWriter, r *http.Request)
	UpdateCountry(w http.ResponseWriter, r *http.Request)
	UpdateCountryStatus(w http.ResponseWriter, r *http.Request)

	ListCities(w http.ResponseWriter, r *http.Request)
	CreateCity(w http.ResponseWriter, r *http.Request)
	GetCity(w http.ResponseWriter, r *http.Request)
	UpdateCity(w http.ResponseWriter, r *http.Request)
	UpdateCityStatus(w http.ResponseWriter, r *http.Request)

	GetCityBySlug(w http.ResponseWriter, r *http.Request)
	ListGovs(w http.ResponseWriter, r *http.Request)
	CreateInvite(w http.ResponseWriter, r *http.Request)
	AcceptInvite(w http.ResponseWriter, r *http.Request)
	GetCityAdmin(w http.ResponseWriter, r *http.Request)
	DeleteCityAdmin(w http.ResponseWriter, r *http.Request)

	GetOwnCityAdmin(w http.ResponseWriter, r *http.Request)
	UpdateOwnCityAdmin(w http.ResponseWriter, r *http.Request)
	RefuseOwnCityAdmin(w http.ResponseWriter, r *http.Request)
}

type Middlewares interface {
	CityAdminRoles(
		UserCtxKey interface{},
		allowedGovRoles map[string]bool,
		allowedSysadminRoles map[string]bool,
	) func(http.Handler) http.Handler
}

func Run(ctx context.Context, cfg internal.Config, log logium.Logger, h Handlers, m Middlewares) {
	svc := mdlv.ServiceGrant(cfg.Service.Name, cfg.JWT.Service.SecretKey)
	auth := mdlv.Auth(meta.UserCtxKey, cfg.JWT.User.AccessToken.SecretKey)
	sysadmin := mdlv.RoleGrant(meta.UserCtxKey, map[string]bool{
		roles.Admin: true,
	})

	user := mdlv.RoleGrant(meta.UserCtxKey, map[string]bool{
		roles.User: true,
	})

	cityMod := m.CityAdminRoles(meta.UserCtxKey, map[string]bool{
		enum.CityAdminRoleHead:      true,
		enum.CityAdminRoleModerator: true,
	}, map[string]bool{
		roles.Admin: true,
	})

	cityStuff := m.CityAdminRoles(meta.UserCtxKey, map[string]bool{
		enum.CityAdminRoleHead:      true,
		enum.CityAdminRoleModerator: true,
		enum.CityAdminMember:        true,
	}, map[string]bool{
		roles.Admin: true,
	})

	r := chi.NewRouter()

	r.Route("/cities-svc/", func(r chi.Router) {
		r.Use(svc)

		r.Route("/v1", func(r chi.Router) {
			r.Route("/countries", func(r chi.Router) {
				r.With(auth, sysadmin).Post("/", h.CreateCountry)

				r.Get("/", h.ListCountries)

				r.Route("/{country_id}", func(r chi.Router) {
					r.Get("/", h.GetCountry)

					r.Group(func(r chi.Router) {
						r.Use(auth, sysadmin)
						r.Put("/", h.UpdateCountry)
						r.Put("/status", h.UpdateCountryStatus)
					})
				})
			})

			r.Route("/cities", func(r chi.Router) {
				r.Get("/", h.ListCities)
				r.Get("/slug/{slug}", h.GetCityBySlug)

				r.With(auth, sysadmin).Post("/", h.CreateCity)

				r.Route("/{city_id}", func(r chi.Router) {
					r.Get("/", h.GetCity)

					r.With(auth, cityMod).Put("/", h.UpdateCity)
					r.With(auth, sysadmin).Put("/status", h.UpdateCityStatus)

					r.Route("/govs", func(r chi.Router) {
						r.Get("/", h.ListGovs)

						r.With(auth).Route("/invite", func(r chi.Router) {
							r.With(user, cityMod).Post("/", h.CreateInvite)
							r.With(user).Post("/{token}", h.AcceptInvite)
						})

						r.With(auth, user, cityStuff).Route("/me", func(r chi.Router) {
							r.Get("/", h.GetOwnCityAdmin)
							r.Put("/", h.UpdateOwnCityAdmin)
							r.Delete("/", h.RefuseOwnCityAdmin)
						})

						r.Route("/{user_id}", func(r chi.Router) {
							r.Get("/", h.GetCityAdmin)
							r.With(auth, user).Delete("/", h.DeleteCityAdmin)
						})
					})
				})
			})
		})
	})

	log.Infof("starting REST service on %s", cfg.Rest.Port)

	<-ctx.Done()

	log.Info("shutting down REST service")
}
