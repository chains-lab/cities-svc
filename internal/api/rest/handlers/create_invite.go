package handlers

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/api/rest/meta"
	"github.com/chains-lab/cities-svc/internal/api/rest/requests"
	"github.com/chains-lab/cities-svc/internal/api/rest/responses"
	"github.com/google/uuid"
)

func (a Adapter) CreateInvite(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.Log(r).WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.CreateInvite(r)
	if err != nil {
		a.Log(r).WithError(err).Error("failed to parse create gov request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	userID, err := uuid.Parse(req.Data.Attributes.UserId)
	if err != nil {
		a.Log(r).WithError(err).Error("invalid user ID format")
		ape.RenderErr(w, problems.InvalidPointer("/data/attributes/user_id", err))

		return
	}

	//TODO
	gov, _, err := a.app.CreateInvite(r.Context(), initiator.ID, userID, req.Data.Attributes.Role)
	if err != nil {
		a.Log(r).WithError(err).Error("failed to create gov")
		switch {

		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	a.Log(r).Infof("gov %s created successfully by user %s", gov.ID, initiator.ID)

	ape.Render(w, http.StatusCreated, responses.Invite(gov))
}
