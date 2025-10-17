package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/rest/meta"
	"github.com/chains-lab/cities-svc/internal/rest/responses"
)

func (a Service) GetMyCityAdmin(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	res, err := a.domain.moder.GetInitiator(r.Context(), initiator.ID)
	if err != nil {
		a.log.WithError(err).Error("failed to get own active admin")

		switch {
		case errors.Is(err, errx.ErrorInitiatorIsNotCityAdmin):
			ape.RenderErr(w, problems.NotFound("no active city government for the user"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.CityAdmin(res))
}
