package controller

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func (a Service) DeleteCityAdmin(w http.ResponseWriter, r *http.Request) {
	cityID, err := uuid.Parse(chi.URLParam(r, "city_id"))
	if err != nil {
		a.log.WithError(err).Error("invalid city_id")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"city_id": err,
		})...)

		return
	}

	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		a.log.WithError(err).Error("invalid user_id")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"user_id": err,
		})...)

		return
	}

	err = a.domain.moder.Delete(r.Context(), userID, cityID)
	if err != nil {
		a.log.WithError(err).Error("failed to delete admin")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}
