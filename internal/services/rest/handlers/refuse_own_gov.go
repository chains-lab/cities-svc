package handlers

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/services/rest/meta"
)

func (a Adapter) RefuseOwnCurrentGov(w http.ResponseWriter, r *http.Request) {
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
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}
}
