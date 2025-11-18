package controller

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/domain/services/admin"
	"github.com/chains-lab/cities-svc/internal/rest/responses"
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/chains-lab/restkit/pagi"
	"github.com/google/uuid"
)

func (s Service) ListAdmins(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query()

	var filters admin.FilterParams

	if userIDs := q["user_id"]; len(userIDs) > 0 {
		ids := make([]uuid.UUID, 0, len(userIDs))
		for _, idStr := range userIDs {
			id, err := uuid.Parse(idStr)
			if err != nil {
				ape.RenderErr(w, problems.BadRequest(
					validation.Errors{
						"user_id": err,
					})...)
				return
			}
			ids = append(ids, id)
		}

		filters.UserID = ids
	}

	if cityIDs := q["city_id"]; len(cityIDs) > 0 {
		ids := make([]uuid.UUID, 0, len(cityIDs))
		for _, idStr := range cityIDs {
			id, err := uuid.Parse(idStr)
			if err != nil {
				ape.RenderErr(w, problems.BadRequest(
					validation.Errors{
						"city_id": err,
					})...)
				return
			}
			ids = append(ids, id)
		}

		filters.CityID = ids
	}

	if role := q["role"]; len(role) > 0 {
		roles := make([]string, 0, len(role))
		for _, r := range role {
			roles = append(roles, r)
		}

		filters.Roles = roles
	}

	page, size := pagi.GetPagination(r)

	admins, err := s.domain.admin.Filter(ctx, filters, page, size)
	if err != nil {
		s.log.WithError(err).Error("failed to search admins")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.CityAdminsCollection(admins))
}
