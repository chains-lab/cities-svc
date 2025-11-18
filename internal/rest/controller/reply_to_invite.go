package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/rest/meta"
	"github.com/chains-lab/cities-svc/internal/rest/requests"
	"github.com/chains-lab/cities-svc/internal/rest/responses"
)

func (s Service) ReplyInvite(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.ReplyToInvite(r)
	if err != nil {
		s.log.WithError(err).Error("invalid answer invite request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	res, err := s.domain.invite.Reply(r.Context(), initiator.ID, req.Data.Id, req.Data.Attributes.Answer)
	if err != nil {
		s.log.WithError(err).Error("failed to answer to invite")
		switch {
		case errors.Is(err, errx.ErrorInviteNotFound):
			ape.RenderErr(w, problems.NotFound("invite not found"))
		case errors.Is(err, errx.ErrorInviteAlreadyReplied):
			ape.RenderErr(w, problems.Conflict("invite already answered"))
		case errors.Is(err, errx.ErrorInviteExpired):
			ape.RenderErr(w, problems.Conflict("invite expired"))
		case errors.Is(err, errx.ErrorCityAdminAlreadyExists):
			ape.RenderErr(w, problems.Conflict("user is already s city admin"))
		case errors.Is(err, errx.ErrorCityIsNotSupported):
			ape.RenderErr(w, problems.Forbidden("cannot accept invite for not official support city"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusCreated, responses.Invite(res))
}
