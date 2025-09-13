package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/api/rest/requests"
	"github.com/chains-lab/cities-svc/internal/api/rest/responses"
	"github.com/chains-lab/cities-svc/internal/app"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (a Adapter) UpdateCountry(w http.ResponseWriter, r *http.Request) {
	req, err := requests.UpdateCountry(r)
	if err != nil {
		a.log.WithError(err).Error("failed to parse update country request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	if req.Data.Id != chi.URLParam(r, "country_id") {
		a.log.Error("body id does not match url country_id")
		ape.RenderErr(w,
			problems.InvalidParameter("country_id", fmt.Errorf("data/id does not match url country_id")),
			problems.InvalidPointer("/data/id", fmt.Errorf("data/id does not match url country_id")),
		)

		return
	}

	countryID, err := uuid.Parse(req.Data.Id)
	if err != nil {
		a.log.WithError(err).Error("invalid country_id")
		ape.RenderErr(w, problems.InvalidParameter("country_id", err))

		return
	}

	params := app.UpdateCountryParams{}
	if req.Data.Attributes.Name != nil {
		params.Name = req.Data.Attributes.Name
	}

	country, err := a.app.UpdateCountry(r.Context(), countryID, params)
	if err != nil {
		a.log.WithError(err).Error("failed to update country")
		switch {
		case errors.Is(err, errx.ErrorCountryNotFound):
			ape.RenderErr(w, problems.NotFound("country not found"))
		case errors.Is(err, errx.ErrorCountryAlreadyExistsWithThisName):
			ape.RenderErr(w, problems.Conflict("country with the same name already exists"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Country(country))
}
