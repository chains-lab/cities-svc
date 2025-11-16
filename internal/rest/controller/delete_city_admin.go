package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/rest/meta"
	"github.com/chains-lab/restkit/roles"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func (s Service) DeleteCityAdmin(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	cityID, err := uuid.Parse(chi.URLParam(r, "city_id"))
	if err != nil {
		s.log.WithError(err).Error("invalid city_id")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"city_id": err,
		})...)

		return
	}

	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		s.log.WithError(err).Error("invalid user_id")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"user_id": err,
		})...)

		return
	}

	switch initiator.Role {
	case roles.SystemUser:
		err = s.domain.admin.DeleteByCityAdmin(r.Context(), cityID, userID, initiator.ID)
	default:
		err = s.domain.admin.DeleteBySysAdmin(r.Context(), userID, cityID)
	}

	if err != nil {
		s.log.WithError(err).Error("failed to delete admin")
		switch {
		case errors.Is(err, errx.ErrorCityAdminNotFound):
			ape.RenderErr(w, problems.NotFound("city admin not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}
