package rest

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/chains-lab/cities-svc/internal"
	"github.com/chains-lab/cities-svc/internal/rest/meta"
	"github.com/chains-lab/logium"
	"github.com/chains-lab/restkit/roles"
	"github.com/go-chi/chi/v5"
)

type Handlers interface {
	ListCities(w http.ResponseWriter, r *http.Request)
	CreateCity(w http.ResponseWriter, r *http.Request)
	GetCity(w http.ResponseWriter, r *http.Request)
	UpdateCity(w http.ResponseWriter, r *http.Request)
	UpdateCityStatus(w http.ResponseWriter, r *http.Request)

	GetCityBySlug(w http.ResponseWriter, r *http.Request)
	ListAdmins(w http.ResponseWriter, r *http.Request)
	SentInvite(w http.ResponseWriter, r *http.Request)
	AnswerInvite(w http.ResponseWriter, r *http.Request)
	GetCityAdmin(w http.ResponseWriter, r *http.Request)
	DeleteCityAdmin(w http.ResponseWriter, r *http.Request)

	GetMyCityAdmin(w http.ResponseWriter, r *http.Request)
	UpdateCityAdmin(w http.ResponseWriter, r *http.Request)
	UpdateMyCityAdmin(w http.ResponseWriter, r *http.Request)
	RefuseMyCityAdmin(w http.ResponseWriter, r *http.Request)
}

type Middlewares interface {
	Auth(userCtxKey interface{}, skUser string) func(http.Handler) http.Handler
	RoleGrant(userCtxKey interface{}, allowedRoles map[string]bool) func(http.Handler) http.Handler
}

func Run(ctx context.Context, cfg internal.Config, log logium.Logger, m Middlewares, h Handlers) {
	auth := m.Auth(meta.UserCtxKey, cfg.JWT.User.AccessToken.SecretKey)

	sysadmin := m.RoleGrant(meta.UserCtxKey, map[string]bool{
		roles.SystemAdmin: true,
	})

	r := chi.NewRouter()

	r.Route("/cities-svc/", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Get("/city/{slug}", h.GetCityBySlug)

			r.Route("/cities", func(r chi.Router) {
				r.Get("/", h.ListCities)

				r.With(auth, sysadmin).Post("/", h.CreateCity)

				r.Route("/{city_id}", func(r chi.Router) {
					r.Get("/", h.GetCity)

					r.With(auth).Put("/", h.UpdateCity)
					r.With(auth, sysadmin).Patch("/status", h.UpdateCityStatus)

					r.Route("/admins", func(r chi.Router) {
						r.Get("/", h.ListAdmins)

						r.With(auth).Route("/invite", func(r chi.Router) {
							r.Post("/", h.SentInvite)
							r.Post("/", h.AnswerInvite)
						})

						r.With(auth).Route("/me", func(r chi.Router) {
							r.Get("/", h.GetMyCityAdmin)
							r.Put("/", h.UpdateMyCityAdmin)
							r.Delete("/", h.RefuseMyCityAdmin)
						})

						r.Route("/{user_id}", func(r chi.Router) {
							r.Get("/", h.GetCityAdmin)
							r.With(auth).Put("/", h.UpdateCityAdmin)
							r.With(auth).Delete("/", h.DeleteCityAdmin)
						})
					})
				})
			})
		})
	})

	srv := &http.Server{
		Addr:              cfg.Rest.Port,
		Handler:           r,
		ReadTimeout:       cfg.Rest.Timeouts.Read,
		ReadHeaderTimeout: cfg.Rest.Timeouts.ReadHeader,
		WriteTimeout:      cfg.Rest.Timeouts.Write,
		IdleTimeout:       cfg.Rest.Timeouts.Idle,
	}

	log.Infof("starting REST service on %s", cfg.Rest.Port)

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		} else {
			errCh <- nil
		}
	}()

	select {
	case <-ctx.Done():
		log.Info("shutting down REST service...")
	case err := <-errCh:
		if err != nil {
			log.Errorf("REST server error: %v", err)
		}
	}

	shCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shCtx); err != nil {
		log.Errorf("REST shutdown error: %v", err)
	} else {
		log.Info("REST server stopped")
	}
}
