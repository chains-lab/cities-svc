package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/api/rest/meta"
	"github.com/chains-lab/cities-svc/internal/api/rest/responses"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
)

func (a Adapter) InviteMayor(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.Log(r).WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	cityID, err := uuid.Parse(r.URL.Query().Get("city_id"))
	if err != nil {
		a.Log(r).WithError(err).Errorf("invalid city ID format: %s", r.URL.Query().Get("city_id"))
		ape.RenderErr(w, problems.InvalidParameter("city_id", err))

		return
	}

	access := false
	if initiator.Role == roles.Admin || initiator.Role == roles.SuperUser {
		access = true
	} else {
		gov, err := a.app.GetInitiatorGov(r.Context(), initiator.ID)
		if err != nil {
			a.Log(r).WithError(err).Error("failed to get initiator gov from context")
			switch {
			case errors.Is(err, errx.ErrorInitiatorIsNotActiveCityGov):
				ape.RenderErr(w, problems.Forbidden("initiator is not active city government"))
			default:
				ape.RenderErr(w, problems.InternalError())
			}

			return
		}

		if gov.CityID == cityID && gov.Role == enum.CityGovRoleMayor {
			access = true
		}
	}

	if !access {
		a.Log(r).Warnf("user %s with role %s attempted to invite mayor for city %s", initiator.ID, initiator.Role, cityID)
		ape.RenderErr(w, problems.PreconditionFailed("only admin, superuser or mayor of the city can invite new mayor"))

		return
	}

	inv, err := a.app.CreateInviteMayor(r.Context(), cityID)
	if err != nil {
		a.Log(r).WithError(err).Errorf("failed to create invite for city %s", cityID)
		switch {
		case errors.Is(err, errx.ErrorCityNotFound):
			ape.RenderErr(w, problems.NotFound("city not found"))
		case errors.Is(err, errx.ErrorMayorInviteAlreadyExists):
			ape.RenderErr(w, problems.Conflict("mayor invite already exists"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	a.Log(r).Infof("mayor invite %s created successfully for city %s by user %s", inv.ID, cityID, initiator.ID)

	ape.Render(w, http.StatusCreated, responses.Invite(inv))
}
