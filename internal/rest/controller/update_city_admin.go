package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/cities-svc/internal/domain/services/admin"
	"github.com/chains-lab/cities-svc/internal/rest/meta"
	"github.com/chains-lab/cities-svc/internal/rest/requests"
	"github.com/chains-lab/cities-svc/internal/rest/responses"
	"github.com/chains-lab/restkit/roles"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func (s Service) UpdateCityAdmin(w http.ResponseWriter, r *http.Request) {
	req, err := requests.UpdateCityAmin(r)
	if err != nil {
		s.log.WithError(err).Error("failed to parse update city admin request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

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
	if userID.String() != chi.URLParam(r, "user_id") {
		s.log.Error("user_id in url and user_id in body do not match")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"id": fmt.Errorf("user_id in url and user_id in body do not match"),
		})...)

		return
	}

	var result models.CityAdmin
	switch initiator.Role {
	case roles.SystemUser:
		result, err = s.domain.admin.UpdateByCityAdmin(r.Context(), initiator.ID, userID, cityID, admin.UpdateParams{
			Label:    req.Data.Attributes.Label,
			Position: req.Data.Attributes.Position,
		})
	default:
		result, err = s.domain.admin.UpdateBySysAdmin(r.Context(), userID, cityID, admin.UpdateParams{
			Role:     req.Data.Attributes.Role,
			Label:    req.Data.Attributes.Label,
			Position: req.Data.Attributes.Position,
		})
	}
	if err != nil {
		s.log.WithError(err).Error("failed to update city admin")
		ape.RenderErr(w, problems.InternalError())

		return
	}

	ape.Render(w, http.StatusAccepted, responses.CityAdmin(result))
}
