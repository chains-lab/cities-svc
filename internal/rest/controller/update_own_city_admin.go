package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/domain/services/admin"
	"github.com/chains-lab/cities-svc/internal/rest/meta"
	"github.com/chains-lab/cities-svc/internal/rest/requests"
	"github.com/chains-lab/cities-svc/internal/rest/responses"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
)

func (s Service) UpdateMyCityAdmin(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.UpdateOwnCityAdmin(r)
	if err != nil {
		s.log.WithError(err).Error("failed to parse update own active admin request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	ids := strings.Split(req.Data.Id, ":")
	if len(ids) != 2 {
		s.log.Error("invalid city admin id format")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"id": fmt.Errorf("invalid id: %s, need format uuid:uuid look like user_id:city_id", req.Data.Id),
		})...)

		return
	}
	userID, err := uuid.Parse(ids[0])
	if err != nil {
		s.log.WithError(err).Error("failed to parse user id from city admin id")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"id": fmt.Errorf("invalid user id in city admin id: %s", ids[0]),
		})...)

		return
	}
	cityID, err := uuid.Parse(ids[1])
	if err != nil {
		s.log.WithError(err).Error("failed to parse city id from city admin id")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"id": fmt.Errorf("invalid city id in city admin id: %s", ids[1]),
		})...)

		return
	}
	if cityID.String() != chi.URLParam(r, "city_id") {
		s.log.Error("city_id in url and city_id in body do not match")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"id": fmt.Errorf("city_id in url and city_id in body do not match"),
		})...)

		return
	}
	if initiator.ID.String() != userID.String() {
		s.log.Error("user ID does not match request ID")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"id": errors.New("user ID does not match request ID"),
		})...)
		return
	}

	params := admin.UpdateOwnParams{
		Label:    req.Data.Attributes.Label,
		Position: req.Data.Attributes.Position,
	}

	res, err := s.domain.admin.UpdateOwn(r.Context(), userID, cityID, params)
	if err != nil {
		s.log.WithError(err).Error("failed to update own active admin")
		switch {
		case errors.Is(err, errx.ErrorNotEnoughRight):
			ape.RenderErr(w, problems.Forbidden("only active city admin can update their admin info"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.CityAdmin(res))
}
