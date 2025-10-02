package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/cities-svc/internal/api/rest/meta"
	"github.com/chains-lab/cities-svc/internal/api/rest/requests"
	"github.com/chains-lab/cities-svc/internal/api/rest/responses"
	"github.com/chains-lab/cities-svc/internal/domain/services/citymod"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
)

func (a Service) UpdateOwnGov(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.UpdateOwnGov(r)
	if err != nil {
		a.log.WithError(err).Error("failed to parse update own active citymod request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	if initiator.ID.String() != req.Data.Id.String() {
		a.log.Error("user ID does not match request ID")
		ape.RenderErr(w, problems.InvalidPointer("/data/id",
			fmt.Errorf("user ID does not match request ID")),
		)

		return
	}

	params := citymod.UpdateCityModerParams{}
	if req.Data.Attributes.Label != nil {
		params.Label = req.Data.Attributes.Label
	}

	res, err := a.domain.moder.UpdateOwn(r.Context(), initiator.ID, params)
	if err != nil {
		a.log.WithError(err).Error("failed to update own active citymod")
		switch {
		case errors.Is(err, errx.ErrorInitiatorIsNotActiveCityGov):
			ape.RenderErr(w, problems.Forbidden("only active city citymod can update their citymod info"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.Gov(res))
}
