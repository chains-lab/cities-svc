package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/api/rest/meta"
	"github.com/chains-lab/cities-svc/internal/api/rest/responses"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/google/uuid"
)

func (a Adapter) InviteMayor(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	cityID, err := uuid.Parse(r.URL.Query().Get("city_id"))
	if err != nil {
		a.log.WithError(err).Errorf("invalid city ID format: %s", r.URL.Query().Get("city_id"))
		ape.RenderErr(w, problems.InvalidParameter("city_id", err))

		return
	}

	inv, err := a.app.CreateInviteMayor(r.Context(), cityID, initiator.ID, initiator.Role)
	if err != nil {
		a.log.WithError(err).Errorf("failed to create invite for city %s", cityID)
		switch {
		case errors.Is(err, errx.ErrorInitiatorIsNotActiveCityGov):
			ape.RenderErr(w, problems.Forbidden("initiator is not an active city governor"))
		case errors.Is(err, errx.ErrorInitiatorIsNotThisCityGov):
			ape.RenderErr(w, problems.Forbidden("initiator is not the city governor"))
		case errors.Is(err, errx.ErrorInitiatorGovRoleHaveNotEnoughRights):
			ape.RenderErr(w, problems.Forbidden("initiator governor role have not enough rights"))
		case errors.Is(err, errx.ErrorCannotCreateMayorInviteForNotOfficialCity):
			ape.RenderErr(w, problems.Conflict("cannot create mayor invite for not official city"))
		case errors.Is(err, errx.ErrorCityNotFound):
			ape.RenderErr(w, problems.NotFound("city not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	a.log.Infof("mayor invite %s created successfully for city %s by user %s", inv.ID, cityID, initiator.ID)

	ape.Render(w, http.StatusCreated, responses.Invite(inv))
}
