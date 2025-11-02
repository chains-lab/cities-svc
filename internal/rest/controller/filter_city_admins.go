package controller

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/domain/services/admin"
	"github.com/chains-lab/cities-svc/internal/rest/responses"

	"github.com/chains-lab/restkit/pagi"
	"github.com/google/uuid"
)

func (a Service) ListAdmins(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query()

	var filters admin.FilterParams

	if cityID, err := uuid.Parse(q.Get("city_id")); err != nil {
		filters.CityID = &cityID
	}

	if role := q["role"]; len(role) > 0 {
		roles := make([]string, 0, len(role))
		for _, r := range role {
			roles = append(roles, r)
		}

		filters.Roles = roles
	}

	page, size := pagi.GetPagination(r)

	admins, err := a.domain.moder.Filter(ctx, filters, page, size)
	if err != nil {
		a.log.WithError(err).Error("failed to search admins")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.CityAdminsCollection(admins))
}
