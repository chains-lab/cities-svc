package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/api/rest/meta"
	"github.com/chains-lab/cities-svc/internal/api/rest/requests"
	"github.com/chains-lab/cities-svc/internal/api/rest/responses"
	"github.com/chains-lab/cities-svc/internal/app"
	"github.com/chains-lab/cities-svc/internal/errx"
)

func (a Adapter) CreateInvite(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.CreateInvite(r)
	if err != nil {
		a.log.WithError(err).Error("failed to parse create gov request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	inv, err := a.app.CreateInvite(r.Context(), app.CreateInviteParams{
		InitiatorID: initiator.ID,
		Role:        req.Data.Attributes.Role,
	})
	if err != nil {
		a.log.WithError(err).Error("failed to create gov")
		switch {
		case errors.Is(err, errx.ErrorCannotCreateInviteForNotOfficialCity):
			ape.RenderErr(w, problems.Conflict("cannot create invite for not official city"))
		case errors.Is(err, errx.ErrorInitiatorIsNotThisCityGov):
			ape.RenderErr(w, problems.Forbidden("only city gov can create invite"))
		case errors.Is(err, errx.ErrorInitiatorGovRoleHaveNotEnoughRights):
			ape.RenderErr(w, problems.NotFound("initiator role have not enough rights to invite this role"))
		case errors.Is(err, errx.ErrorInvalidGovRole):
			ape.RenderErr(w, problems.InvalidPointer("/data/attributes/role", err))
		case errors.Is(err, errx.ErrorGovAlreadyExists):
			ape.RenderErr(w, problems.Conflict("gov already exists for this user"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	a.log.Infof("gov %s created successfully by user %s", inv.ID, initiator.ID)

	ape.Render(w, http.StatusCreated, responses.Invite(inv))
}
