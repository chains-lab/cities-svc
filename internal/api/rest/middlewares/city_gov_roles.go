package middlewares

import (
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/domain/services/citymod"
	"github.com/chains-lab/gatekit/auth"
	"github.com/chains-lab/gatekit/roles"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func (m Middleware) CityGovRoles(UserCtxKey interface{}, allowedGovRoles map[string]bool, allowedSysadminRoles map[string]bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			user, ok := ctx.Value(UserCtxKey).(auth.UserData)
			if !ok {
				ape.RenderErr(w,
					problems.Unauthorized("Missing AuthorizationHeader header"),
				)

				return
			}

			if err := roles.ParseRole(user.Role); err != nil {
				ape.RenderErr(w,
					problems.Unauthorized("User role not valid"),
				)

				return
			}

			cityID, err := uuid.Parse(chi.URLParam(r, "city_id"))
			if err != nil {
				ape.RenderErr(w,
					problems.BadRequest(validation.Errors{
						"city_id": fmt.Errorf("city id is invalid format"),
					})...)
				return
			}

			if allowedSysadminRoles[user.Role] {
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			initiator, err := m.domain.gov.Get(ctx, citymod.GetFilters{
				CityID: &cityID,
				UserID: &user.ID,
			})
			if err != nil {
				m.log.WithError(err).Errorf("failed to get gov for user %s and city %s", user.ID, cityID)
				ape.RenderErr(w,
					problems.Unauthorized("User is not a city moderator"),
				)

				return
			}

			if !allowedGovRoles[initiator.Role] {
				ape.RenderErr(w,
					problems.Forbidden("user moderation role is not allowed"),
				)

				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
