package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/cities-svc/internal/services/rest/responses"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (a Adapter) GetCountry(w http.ResponseWriter, r *http.Request) {
	countryID, err := uuid.Parse(chi.URLParam(r, "country_id"))
	if err != nil {
		a.log.WithError(err).Error("invalid country_id")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	country, err := a.app.GetCountryByID(r.Context(), countryID)
	if err != nil {
		a.log.WithError(err).Error("failed to get country")
		switch {
		case errors.Is(err, errx.ErrorCountryNotFound):
			ape.RenderErr(w, problems.NotFound("country not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Country(country))
}
