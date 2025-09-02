package rest

import (
	"context"
	"net/http"

	"github.com/chains-lab/cities-svc/internal/services/rest/handlers"
	"github.com/go-chi/chi/v5"
)

type Handlers interface {
	CreateCountry(w http.ResponseWriter, r *http.Request)
	SearchCountries(w http.ResponseWriter, r *http.Request)
	GetCountry(w http.ResponseWriter, r *http.Request)
	UpdateCountry(w http.ResponseWriter, r *http.Request)
	ChangeCountryStatus(w http.ResponseWriter, r *http.Request)

	SearchCities(w http.ResponseWriter, r *http.Request)
	GetNearbyCity(w http.ResponseWriter, r *http.Request)
	CreateCity(w http.ResponseWriter, r *http.Request)
	GetCity(w http.ResponseWriter, r *http.Request)
	UpdateCity(w http.ResponseWriter, r *http.Request)
	ChangeCityStatus(w http.ResponseWriter, r *http.Request)
	GetCityBySlug(w http.ResponseWriter, r *http.Request)

	SearchGovs(w http.ResponseWriter, r *http.Request)
	CreateGov(w http.ResponseWriter, r *http.Request)
	GetGov(w http.ResponseWriter, r *http.Request)
	UpdateGov(w http.ResponseWriter, r *http.Request)
	CreateMayor(w http.ResponseWriter, r *http.Request)
	TransferMayorRole(w http.ResponseWriter, r *http.Request)
}

func (s *Service) Api(ctx context.Context, h handlers.Adapter) {
	politics := s.buildPolicies()

	s.router.Route("/cities-svc/", func(r chi.Router) {
		r.Use(politics.Service)

		r.Route("/v1", func(r chi.Router) {
			r.Route("/countries", func(r chi.Router) {
				r.With(politics.Sysadmin).Post("/", h.CreateCountry)

				r.Get("/", h.SearchCountries)

				r.Route("/{country_id}", func(r chi.Router) {
					r.Get("/", h.GetCountry)

					r.With(politics.Sysadmin).Put("/", h.UpdateCountry)
					r.With(politics.Sysadmin).Put("/status", h.UpdateCountryStatus)
				})
			})

			r.Route("/cities", func(r chi.Router) {
				r.Get("/", h.SearchCities)

				r.With(politics.Sysadmin).Post("/", h.CreateCity)

				r.Route("/{city_id}", func(r chi.Router) {
					r.Get("/", h.GetCity)
					r.Put("/", h.UpdateCity) //TODO

					r.Put("/status", h.UpdateCityStatus) //TODO
				})
			})

			r.Get("slug/{slug}", h.GetCityBySlug)

			r.Route("/govs", func(r chi.Router) {
				r.Get("/", h.SearchGovs)
				r.Post("/", h.CreateGov)

				r.Route("/{gov_id}", func(r chi.Router) {
					r.Get("/", h.GetGov)
					r.Put("/", h.UpdateGov)
				})

				r.Route("/mayor", func(r chi.Router) {
					r.With(politics.Sysadmin).
						Post("/", h.CreateMayor)

					r.Post("/transfer/{user_id}", h.TransferMayor)
				})

				r.Route("/me", func(r chi.Router) {
					r.Get("/", h.GetOwnCurrentGov)
					r.Put("/", h.UpdateOwnCurrenteGov)
					r.Delete("/", h.RefuseOwnCurrentGov)
				})
			})
		})
	})

	s.Start(ctx)

	<-ctx.Done()
	s.Stop(ctx)
}
