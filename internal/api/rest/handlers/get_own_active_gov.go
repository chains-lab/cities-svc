package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/api/rest/meta"
	"github.com/chains-lab/cities-svc/internal/api/rest/responses"
	"github.com/chains-lab/cities-svc/internal/errx"
)

func (a Adapter) GetOwnGov(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	gov, err := a.app.GetInitiatorGov(r.Context(), initiator.ID)
	if err != nil {
		a.log.WithError(err).Error("failed to get own active gov")

		switch {
		case errors.Is(err, errx.ErrorInitiatorIsNotActiveCityGov):
			ape.RenderErr(w, problems.NotFound("no active city government for the user"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.Gov(gov))
}
