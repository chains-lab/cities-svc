package handlers

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/app"
	"github.com/chains-lab/cities-svc/internal/services/rest/meta"
	"github.com/chains-lab/cities-svc/internal/services/rest/requests"
	"github.com/chains-lab/cities-svc/internal/services/rest/responses"
	"github.com/google/uuid"
)

func (a Adapter) CreateGov(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.Log(r).WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.CreateGov(r)
	if err != nil {
		a.Log(r).WithError(err).Error("failed to parse create gov request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	cityID, err := uuid.Parse(req.Data.Attributes.CityId)
	if err != nil {
		a.Log(r).WithError(err).Error("invalid city ID format")
		ape.RenderErr(w, problems.InvalidPointer("/data/attributes/city_id", err))

		return
	}

	userID, err := uuid.Parse(req.Data.Attributes.UserId)
	if err != nil {
		a.Log(r).WithError(err).Error("invalid user ID format")
		ape.RenderErr(w, problems.InvalidPointer("/data/attributes/user_id", err))

		return
	}

	gov, err := a.app.CreateGov(r.Context(), initiator.ID, app.CreateGovParams{
		CityID: cityID,
		UserID: userID,
		Label:  req.Data.Attributes.Label,
		Role:   req.Data.Attributes.Role,
	})
	if err != nil {
		a.Log(r).WithError(err).Error("failed to create gov")
		switch {

		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	a.log.Infof("gov %s created for city %s by initiator %s", gov.ID, cityID, initiator.ID)

	ape.Render(w, http.StatusCreated, responses.Gov(gov))
}
