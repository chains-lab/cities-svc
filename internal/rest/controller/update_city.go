package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/cities-svc/internal/domain/services/city"
	"github.com/chains-lab/cities-svc/internal/rest/meta"
	"github.com/chains-lab/cities-svc/internal/rest/requests"
	"github.com/chains-lab/cities-svc/internal/rest/responses"
	"github.com/chains-lab/restkit/roles"
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/paulmach/orb"
)

func (s Service) UpdateCity(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.UpdateCity(r)
	if err != nil {
		s.log.WithError(err).Error("failed to parse update city request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	param := city.UpdateParams{}

	if req.Data.Attributes.Name != nil {
		param.Name = req.Data.Attributes.Name
	}
	if req.Data.Attributes.Point != nil {
		param.Point = &orb.Point{
			req.Data.Attributes.Point.Longitude,
			req.Data.Attributes.Point.Latitude,
		}
	}
	if req.Data.Attributes.Timezone != nil {
		param.Timezone = req.Data.Attributes.Timezone
	}
	if req.Data.Attributes.Icon != nil {
		param.Icon = req.Data.Attributes.Icon
	}
	if req.Data.Attributes.Slug != nil {
		param.Slug = req.Data.Attributes.Slug
	}

	var res models.City
	switch initiator.Role {
	case roles.SystemUser:
		res, err = s.domain.city.UpdateByCityAdmin(r.Context(), initiator.ID, req.Data.Id, param)
	default:
		res, err = s.domain.city.UpdateByAdmin(r.Context(), req.Data.Id, param)
	}
	if err != nil {
		s.log.WithError(err).Error("failed to update city")
		switch {
		case errors.Is(err, errx.ErrorNotEnoughRight):
			ape.RenderErr(w, problems.Forbidden("not enough rights to update city"))
		case errors.Is(err, errx.ErrorCityNotFound):
			ape.RenderErr(w, problems.NotFound("city not found"))
		case errors.Is(err, errx.ErrorInvalidPoint):
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"data/attributes/point": err,
			})...)
		case errors.Is(err, errx.ErrorInvalidTimeZone):
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"data/attributes/timezone": err,
			})...)
		case errors.Is(err, errx.ErrorInvalidCityStatus):
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"data/attributes/status": err,
			})...)
		case errors.Is(err, errx.ErrorInvalidCityName):
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"data/attributes/name": err,
			})...)
		case errors.Is(err, errx.ErrorInvalidSlug):
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"data/attributes/slug": err,
			})...)

		case errors.Is(err, errx.ErrorCityAlreadyExistsWithThisSlug):
			ape.RenderErr(w, problems.Conflict("city with the given slug already exists"))

		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	s.log.Infof("city %s updated by user %s", res.ID, initiator.ID)

	ape.Render(w, http.StatusOK, responses.City(res))
}
