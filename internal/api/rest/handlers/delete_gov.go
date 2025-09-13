package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/api/rest/meta"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (a Adapter) DeleteGov(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	cityID, err := uuid.Parse(chi.URLParam(r, "city_id"))
	if err != nil {
		a.log.WithError(err).Error("invalid city_id")
		ape.RenderErr(w, problems.InvalidParameter("city_id", err))

		return
	}

	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		a.log.WithError(err).Error("invalid user_id")
		ape.RenderErr(w, problems.InvalidParameter("user_id", err))

		return
	}

	err = a.app.DeleteGov(r.Context(), initiator.ID, userID, cityID)
	if err != nil {
		a.log.WithError(err).Error("failed to delete gov")
		switch {
		case errors.Is(err, errx.ErrorInitiatorAndUserHaveDifferentCity):
			ape.RenderErr(w, problems.Conflict("initiator and user have different city"))
		case errors.Is(err, errx.ErrorInitiatorIsNotThisCityGov):
			ape.RenderErr(w, problems.Conflict("initiator is not this city gov"))
		case errors.Is(err, errx.ErrorCityGovNotFound):
			ape.RenderErr(w, problems.NotFound("city gov not found"))
		case errors.Is(err, errx.ErrorInitiatorIsNotActiveCityGov):
			ape.RenderErr(w, problems.Conflict("initiator is not active city gov"))
		case errors.Is(err, errx.ErrorInitiatorGovRoleHaveNotEnoughRights):
			ape.RenderErr(w, problems.Conflict("initiator role have not enough rights"))
		case errors.Is(err, errx.ErrorCannotRefuseMayor):
			ape.RenderErr(w, problems.Conflict("cannot delete mayor, assign new mayor first"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}
