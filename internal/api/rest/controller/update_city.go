package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/api/rest/requests"
	"github.com/chains-lab/cities-svc/internal/api/rest/responses"
	"github.com/chains-lab/cities-svc/internal/domain/services/city"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/paulmach/orb"
)

func (a Service) UpdateCity(w http.ResponseWriter, r *http.Request) {
	//initiator, err := meta.User(r.Context())
	//if err != nil {
	//	a.log.WithError(err).Error("failed to get user from context")
	//	ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))
	//
	//	return
	//}

	req, err := requests.UpdateCity(r)
	if err != nil {
		a.log.WithError(err).Error("failed to parse update city request")
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

	res, err := a.domain.city.Update(r.Context(), req.Data.Id, param)
	if err != nil {
		a.log.WithError(err).Error("failed to update city")
		switch {
		case errors.Is(err, errx.ErrorInitiatorIsNotActiveCityGov):
			ape.RenderErr(w, problems.Forbidden("initiator is not an active city governor"))
		case errors.Is(err, errx.ErrorInitiatorIsNotThisCityGov):
			ape.RenderErr(w, problems.Forbidden("initiator is not the city governor"))
		case errors.Is(err, errx.ErrorInitiatorGovRoleHaveNotEnoughRights):
			ape.RenderErr(w, problems.Forbidden("initiator governor role have not enough rights"))
		case errors.Is(err, errx.ErrorCityNotFound):
			ape.RenderErr(w, problems.NotFound("city not found"))
		case errors.Is(err, errx.ErrorInvalidPoint):
			ape.RenderErr(w, problems.InvalidPointer("data/attributes/point", err))
		case errors.Is(err, errx.ErrorInvalidTimeZone):
			ape.RenderErr(w, problems.InvalidPointer("data/attributes/timezone", err))
		case errors.Is(err, errx.ErrorCityAlreadyExistsWithThisSlug):
			ape.RenderErr(w, problems.Conflict("city with the given slug already exists"))
		case errors.Is(err, errx.ErrorInvalidCityStatus):
			ape.RenderErr(w, problems.InvalidPointer("data/attributes/status", err))
		case errors.Is(err, errx.ErrorInvalidCityName):
			ape.RenderErr(w, problems.InvalidPointer("data/attributes/name", err))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.City(res))
}
