package app

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/app/entities/gov"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

type SentInviteParams struct {
	InitiatorID uuid.UUID
	CityID      uuid.UUID
	Role        string
}

func (a App) SentInvite(ctx context.Context, params SentInviteParams) (models.Invite, error) {
	city, err := a.GetCityByID(ctx, params.CityID)
	if err != nil {
		return models.Invite{}, err
	}

	if city.Status != enum.CityStatusOfficial {
		return models.Invite{}, errx.ErrorCannotCreateInviteForNotOfficialCity.Raise(
			fmt.Errorf("invite for city %s is not official", params.CityID),
		)
	}

	p := gov.SentInviteParams{
		InitiatorID: params.InitiatorID,
		CityID:      params.CityID,
		Role:        params.Role,
	}

	newInvite, err := a.gov.SentInvite(ctx, p)
	if err != nil {
		return models.Invite{}, err
	}

	return newInvite, nil
}
