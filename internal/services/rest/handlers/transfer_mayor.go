package handlers

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/services/rest/meta"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (a Adapter) TransferMayor(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.Log(r).WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		a.Log(r).WithError(err).Error("invalid user_id")
		ape.RenderErr(w, problems.InvalidParameter("user_id", err))

		return
	}

	err = a.app.TransferGovMayor(r.Context(), initiator.ID, userID)
	if err != nil {
		a.Log(r).WithError(err).Error("failed to transfer mayor")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}
}
