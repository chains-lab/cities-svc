package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/api/rest/meta"
	"github.com/chains-lab/cities-svc/internal/api/rest/requests"
	"github.com/chains-lab/cities-svc/internal/api/rest/responses"
	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/services/city"
	"github.com/paulmach/orb"
)

func (a Service) CreateCity(w http.ResponseWriter, r *http.Request) {
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

	c, err := a.domain.city.Create(r.Context(), city.CreateParams{
		Name:      req.Data.Attributes.Name,
		CountryID: req.Data.Attributes.CountryId,
		Point: orb.Point{
			req.Data.Attributes.Point.Longitude,
			req.Data.Attributes.Point.Latitude,
		},
		Timezone: req.Data.Attributes.Timezone,
	})
	if err != nil {
		a.log.WithError(err).Error("error creating city")
		switch {
		case errors.Is(err, errx.ErrorCountryNotSupported):
			ape.RenderErr(w, problems.Conflict("cannot create city in unsupported country"))
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

	a.log.Infof("created city with name %s by user %s", c.Name, initiator.ID)

	ape.Render(w, http.StatusCreated, responses.City(c))
}
