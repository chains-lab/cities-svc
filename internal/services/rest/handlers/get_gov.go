package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/cities-svc/internal/services/rest/responses"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (a Adapter) GetGov(w http.ResponseWriter, r *http.Request) {
	govID, err := uuid.Parse(chi.URLParam(r, "gov_id"))
	if err != nil {
		a.Log(r).WithError(err).Error("invalid gov_id")
		ape.RenderErr(w, problems.InvalidParameter("gov_id", err))

		return
	}

	gov, err := a.app.GetGov(r.Context(), govID)
	if err != nil {
		a.Log(r).WithError(err).Error("failed to get gov")
		switch {
		case errors.Is(err, errx.ErrorCityGovNotFound):
			ape.RenderErr(w, problems.NotFound("city government not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Gov(gov))
}
