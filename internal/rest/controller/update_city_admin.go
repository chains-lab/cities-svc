package controller

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/cities-svc/internal/domain/services/admin"
	"github.com/chains-lab/cities-svc/internal/rest/meta"
	"github.com/chains-lab/cities-svc/internal/rest/requests"
	"github.com/chains-lab/cities-svc/internal/rest/responses"
	"github.com/chains-lab/restkit/roles"
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

	var result models.CityAdmin
	switch initiator.Role {
	case roles.SystemUser:
		result, err = s.domain.admin.Update(r.Context(), req.Data.Id, initiator.ID, admin.UpdateParams{
			Label:    req.Data.Attributes.Label,
			Position: req.Data.Attributes.Position,
		})
	default:
		result, err = s.domain.admin.UpdateBySysAdmin(r.Context(), req.Data.Id, admin.UpdateParams{
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
