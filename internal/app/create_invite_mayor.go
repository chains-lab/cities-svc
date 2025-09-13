package app

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
)

func (a App) CreateInviteMayor(ctx context.Context, cityID, initiatorID uuid.UUID, role string) (models.Invite, error) {
	if role == roles.User {
		gov, err := a.GetInitiatorGov(ctx, initiatorID)
		if err != nil {
			return models.Invite{}, err
		}
		if gov.CityID != cityID {
			return models.Invite{}, errx.ErrorInitiatorIsNotThisCityGov.Raise(
				fmt.Errorf("initiator %s is not the city %s", initiatorID.String(), cityID.String()),
			)
		}
		if gov.Role != enum.CityGovRoleMayor {
			return models.Invite{}, errx.ErrorInitiatorGovRoleHaveNotEnoughRights.Raise(
				fmt.Errorf("only mayor can invite new mayor"),
			)
		}
	}

	city, err := a.cities.GetByID(ctx, cityID)
	if err != nil {
		return models.Invite{}, err
	}

	if city.Status != enum.CityStatusOfficial {
		return models.Invite{}, errx.ErrorCannotCreateMayorInviteForNotOfficialCity.Raise(
			fmt.Errorf("cannot create mayor invite for city with status %s", city.Status),
		)
	}
	newInvite, err := a.gov.CreateMayorInvite(ctx, cityID)
	if err != nil {
		return models.Invite{}, err
	}

	return newInvite, nil
}
