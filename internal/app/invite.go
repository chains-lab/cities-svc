package app

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/app/domain/gov"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
)

type CreateInviteParams struct {
	InitiatorID uuid.UUID
	CityID      uuid.UUID
	Role        string
}

func (a App) CreateInvite(ctx context.Context, params CreateInviteParams) (models.Invite, error) {
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

	newInvite, err := a.gov.CreateInvite(ctx, p)
	if err != nil {
		return models.Invite{}, err
	}

	return newInvite, nil
}

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

func (a App) AcceptInvite(ctx context.Context, initiatorID uuid.UUID, token string) (models.Invite, error) {
	var invite models.Invite
	var err error

	txErr := a.transaction(func(ctx context.Context) error {
		invite, err = a.gov.AcceptInvite(ctx, initiatorID, token)
		if err != nil {
			return err
		}

		city, err := a.GetCityByID(ctx, invite.CityID)
		if err != nil {
			return err
		}
		if city.Status != enum.CityStatusOfficial {
			return errx.ErrorAnswerToInviteForNotOffSupCity.Raise(
				fmt.Errorf("cannot answer to invite for not official support city %s", city.ID),
			)
		}

		return nil
	})
	if txErr != nil {
		return models.Invite{}, txErr
	}

	return invite, err
}
