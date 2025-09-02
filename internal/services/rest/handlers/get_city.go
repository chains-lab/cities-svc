package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/cities-svc/internal/services/rest/responses"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (a Adapter) GetCity(w http.ResponseWriter, r *http.Request) {
	cityID, err := uuid.Parse(chi.URLParam(r, "city_id"))
	if err != nil {
		a.Log(r).WithError(err).Error("invalid city_id")
		ape.RenderErr(w, problems.InvalidParameter("city_id", err))

		return
	}

	city, err := a.app.GetCityByID(r.Context(), cityID)
	if err != nil {
		a.Log(r).WithError(err).Error("failed to get city")
		switch {
		case errors.Is(err, errx.ErrorCityNotFound):
			ape.RenderErr(w, problems.NotFound("city not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.City(city))
}
