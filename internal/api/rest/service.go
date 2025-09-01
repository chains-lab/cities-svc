package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/chains-lab/cities-svc/internal/api/rest/handlers"
	"github.com/chains-lab/cities-svc/internal/api/rest/meta"
	"github.com/chains-lab/cities-svc/internal/app"
	"github.com/chains-lab/cities-svc/internal/config"
	"github.com/chains-lab/gatekit/mdlv"
	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/logium"
	"github.com/go-chi/chi/v5"
)

type Rest struct {
	server   *http.Server
	router   *chi.Mux
	handlers handlers.Service

	log logium.Logger
	cfg config.Config
}

func NewRest(cfg config.Config, log logium.Logger, app *app.App) Rest {
	logger := log.WithField("module", "api")
	router := chi.NewRouter()
	server := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: router,
	}
	hands := handlers.NewService(cfg, logger, app)

	router.Use()

	return Rest{
		handlers: hands,
		router:   router,
		server:   server,
		log:      logger,
		cfg:      cfg,
	}
}

func (a *Rest) Run(ctx context.Context) {
	//svcAuth := mdlv.ServiceAuthMdl(constant.ServiceName, a.cfg.JWT.Service.SecretKey)
	userAuth := mdlv.AuthMdl(meta.UserCtxKey, a.cfg.JWT.User.AccessToken.SecretKey)
	adminGrant := mdlv.AccessGrant(meta.UserCtxKey, roles.Admin, roles.SuperUser)

	a.router.Route("/cities-svc/", func(r chi.Router) {
		//r.Use(svcAuth)
		r.Route("/v1", func(r chi.Router) {
			r.Route("/countries", func(r chi.Router) {
				r.Get("/", a.handlers.SearchCountries)
				r.Post("/", a.handlers.CreateCountry)

				r.Route("/{country_id}", func(r chi.Router) {
					r.Get("/", a.handlers.GetCountry)
					r.Put("/", a.handlers.UpdateCountry)
					r.Route("/status", func(r chi.Router) {
						r.Use(userAuth, adminGrant)
						r.Put("/{status}", a.handlers.ChangeCountryStatus)
					})
				})
			})

			r.Route("/cities", func(r chi.Router) {
				r.Get("/", a.handlers.SearchCities)
				r.Get("/nearby", a.handlers.GetNearbyCities)

				r.Post("/", a.handlers.CreateCity)

				r.Route("/{city_id}", func(r chi.Router) {
					r.Get("/", a.handlers.GetCity)
					r.Put("/", a.handlers.UpdateCity)

					r.Route("/status", func(r chi.Router) {
						r.Use(userAuth, adminGrant)
						r.Put("/{status}", a.handlers.ChangeCityStatus)
					})

					r.Route("/govs", func(r chi.Router) {
						r.Get("/", a.handlers.SearchGovs)
						r.Post("/", a.handlers.CreateGov)

						r.Route("/{gov_id}", func(r chi.Router) {
							r.Get("/", a.handlers.GetGov)
							r.Put("/", a.handlers.UpdateGov)
						})
					})
				})
			})

			r.Get("/{slug}", a.handlers.GetCityBySlug)
		})
	})

	a.Start(ctx)

	<-ctx.Done()
	a.Stop(ctx)
}

func (a *Rest) Start(ctx context.Context) {
	go func() {
		a.log.Infof("Starting server on port %s", a.cfg.Server.Port)
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Fatalf("Server failed to start: %v", err)
		}
	}()
}

func (a *Rest) Stop(ctx context.Context) {
	a.log.Info("Shutting down server...")
	if err := a.server.Shutdown(ctx); err != nil {
		a.log.Errorf("Server shutdown failed: %v", err)
	}
}
