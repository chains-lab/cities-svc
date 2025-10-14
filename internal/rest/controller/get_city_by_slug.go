package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/rest/responses"
	"github.com/go-chi/chi/v5"
)

func (a Service) GetCityBySlug(w http.ResponseWriter, r *http.Request) {
	city, err := a.domain.city.GetBySlug(r.Context(), chi.URLParam(r, "slug"))
	if err != nil {
		a.log.WithError(err).Error("failed to get city")
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
