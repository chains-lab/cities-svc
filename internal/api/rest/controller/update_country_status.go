package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/api/rest/requests"
	"github.com/chains-lab/cities-svc/internal/api/rest/responses"
	"github.com/chains-lab/cities-svc/internal/domain/errx"
)

func (a Service) UpdateCountryStatus(w http.ResponseWriter, r *http.Request) {
	req, err := requests.UpdateCountryStatus(r)
	if err != nil {
		a.log.WithError(err).Error("failed to parse update country status request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	res, err := a.domain.country.UpdateStatus(r.Context(), req.Data.Id, req.Data.Attributes.Status)
	if err != nil {
		a.log.WithError(err).Error("failed to update country status")
		switch {
		case errors.Is(err, errx.ErrorCountryNotFound):
			ape.RenderErr(w, problems.NotFound("country not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Country(res))
}
