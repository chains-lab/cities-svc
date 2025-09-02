package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/api/rest/responses"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/go-chi/chi/v5"
)

func (a Adapter) GetCityBySlug(w http.ResponseWriter, r *http.Request) {
	city, err := a.app.GetCityBySlug(r.Context(), chi.URLParam(r, "slug"))
	if err != nil {
		a.Log(r).WithError(err).Error("failed to get city")
		switch {
		case errors.Is(err, errx.ErrorCityNotFound):
			ape.RenderErr(w, problems.NotFound("city not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.City(city))
}
