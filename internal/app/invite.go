package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/app/entities"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/constant"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/google/uuid"
)

func (a App) CreateInvite(ctx context.Context, initiatorID, userID uuid.UUID, role string) (models.Invite, string, error) {
	initiator, err := a.GetInitiatorGov(ctx, initiatorID)
	if err != nil {
		return models.Invite{}, "", err
	}

	err = constant.CheckCityGovRole(role)
	if err != nil {
		return models.Invite{}, "", errx.ErrorInvalidGovRole.Raise(
			fmt.Errorf("check city gov role: %w", err),
		)
	}

	_, err = a.gov.Get(ctx, entities.GetGovFilters{
		CityID: &initiator.CityID,
		UserID: &userID,
		Role:   &role,
	})
	if err != nil && !errors.Is(err, errx.ErrorCityGovNotFound) {
		return models.Invite{}, "", err
	}
	if err == nil {
		return models.Invite{}, "", errx.ErrorGovAlreadyExists.Raise(
			fmt.Errorf("user %s already has role %s in city %s", userID, role, initiator.CityID),
		)
	}

	access, err := constant.CompareCityGovRoles(initiator.Role, role)
	if err != nil {
		return models.Invite{}, "", errx.ErrorInvalidGovRole.Raise(
			fmt.Errorf("compare city gov roles: %w", err),
		)
	}
	if access != 0 {
		return models.Invite{}, "", errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("no access to invite user %s with role %s, initiator role is %s", userID, role, initiator.Role),
		)
	}

	newInvite, token, err := a.invite.Create(ctx, entities.CreateInviteParams{
		InitiatorID: initiatorID,
		CityID:      initiator.CityID,
		Role:        role,
		TimeLife:    24 * time.Hour,
	})
	if err != nil {
		return models.Invite{}, "", fmt.Errorf("create invite: %w", err)
	}

	return newInvite, token, nil
}

func (a App) AnswerToInvite(ctx context.Context, initiatorID uuid.UUID, status, token string) (models.Invite, error) {
	var invite models.Invite
	var err error

	txErr := a.transaction(func(ctx context.Context) error {
		invite, err = a.invite.Answered(ctx, initiatorID, token, status)
		if err != nil {
			return err
		}

		if status == constant.InviteStatusAccepted {
			_, err = a.gov.CreateGov(ctx, entities.CreateGovParams{
				UserID: invite.ID,
				CityID: invite.CityID,
				Role:   invite.Role,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
	if txErr != nil {
		return models.Invite{}, txErr
	}

	return invite, err
}
