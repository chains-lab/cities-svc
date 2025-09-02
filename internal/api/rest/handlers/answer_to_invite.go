package handlers

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/api/rest/meta"
	"github.com/chains-lab/cities-svc/internal/api/rest/requests"
	"github.com/chains-lab/cities-svc/internal/api/rest/responses"
	"github.com/go-chi/chi/v5"
)

func (a Adapter) AnswerToInvite(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.Log(r).WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.AnswerToInvite(r)
	if err != nil {
		a.Log(r).WithError(err).Error("failed to decode answer to invite request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	token := chi.URLParam(r, "token")

	invite, err := a.app.AnswerToInvite(r.Context(), initiator.ID, req.Data.Attributes.Answer, token)
	if err != nil {
		a.Log(r).WithError(err).Error("failed to answer to invite")
		ape.RenderErr(w, problems.InternalError())

		return
	}

	ape.Render(w, http.StatusCreated, responses.Invite(invite))
}
