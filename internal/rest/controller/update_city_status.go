package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/rest/meta"
	"github.com/chains-lab/cities-svc/internal/rest/requests"
	"github.com/chains-lab/cities-svc/internal/rest/responses"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (s Service) UpdateCityStatus(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.UpdateCityStatus(r)
	if err != nil {
		s.log.WithError(err).Error("failed to parse update city request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}
	if req.Data.Id.String() != chi.URLParam(r, "city_id") {
		s.log.Error("city_id in url and city_id in body do not match")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"id": fmt.Errorf("city_id in url and city_id in body do not match"),
		})...)

		return
	}

	res, err := s.domain.city.UpdateStatusByCityAdmin(r.Context(), initiator.ID, req.Data.Id, req.Data.Attributes.Status)
	if err != nil {
		s.log.WithError(err).Error("failed to update city status")
		switch {
		case errors.Is(err, errx.ErrorCityNotFound):
			ape.RenderErr(w, problems.NotFound("city not found"))
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
