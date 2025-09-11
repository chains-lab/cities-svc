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
	"github.com/paulmach/orb"
)

func (a Adapter) UpdateCity(w http.ResponseWriter, r *http.Request) {
	req, err := requests.UpdateCity(r)
	if err != nil {
		a.Log(r).WithError(err).Error("failed to parse update city request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	if req.Data.Id != chi.URLParam(r, "city_id") {
		a.Log(r).Error("body id does not match url city_id")
		ape.RenderErr(w,
			problems.InvalidParameter("city_id", fmt.Errorf("data/id does not match url city_id")),
			problems.InvalidPointer("/data/id", fmt.Errorf("data/id does not match url city_id")),
		)
		return
	}

	cityID, err := uuid.Parse(req.Data.Id)
	if err != nil {
		a.Log(r).WithError(err).Error("invalid city_id")
		ape.RenderErr(w, problems.InvalidParameter("city_id", err))

		return
	}

	param := app.UpdateCityParams{}

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

	city, err := a.app.UpdateCity(r.Context(), cityID, param)
	if err != nil {
		a.Log(r).WithError(err).Error("failed to update city")
		switch {
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

	ape.Render(w, http.StatusOK, responses.City(city))
}
