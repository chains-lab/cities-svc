package rest

import (
	"context"
	"net/http"

	"github.com/chains-lab/cities-svc/internal/api/rest/meta"
	"github.com/chains-lab/cities-svc/internal/config"
	"github.com/chains-lab/cities-svc/internal/constant"
	"github.com/chains-lab/gatekit/mdlv"
	"github.com/chains-lab/gatekit/roles"
	"github.com/go-chi/chi/v5"
)

type Handlers interface {
	CreateCountry(w http.ResponseWriter, r *http.Request)
	SearchCountries(w http.ResponseWriter, r *http.Request)
	GetCountry(w http.ResponseWriter, r *http.Request)
	UpdateCountry(w http.ResponseWriter, r *http.Request)
	UpdateCountryStatus(w http.ResponseWriter, r *http.Request)

	SearchCities(w http.ResponseWriter, r *http.Request)
	CreateCity(w http.ResponseWriter, r *http.Request)
	GetCity(w http.ResponseWriter, r *http.Request)
	UpdateCity(w http.ResponseWriter, r *http.Request)
	UpdateCityStatus(w http.ResponseWriter, r *http.Request)

	GetCityBySlug(w http.ResponseWriter, r *http.Request)
	SearchGovs(w http.ResponseWriter, r *http.Request)
	CreateInvite(w http.ResponseWriter, r *http.Request)
	AnswerToInvite(w http.ResponseWriter, r *http.Request)
	GetGov(w http.ResponseWriter, r *http.Request)
	CreateMayor(w http.ResponseWriter, r *http.Request)
	TransferMayor(w http.ResponseWriter, r *http.Request)

	GetOwnGov(w http.ResponseWriter, r *http.Request)
	UpdateOwnGov(w http.ResponseWriter, r *http.Request)
	RefuseOwnGov(w http.ResponseWriter, r *http.Request)
}

func (s *Service) Api(ctx context.Context, cfg config.Config, h Handlers) {
	svc := mdlv.ServiceGrant(constant.ServiceName, cfg.JWT.Service.SecretKey)
	auth := mdlv.Auth(meta.UserCtxKey, cfg.JWT.User.AccessToken.SecretKey)
	sysadmin := mdlv.RoleGrant(meta.UserCtxKey, map[string]bool{
		roles.Admin:     true,
		roles.SuperUser: true,
	})

	s.router.Route("/cities-svc/", func(r chi.Router) {
		r.Use(svc)

		r.Route("/v1", func(r chi.Router) {
			r.Route("/countries", func(r chi.Router) {
				r.With(auth, sysadmin).Post("/", h.CreateCountry)

				r.Get("/", h.SearchCountries)

				r.Route("/{country_id}", func(r chi.Router) {
					r.Get("/", h.GetCountry)

					r.Group(func(r chi.Router) {
						r.Use(auth, sysadmin)
						r.Patch("/", h.UpdateCountry) // частичное обновление
						r.Patch("/status", h.UpdateCountryStatus)
					})
				})
			})

			r.Route("/cities", func(r chi.Router) {
				r.Get("/", h.SearchCities)
				r.Get("/slug/{slug}", h.GetCityBySlug)

				r.With(auth, sysadmin).Post("/", h.CreateCity)

				r.Route("/{city_id}", func(r chi.Router) {
					r.Get("/", h.GetCity)

					r.Group(func(r chi.Router) {
						r.Use(auth)
						r.Patch("/", h.UpdateCity)
						r.Patch("/status", h.UpdateCityStatus)
					})

					r.Route("/govs", func(r chi.Router) {
						r.Get("/", h.SearchGovs)

						r.With(auth).Route("/invite", func(r chi.Router) {
							r.Post("/", h.CreateInvite)
							r.Post("/{token}", h.AnswerToInvite)
						})

						r.Route("/mayor", func(r chi.Router) {
							r.Get("/", h.GetCityMayor) // get current mayor
							r.With(auth, sysadmin).Post("/", h.CreateMayor)
							r.With(auth).Post("/transfer", h.TransferMayor)
						})

						r.With(auth).Route("/me", func(r chi.Router) {
							r.Get("/", h.GetOwnGov)
							r.Put("/", h.UpdateOwnGov)
							r.Delete("/", h.RefuseOwnGov)
						})

						r.Route("/{user_id}", func(r chi.Router) {
							r.Get("/", h.GetGov)
							r.With(auth).Delete("/", h.DeleteGov)
						})

					})
				})
			})
		})
	})

	s.Start(ctx)

	<-ctx.Done()
	s.Stop(ctx)
}
