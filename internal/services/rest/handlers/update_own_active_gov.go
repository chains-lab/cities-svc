package handlers

import (
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/app"
	"github.com/chains-lab/cities-svc/internal/services/rest/meta"
	"github.com/chains-lab/cities-svc/internal/services/rest/requests"
	"github.com/chains-lab/cities-svc/internal/services/rest/responses"
)

func (a Adapter) UpdateOwnCurrenteGov(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.Log(r).WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.UpdateOwnGov(r)
	if err != nil {
		a.Log(r).WithError(err).Error("failed to parse update own active gov request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	if initiator.ID.String() != req.Data.Id {
		a.Log(r).Error("user ID does not match request ID")
		ape.RenderErr(w, problems.InvalidPointer("/data/id",
			fmt.Errorf("user ID does not match request ID")),
		)

		return
	}

	params := app.UpdateOwnActiveGovParams{}
	if req.Data.Attributes.Label != nil {
		params.Label = *req.Data.Attributes.Label
	}

	gov, err := a.app.UpdateOwnActiveGov(r.Context(), initiator.ID, params)
	if err != nil {
		a.Log(r).WithError(err).Error("failed to update own active gov")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.Gov(gov))
}
