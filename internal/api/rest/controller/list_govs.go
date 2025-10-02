package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/api/rest/responses"
	"github.com/chains-lab/cities-svc/internal/domain/services/citymod"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
)

func (a Service) ListGovs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query()

	var filters citymod.FilterParams

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

	govs, err := a.domain.moder.Filter(ctx, filters, page, size)
	if err != nil {
		a.log.WithError(err).Error("failed to search govs")
		switch {
		case errors.Is(err, errx.ErrorInvalidGovRole):
			ape.RenderErr(w, problems.InvalidParameter("role", err))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.GovsCollection(govs))
}
