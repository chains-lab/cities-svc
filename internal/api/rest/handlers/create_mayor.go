package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/api/rest/meta"
	"github.com/chains-lab/cities-svc/internal/api/rest/requests"
	"github.com/chains-lab/cities-svc/internal/api/rest/responses"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/google/uuid"
)

func (a Adapter) CreateMayor(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.Log(r).WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.CreateMayor(r)
	if err != nil {
		a.Log(r).WithError(err).Error("failed to parse create mayor request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	userID, err := uuid.Parse(req.Data.Attributes.UserId)
	if err != nil {
		a.Log(r).WithError(err).Error("invalid user ID format")
		ape.RenderErr(w, problems.InvalidPointer("/data/attributes/user_id", err))

		return
	}

	cityID, err := uuid.Parse(req.Data.Attributes.CityId)
	if err != nil {
		a.Log(r).WithError(err).Error("invalid city ID format")
		ape.RenderErr(w, problems.InvalidPointer("/data/attributes/city_id", err))

		return
	}

	mayor, err := a.app.CreateGovMayor(r.Context(), cityID, userID)
	if err != nil {
		a.Log(r).WithError(err).Error("failed to create mayor")
		switch {
		case errors.Is(err, errx.ErrorGovAlreadyExists):
			ape.RenderErr(w, problems.Conflict("a mayor for this city already exists"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	a.log.Infof("mayor %s created for city %s by initiator %s", userID, cityID, initiator.ID)

	ape.Render(w, http.StatusCreated, responses.Gov(mayor))
}
