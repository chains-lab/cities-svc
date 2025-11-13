package controller

import (
	"errors"
	"net/http"
	"time"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/rest/meta"
	"github.com/chains-lab/cities-svc/internal/rest/requests"
	"github.com/chains-lab/cities-svc/internal/rest/responses"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
)

func (a Service) SentInvite(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.CreateInvite(r)
	if err != nil {
		a.log.WithError(err).Error("failed to parse create city moder request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	CityID, err := uuid.Parse(chi.URLParam(r, "city_id"))
	if err != nil {
		a.log.WithError(err).Error("invalid city_id")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"city_id": err,
		})...)

		return
	}

	inv, err := a.domain.invite.Create(r.Context(), CityID, req.Data.Attributes.UserId, req.Data.Attributes.Role, 24*time.Hour)
	if err != nil {
		a.log.WithError(err).Error("failed to create city moder")
		switch {
		case errors.Is(err, errx.ErrorInvalidCityAdminRole):
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"role": err,
			})...)
		case errors.Is(err, errx.ErrorCityIsNotSupported):
			ape.RenderErr(w, problems.Forbidden("cannot create invite for not official city"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	a.log.Infof("admin %s created successfully by user %s", inv.ID, initiator.ID)

	ape.Render(w, http.StatusCreated, responses.Invite(inv))
}
