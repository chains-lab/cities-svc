package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/app"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/cities-svc/internal/services/rest/meta"
	"github.com/chains-lab/cities-svc/internal/services/rest/requests"
	"github.com/chains-lab/cities-svc/internal/services/rest/responses"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

func (a Adapter) CreateCity(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.CreateCity(r)
	if err != nil {
		a.log.WithError(err).Error("error creating city")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	countryID, err := uuid.Parse(req.Data.Attributes.CountryId)
	if err != nil {
		a.log.WithError(err).Error("invalid country_id")
		ape.RenderErr(w, problems.InvalidPointer("data/attributes/country_id", err))

		return
	}

	point := orb.Point{
		req.Data.Attributes.Point.Longitude,
		req.Data.Attributes.Point.Latitude,
	}

	city, err := a.app.CreateCity(r.Context(), app.CreateCityParams{
		Name:      req.Data.Attributes.Name,
		CountryID: countryID,
		Point:     point,
		Timezone:  req.Data.Attributes.Timezone,
	})
	if err != nil {
		a.log.WithError(err).Error("error creating city")
		switch {
		case errors.Is(err, errx.ErrorInvalidTimeZone):
			ape.RenderErr(w, problems.InvalidPointer("data/attributes/timezone", err))
		case errors.Is(err, errx.ErrorInvalidPoint):
			ape.RenderErr(w, problems.InvalidPointer("data/attributes/point", err))
		case errors.Is(err, errx.ErrorInvalidCityName):
			ape.RenderErr(w, problems.InvalidPointer("data/attributes/name", err))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	a.log.Infof("created city with name %s by user %s", city.Name, initiator.ID)

	ape.Render(w, http.StatusCreated, responses.City(city))
}
