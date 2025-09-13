package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/api/rest/requests"
	"github.com/chains-lab/cities-svc/internal/api/rest/responses"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (a Adapter) UpdateCountryStatus(w http.ResponseWriter, r *http.Request) {
	req, err := requests.UpdateCountryStatus(r)
	if err != nil {
		a.log.WithError(err).Error("failed to parse update country status request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	if req.Data.Id != chi.URLParam(r, "country_id") {
		a.log.Error("body id does not match url country_id")
		ape.RenderErr(w,
			problems.InvalidParameter("country_id", fmt.Errorf("data/id does not match url country_id")),
			problems.InvalidPointer("/data/id", fmt.Errorf("data/id does not match url country_id")),
		)

		return
	}

	countryID, err := uuid.Parse(req.Data.Id)
	if err != nil {
		a.log.WithError(err).Error("invalid country_id")
		ape.RenderErr(w, problems.InvalidParameter("country_id", err))

		return
	}

	var country models.Country

	switch req.Data.Attributes.Status {
	case enum.CountryStatusSupported:
		country, err = a.app.SetCountryStatusSupported(r.Context(), countryID)
	case enum.CountryStatusDeprecated:
		country, err = a.app.SetCountryStatusDeprecated(r.Context(), countryID)
	default:
		a.log.Error("invalid country status")
		ape.RenderErr(w, problems.InvalidPointer("data/attributes/status",
			fmt.Errorf("invalid country status for update, allowed values are: %s, %s",
				enum.CountryStatusSupported, enum.CountryStatusDeprecated),
		),
		)

		return
	}

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

	ape.Render(w, http.StatusOK, responses.Country(country))
}
