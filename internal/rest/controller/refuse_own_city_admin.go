package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/rest/meta"
)

func (a Service) RefuseMyCityAdmin(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.log.WithError(err).Error("failed to get user from context")
		http.Error(w, "failed to get user from context", http.StatusUnauthorized)

		return
	}

	err = a.domain.moder.RefuseOwn(r.Context(), initiator.ID)
	if err != nil {
		a.log.WithError(err).Error("failed to refuse own admin")
		switch {
		case errors.Is(err, errx.ErrorInitiatorIsNotCityAdmin):
			ape.RenderErr(w, problems.Forbidden("no active city adminernment for the user"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	a.log.Infof("user %s refused own admin successfully", initiator.ID)

	w.WriteHeader(http.StatusNoContent)
}
