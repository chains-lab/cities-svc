package controller

import (
	"net/http"
	"strings"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/api/rest/responses"
	"github.com/chains-lab/cities-svc/internal/domain/services/country"

	"github.com/chains-lab/pagi"
)

func (a Service) ListCountries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	q := r.URL.Query()

	filters := country.FilterParams{}

	if name := strings.TrimSpace(q.Get("name")); name != "" {
		filters.Name = &name
	}

	if sts := q["status"]; len(sts) > 0 {
		filters.Statuses = sts
	}

	page, size := pagi.GetPagination(r)

	countries, err := a.domain.country.Filter(ctx, filters, page, size)
	if err != nil {
		a.log.WithError(err).Error("failed to search countries")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.CountriesCollection(countries))
}
