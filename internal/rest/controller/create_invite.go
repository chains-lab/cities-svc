package controller

import (
	"errors"
	"net/http"
	"time"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/cities-svc/internal/rest/meta"
	"github.com/chains-lab/cities-svc/internal/rest/requests"
	"github.com/chains-lab/cities-svc/internal/rest/responses"
	"github.com/chains-lab/restkit/roles"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (s Service) SentInvite(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.CreateInvite(r)
	if err != nil {
		s.log.WithError(err).Error("failed to parse create city admin request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	var result models.Invite
	switch initiator.Role {
	case roles.SystemUser:
		result, err = s.domain.invite.CreateByCityAdmin(r.Context(),
			req.Data.Attributes.UserId,
			req.Data.Attributes.CityId,
			initiator.ID,
			req.Data.Attributes.Role,
			24*time.Hour,
		)
	default:
		result, err = s.domain.invite.CreateBySysAdmin(
			r.Context(),
			req.Data.Attributes.UserId,
			req.Data.Attributes.CityId,
			req.Data.Attributes.Role,
			24*time.Hour,
		)
	}
	if err != nil {
		s.log.WithError(err).Error("failed to create city admin")
		switch {
		case errors.Is(err, errx.ErrorInitiatorHasNoRights):
			ape.RenderErr(w, problems.Conflict("initiator have no rights for this action"))
		case errors.Is(err, errx.ErrorCityAdminAlreadyExists):
			ape.RenderErr(w, problems.Conflict("city admin already exists"))
		case errors.Is(err, errx.ErrorInvalidCityAdminRole):
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"role": err,
			})...)
		case errors.Is(err, errx.ErrorCityIsNotSupported):
			ape.RenderErr(w, problems.Forbidden("cannot create invite for not official city"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	s.log.Infof("admin %s created successfully by user %s", result.ID, initiator.ID)

	ape.Render(w, http.StatusCreated, responses.Invite(result))
}
