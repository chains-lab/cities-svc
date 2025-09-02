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
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (a Adapter) UpdateGov(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.Log(r).WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.UpdateGov(r)
	if err != nil {
		a.Log(r).WithError(err).Error("failed to parse update gov request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	if req.Data.Id != chi.URLParam(r, "gov_id") {
		a.Log(r).Error("body id does not match url gov_id")
		ape.RenderErr(w,
			problems.InvalidParameter("gov_id", fmt.Errorf("data/id does not match url gov_id")),
			problems.InvalidPointer("/data/id", fmt.Errorf("data/id does not match url gov_id")),
		)

		return
	}

	govID, err := uuid.Parse(chi.URLParam(r, "gov_id"))
	if err != nil {
		a.Log(r).WithError(err).Error("invalid gov_id")
		ape.RenderErr(w, problems.InvalidParameter("gov_id", err))

		return
	}

	params := app.UpdateGovParams{}
	if params.Label != nil {
		params.Label = req.Data.Attributes.Label
	}
	if params.Active != nil {
		params.Active = req.Data.Attributes.Active
	}

	gov, err := a.app.UpdateGov(r.Context(), initiator.ID, govID, params)
	if err != nil {
		a.Log(r).WithError(err).Error("failed to update gov")
		switch {

		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	a.log.Infof("city government member %s updated successfully", govID)

	ape.Render(w, http.StatusOK, responses.Gov(gov))
}
