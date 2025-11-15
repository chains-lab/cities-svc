package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/domain/services/admin"
	"github.com/chains-lab/cities-svc/internal/rest/meta"
	"github.com/chains-lab/cities-svc/internal/rest/requests"
	"github.com/chains-lab/cities-svc/internal/rest/responses"
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
)

func (s Service) UpdateMyCityAdmin(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.UpdateOwnAdmin(r)
	if err != nil {
		s.log.WithError(err).Error("failed to parse update own active admin request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	if initiator.ID.String() != req.Data.Id.String() {
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

	res, err := s.domain.admin.UpdateOwn(r.Context(), initiator.ID, params)
	if err != nil {
		s.log.WithError(err).Error("failed to update own active admin")
		switch {
		case errors.Is(err, errx.ErrorInitiatorIsNotCityAdmin):
			ape.RenderErr(w, problems.Forbidden("only active city admin can update their admin info"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.CityAdmin(res))
}
