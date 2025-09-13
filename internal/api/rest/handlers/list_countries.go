package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/api/rest/responses"
	"github.com/chains-lab/cities-svc/internal/app"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/pagi"
)

func (a Adapter) ListCountries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	q := r.URL.Query()

	filters := app.FilterCountriesListParams{}

	if name := strings.TrimSpace(q.Get("name")); name != "" {
		filters.Name = name
	}

	if sts := q["status"]; len(sts) > 0 {
		filters.Statuses = sts
	}

	var pag pagi.Request
	if v := q.Get("page"); v != "" {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil && n > 0 {
			pag.Page = n
		}
	}

	if v := q.Get("size"); v != "" {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil && n > 0 {
			pag.Size = n
		}
	}

	var sort []pagi.SortField
	if raw := q.Get("sort"); raw != "" {
		fields := strings.Split(raw, ",")
		for _, f := range fields {
			f = strings.TrimSpace(f)
			if f == "" {
				continue
			}
			asc := true
			if strings.HasPrefix(f, "-") {
				asc = false
				f = f[1:]
			}
			sort = append(sort, pagi.SortField{Field: f, Ascend: asc})
		}
	}

	countries, resp, err := a.app.ListCountries(ctx, filters, pag, sort)
	if err != nil {
		a.log.WithError(err).Error("failed to search countries")
		switch {
		case errors.Is(err, errx.ErrorInvalidCountryStatus):
			ape.RenderErr(w, problems.InvalidParameter("status", err))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.CountriesCollection(countries, resp))
}
