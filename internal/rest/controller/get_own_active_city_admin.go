package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/rest/meta"
	"github.com/chains-lab/cities-svc/internal/rest/responses"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func (s Service) GetMyCityAdmin(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	cityID, err := uuid.Parse(r.URL.Query().Get("city_id"))
	if err != nil {
		s.log.WithError(err).Error("invalid city_id parameter")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"city_id": errors.New("invalid city_id parameter"),
		})...)
		return
	}

	res, err := s.domain.admin.Get(r.Context(), initiator.ID, cityID)
	if err != nil {
		s.log.WithError(err).Error("failed to get own active admin")

		switch {
		case errors.Is(err, errx.ErrorCityAdminNotFound):
			ape.RenderErr(w, problems.Unauthorized("no active city admin for the user"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.CityAdmin(res))
}
