package middlewares

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/cities-svc/internal/services/rest/meta"
	"github.com/chains-lab/gatekit/auth"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (m Adapter) Govs(govRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			user, ok := ctx.Value(meta.UserCtxKey).(auth.UserData)
			if !ok {
				m.log.Error("missing user data in context")
				ape.RenderErr(w, problems.Unauthorized("Missing Authorization header"))

				return
			}

			cityID, err := uuid.Parse(chi.URLParam(r, "city_id"))
			if err != nil {
				m.log.WithError(err).Error("invalid city_id")
				ape.RenderErr(w, problems.InvalidParameter("city_id", err))

				return
			}

			initiator, err := m.app.GetInitiator(ctx, cityID, user.ID)
			if err != nil {
				m.log.WithError(err).Error("failed to get initiator")
				switch {
				case errors.Is(err, errx.ErrorNotActiveCityGovInitiator):
					ape.RenderErr(w, problems.Forbidden("The initiator must be an active city government member"))
				default:
					ape.RenderErr(w, problems.InternalError())
				}

				return
			}

			access := false
			for _, role := range govRoles {
				if initiator.Role == role {
					access = true
					break
				}
			}

			if !access {
				m.log.Errorf("the initiator %s  for city %s dosent have enought permisions", user.ID, cityID)
				ape.RenderErr(w, problems.Forbidden("initiator does not have enough permissions"))

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
