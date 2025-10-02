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
	validation "github.com/go-ozzo/ozzo-validation/v4"
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
		case errors.Is(err, errx.ErrorCityNotFound):
			ape.RenderErr(w, problems.NotFound("city not found"))
		case errors.Is(err, errx.ErrorCountryIsNotSupported):
			ape.RenderErr(w, problems.Conflict(
				fmt.Sprintf("cannot update city status %s in unsupported country", req.Data.Attributes.Status)),
			)
		case errors.Is(err, errx.ErrorInvalidCityStatus):
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"status": fmt.Errorf("status is not supported %s", err),
			})...)

		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.City(res))
}
