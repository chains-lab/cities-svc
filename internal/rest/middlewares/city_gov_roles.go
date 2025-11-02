package middlewares

import (
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/restkit/roles"
	"github.com/chains-lab/restkit/token"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func (s Service) CityAdminMember(
	userCtxKey interface{},
	AllowedAdminRoles map[string]bool,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(userCtxKey).(token.UserData)
			if !ok {
				ape.RenderErr(w, problems.Unauthorized("Missing AuthorizationHeader header"))
				return
			}

			if err := roles.ParseRole(user.Role); err != nil {
				ape.RenderErr(w, problems.Unauthorized("User role not valid"))
				return
			}

			cityID, err := uuid.Parse(chi.URLParam(r, "city_id"))
			if err != nil {
				ape.RenderErr(w, problems.BadRequest(validation.Errors{
					"city_id": fmt.Errorf("city id is invalid format"),
				})...)
				return
			}

			if user.CityID == nil || user.CityRole == nil {
				ape.RenderErr(w, problems.Forbidden("user is not admin of this city"))
				return
			}

			if *user.CityID != cityID {
				ape.RenderErr(w, problems.Forbidden("user city ID does not match"))
				return
			}

			if !AllowedAdminRoles[*user.CityRole] {
				ape.RenderErr(w, problems.Forbidden("user city role not allowed"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
