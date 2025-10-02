package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/api/rest/responses"
	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/services/citymod"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (a Service) GetGov(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		a.log.WithError(err).Error("invalid user_id")
		ape.RenderErr(w, problems.InvalidParameter("user_id", err))

		return
	}

	res, err := a.domain.moder.Get(r.Context(), citymod.GetFilters{UserID: &userID})
	if err != nil {
		a.log.WithError(err).Error("failed to get citymod")
		switch {
		case errors.Is(err, errx.ErrorCityGovNotFound):
			ape.RenderErr(w, problems.NotFound("city government not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Gov(res))
}
