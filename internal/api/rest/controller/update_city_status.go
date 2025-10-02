package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/api/rest/requests"
	"github.com/chains-lab/cities-svc/internal/api/rest/responses"
	"github.com/chains-lab/cities-svc/internal/domain/errx"
)

func (a Service) UpdateCityStatus(w http.ResponseWriter, r *http.Request) {
	req, err := requests.UpdateCityStatus(r)
	if err != nil {
		a.log.WithError(err).Error("failed to parse update city request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	res, err := a.domain.city.UpdateStatus(r.Context(), req.Data.Id, req.Data.Attributes.Status)
	if err != nil {
		a.log.WithError(err).Error("failed to update city status")
		switch {
		case errors.Is(err, errx.ErrorCannotUpdateCityStatusInUnsupportedCountry):
			ape.RenderErr(w, problems.Conflict("cannot update city status in unsupported country"))
		case errors.Is(err, errx.ErrorInvalidCityStatus):
			ape.RenderErr(w, problems.InvalidPointer("data/attributes/status", fmt.Errorf("invalid city status")))
		case errors.Is(err, errx.ErrorCityNotFound):
			ape.RenderErr(w, problems.NotFound("city not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.City(res))
}
