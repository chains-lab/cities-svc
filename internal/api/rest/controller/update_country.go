package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/api/rest/requests"
	"github.com/chains-lab/cities-svc/internal/api/rest/responses"
	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/services/country"
)

func (a Service) UpdateCountry(w http.ResponseWriter, r *http.Request) {
	req, err := requests.UpdateCountry(r)
	if err != nil {
		a.log.WithError(err).Error("failed to parse update country request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	params := country.UpdateParams{}
	if req.Data.Attributes.Name != nil {
		params.Name = req.Data.Attributes.Name
	}

	res, err := a.domain.country.Update(r.Context(), req.Data.Id, params)
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

	ape.Render(w, http.StatusOK, responses.Country(res))
}
