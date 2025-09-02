package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/api/rest/responses"
	"github.com/chains-lab/cities-svc/internal/app"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
)

func (a Adapter) SearchGovs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query()

	var filters app.SearchGovsFilters

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

	var pag pagi.Request
	if v := q.Get("page"); v != "" {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil && n > 0 {
			pag.Page = n
		} else {
			ape.RenderErr(w, problems.InvalidParameter("page", fmt.Errorf("must be positive integer")))
			return
		}
	}
	if v := q.Get("size"); v != "" {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil && n > 0 {
			pag.Size = n
		} else {
			ape.RenderErr(w, problems.InvalidParameter("size", fmt.Errorf("must be positive integer")))
			return
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

	govs, resp, err := a.app.SearchGovs(ctx, filters, pag, sort)
	if err != nil {
		a.Log(r).WithError(err).Error("failed to search govs")

		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.GovsCollection(govs, resp))
}
