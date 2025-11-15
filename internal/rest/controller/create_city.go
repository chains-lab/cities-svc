package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/services/city"
	"github.com/chains-lab/cities-svc/internal/rest/meta"
	"github.com/chains-lab/cities-svc/internal/rest/requests"
	"github.com/chains-lab/cities-svc/internal/rest/responses"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/paulmach/orb"
)

func (s Service) CreateCity(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.CreateCity(r)
	if err != nil {
		s.log.WithError(err).Error("error creating city")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	c, err := s.domain.city.Create(r.Context(), city.CreateParams{
		Name:      req.Data.Attributes.Name,
		CountryID: req.Data.Attributes.CountryId,
		Status:    req.Data.Attributes.Status,
		Point: orb.Point{
			req.Data.Attributes.Point.Longitude,
			req.Data.Attributes.Point.Latitude,
		},
		Timezone: req.Data.Attributes.Timezone,
	})
	if err != nil {
		s.log.WithError(err).Error("error creating city")
		switch {
		case errors.Is(err, errx.ErrorInvalidTimeZone):
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"data/attributes/timezone": err,
			})...)
		case errors.Is(err, errx.ErrorInvalidPoint):
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"data/attributes/point": err,
			})...)
		case errors.Is(err, errx.ErrorInvalidCityName):
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"data/attributes/name": err,
			})...)
		case errors.Is(err, errx.ErrorInvalidCountryISO3ID):
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"data/attributes/country_id": err,
			})...)
		case errors.Is(err, errx.ErrorInvalidCityStatus):
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"data/attributes/status": err,
			})...)

		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	s.log.Infof("created city with name %s by user %s", c.Name, initiator.ID)

	ape.Render(w, http.StatusCreated, responses.City(c))
}
