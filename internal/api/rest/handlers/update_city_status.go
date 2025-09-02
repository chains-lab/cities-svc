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
	"github.com/chains-lab/cities-svc/internal/constant"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (a Adapter) UpdateCityStatus(w http.ResponseWriter, r *http.Request) {
	req, err := requests.UpdateCityStatus(r)
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

	var city models.City
	switch req.Data.Attributes.Status {
	case constant.CityStatusOfficial:
		city, err = a.app.SetCityStatusOfficial(r.Context(), cityID)
	case constant.CityStatusCommunity:
		city, err = a.app.SetCityStatusCommunity(r.Context(), cityID)
	case constant.CityStatusDeprecated:
		city, err = a.app.SetCityStatusDeprecated(r.Context(), cityID)
	default:
		a.Log(r).Error("invalid city status")
		ape.RenderErr(w, problems.InvalidPointer("data/attributes/status",
			fmt.Errorf("invalid city status for update, allowed values are: %s, %s, %s",
				constant.CityStatusOfficial, constant.CityStatusCommunity, constant.CityStatusDeprecated),
		),
		)
	}

	if err != nil {
		a.Log(r).WithError(err).Error("failed to update city status")
		switch {
		case errors.Is(err, errx.ErrorInvalidCityStatus):
			ape.RenderErr(w, problems.InvalidPointer("data/attributes/status", fmt.Errorf("invalid city status")))
		case errors.Is(err, errx.ErrorCityNotFound):
			ape.RenderErr(w, problems.NotFound("city not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.City(city))
}
