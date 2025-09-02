package handlers

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/services/rest/meta"
	"github.com/chains-lab/cities-svc/internal/services/rest/responses"
)

func (a Adapter) GetOwnCurrentGov(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.Log(r).WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	gov, err := a.app.GetOwnActiveGov(r.Context(), initiator.ID)
	if err != nil {
		a.Log(r).WithError(err).Error("failed to get own active gov")

		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.Gov(gov))
}
