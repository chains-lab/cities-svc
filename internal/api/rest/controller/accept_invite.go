package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/api/rest/meta"
	"github.com/chains-lab/cities-svc/internal/api/rest/responses"
	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/go-chi/chi/v5"
)

func (a Service) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	token := chi.URLParam(r, "token")

	res, err := a.domain.invite.Accept(r.Context(), initiator.ID, token)
	if err != nil {
		a.log.WithError(err).Error("failed to answer to invite")
		switch {
		case errors.Is(err, errx.ErrorInvalidInviteToken):
			ape.RenderErr(w, problems.Unauthorized("invalid invite token"))
		case errors.Is(err, errx.ErrorInviteNotFound):
			ape.RenderErr(w, problems.NotFound("invite not found"))
		case errors.Is(err, errx.ErrorInviteAlreadyAnswered):
			ape.RenderErr(w, problems.Conflict("invite already answered"))
		case errors.Is(err, errx.ErrorInviteExpired):
			ape.RenderErr(w, problems.Conflict("invite expired"))
		case errors.Is(err, errx.ErrorUserIsAlreadyCityAdmin):
			ape.RenderErr(w, problems.Conflict("user is already a city admin"))
		case errors.Is(err, errx.ErrorCityIsNotSupported):
			ape.RenderErr(w, problems.Conflict("cannot accept invite for not official support city"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusCreated, responses.Invite(res))
}
