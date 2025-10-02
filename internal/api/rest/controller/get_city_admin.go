package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/api/rest/responses"
	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/services/admin"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func (a Service) GetCityAdmin(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		a.log.WithError(err).Error("invalid user_id")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"user_id": err,
		})...)

		return
	}

	res, err := a.domain.moder.Get(r.Context(), admin.GetFilters{UserID: &userID})
	if err != nil {
		a.log.WithError(err).Error("failed to get admin")
		switch {
		case errors.Is(err, errx.ErrorCityAdminNotFound):
			ape.RenderErr(w, problems.NotFound("city government not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.CityAdmin(res))
}
