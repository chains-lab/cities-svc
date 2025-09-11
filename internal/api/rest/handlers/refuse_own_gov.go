package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/api/rest/meta"
	"github.com/chains-lab/cities-svc/internal/errx"
)

func (a Adapter) RefuseOwnGov(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.Log(r).WithError(err).Error("failed to get user from context")
		http.Error(w, "failed to get user from context", http.StatusUnauthorized)

		return
	}

	err = a.app.RefuseOwnGov(r.Context(), initiator.ID)
	if err != nil {
		a.Log(r).WithError(err).Error("failed to refuse own gov")
		switch {
		case errors.Is(err, errx.ErrorCannotRefuseMayor):
			ape.RenderErr(w, problems.Forbidden("not enough rights to refuse own gov"))
		case errors.Is(err, errx.ErrorInitiatorIsNotActiveCityGov):
			ape.RenderErr(w, problems.PreconditionFailed("no active city government for the user"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}
}
