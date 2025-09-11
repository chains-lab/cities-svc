package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/api/rest/meta"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (a Adapter) DeleteGov(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.Log(r).WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	gov, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		a.Log(r).WithError(err).Error("invalid user_id")
		ape.RenderErr(w, problems.InvalidParameter("user_id", err))

		return
	}

	err = a.app.DeleteGov(r.Context(), initiator.ID, gov)
	if err != nil {
		a.Log(r).WithError(err).Error("failed to delete gov")
		switch {
		case errors.Is(err, errx.ErrorInitiatorRoleHaveNotEnoughRights):
			ape.RenderErr(w, problems.PreconditionFailed("initiator role have not enough rights"))
		case errors.Is(err, errx.ErrorCannotRefuseMayor):
			ape.RenderErr(w, problems.PreconditionFailed("cannot delete mayor, assign new mayor first"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}
