package controller

import (
	"errors"
	"net/http"
	"time"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/api/rest/meta"
	"github.com/chains-lab/cities-svc/internal/api/rest/requests"
	"github.com/chains-lab/cities-svc/internal/api/rest/responses"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
)

func (a Service) CreateInvite(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.CreateInvite(r)
	if err != nil {
		a.log.WithError(err).Error("failed to parse create city moder request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	CityID, err := uuid.Parse(chi.URLParam(r, "city_id"))
	if err != nil {
		a.log.WithError(err).Error("invalid city_id")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"city_id": err,
		})...)

		return
	}

	//TODO WHAt do with token?
	inv, token, err := a.domain.moder.CreateInvite(r.Context(), req.Data.Attributes.Role, CityID, 24*time.Hour)
	if err != nil {
		a.log.WithError(err).Error("failed to create city moder")
		switch {
		case errors.Is(err, errx.ErrorCannotCreateInviteForNotOfficialCity):
			ape.RenderErr(w, problems.Conflict("cannot create invite for not official city"))
		case errors.Is(err, errx.ErrorInitiatorIsNotThisCityGov):
			ape.RenderErr(w, problems.Forbidden("only city city moder can create invite"))
		case errors.Is(err, errx.ErrorInitiatorGovRoleHaveNotEnoughRights):
			ape.RenderErr(w, problems.NotFound("initiator role have not enough rights to invite this role"))
		case errors.Is(err, errx.ErrorInvalidGovRole):
			ape.RenderErr(w, problems.InvalidPointer("/data/attributes/role", err))
		case errors.Is(err, errx.ErrorGovAlreadyExists):
			ape.RenderErr(w, problems.Conflict("city moder already exists for this user"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	a.log.Infof("citymod %s created successfully by user %s", inv.ID, initiator.ID)

	ape.Render(w, http.StatusCreated, responses.Invite(inv, token))
}
