package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/rest/meta"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func (s Service) RefuseMyCityAdmin(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		http.Error(w, "failed to get user from context", http.StatusUnauthorized)

		return
	}

	cityID, err := uuid.Parse(chi.URLParam(r, "city_id"))
	if err != nil {
		s.log.WithError(err).Error("failed to parse city_id param")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"city_id": err,
		})...)

		return
	}

	err = s.domain.admin.DeleteOwn(r.Context(), initiator.ID, cityID)
	if err != nil {
		s.log.WithError(err).Error("failed to refuse own admin")
		switch {
		case errors.Is(err, errx.ErrorNotEnoughRight):
			ape.RenderErr(w, problems.Forbidden("no active city admin for the user"))
		case errors.Is(err, errx.ErrorCityAdminTechLeadCannotRefuseOwn):
			ape.RenderErr(w, problems.Forbidden("tech lead cannot refuse own admin"))

		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	s.log.Infof("user %s refused own admin successfully", initiator.ID)

	w.WriteHeader(http.StatusNoContent)
}
