package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/app/entities"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/constant"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
)

func (a App) GetInitiator(ctx context.Context, cityID, userID uuid.UUID) (models.Gov, error) {
	initiator, err := a.gov.Get(ctx, entities.GetGovFilters{
		CityID: &cityID,
		UserID: &userID,
		Active: func(b bool) *bool { return &b }(true),
	})
	if err != nil {
		switch {
		case errors.Is(err, errx.ErrorCityGovNotFound):
			return models.Gov{}, errx.ErrorNotActiveCityGovInitiator.Raise(
				fmt.Errorf("no active gov found for user %s in city %s", userID, cityID),
			)
		default:
			return models.Gov{}, fmt.Errorf("get initiator gov: %w", err)
		}
	}
	return initiator, nil
}

func (a App) CreateGovMayor(ctx context.Context, cityID, userID uuid.UUID) (models.Gov, error) {
	_, err := a.gov.GetActiveForUser(ctx, userID)
	if err != nil && !errors.Is(err, errx.ErrorCityGovNotFound) {
		return models.Gov{}, err
	}
	if err == nil {
		return models.Gov{}, errx.ErrorGovAlreadyExists.Raise(
			fmt.Errorf("active gov already exists for user %s", userID),
		)
	}

	_, err = a.gov.GetActiveMayorForCity(ctx, cityID)
	if err != nil && !errors.Is(err, errx.ErrorCityGovNotFound) {
		return models.Gov{}, err
	}
	if err == nil {
		return models.Gov{}, errx.ErrorGovAlreadyExists.Raise(
			fmt.Errorf("active mayor already exists in city %s", cityID),
		)
	}

	newGov, err := a.gov.CreateGov(ctx, entities.CreateGovParams{
		CityID: cityID,
		UserID: userID,
		Role:   constant.CityGovRoleMayor,
	})
	if err != nil {
		return models.Gov{}, fmt.Errorf("create mayor gov: %w", err)
	}

	return newGov, nil
}

func (a App) CreateGovAdvisor(ctx context.Context, cityID, userID uuid.UUID, label *string) (models.Gov, error) {
	_, err := a.gov.GetActiveForUser(ctx, userID)
	if err != nil && errors.Is(err, errx.ErrorCityGovNotFound) {
		return models.Gov{}, err
	}
	if err == nil {
		return models.Gov{}, errx.ErrorGovAlreadyExists.Raise(
			fmt.Errorf("active gov already exists for user %s", userID),
		)
	}

	_, pagination, err := a.gov.SelectGovs(ctx, entities.SelectGovsFilters{
		CityID: &cityID,
		Role:   []string{constant.CityGovRoleAdvisor},
		Active: func(b bool) *bool { return &b }(true),
	}, pagi.Request{
		Page: 1,
		Size: 1,
	}, nil)
	if err != nil {
		return models.Gov{}, err
	}
	if pagination.Total >= 10 {
		return models.Gov{}, errx.ErrorAdvisorMaxNumberReached.Raise(
			fmt.Errorf("city %s already has maximum number of advisors", cityID),
		)
	}

	newGov, err := a.gov.CreateGov(ctx, entities.CreateGovParams{
		CityID: cityID,
		UserID: userID,
		Role:   constant.CityGovRoleAdvisor,
		Label:  label,
	})
	if err != nil {
		return models.Gov{}, fmt.Errorf("create advisor gov: %w", err)
	}

	return newGov, nil
}

func (a App) CreateGovMember(ctx context.Context, cityID, userID uuid.UUID, label *string) (models.Gov, error) {
	_, err := a.gov.GetActiveForUser(ctx, userID)
	if err != nil && errors.Is(err, errx.ErrorCityGovNotFound) {
		return models.Gov{}, err
	}
	if err == nil {
		return models.Gov{}, errx.ErrorGovAlreadyExists.Raise(
			fmt.Errorf("active gov already exists for user %s", userID),
		)
	}

	newGov, err := a.gov.CreateGov(ctx, entities.CreateGovParams{
		CityID: cityID,
		UserID: userID,
		Role:   constant.CityGovRoleMember,
		Label:  label,
	})
	if err != nil {
		return models.Gov{}, fmt.Errorf("create member gov: %w", err)
	}

	return newGov, nil
}

func (a App) CreateGovModerator(ctx context.Context, cityID, userID uuid.UUID, label *string) (models.Gov, error) {
	_, err := a.gov.GetActiveForUser(ctx, userID)
	if err != nil && errors.Is(err, errx.ErrorCityGovNotFound) {
		return models.Gov{}, err
	}
	if err == nil {
		return models.Gov{}, errx.ErrorGovAlreadyExists.Raise(
			fmt.Errorf("active gov already exists for user %s", userID),
		)
	}

	newGov, err := a.gov.CreateGov(ctx, entities.CreateGovParams{
		CityID: cityID,
		UserID: userID,
		Role:   constant.CityGovRoleModerator,
		Label:  label,
	})
	if err != nil {
		return models.Gov{}, fmt.Errorf("create moderator gov: %w", err)
	}

	return newGov, nil
}

func (a App) DeactivateGov(ctx context.Context, govID uuid.UUID) error {
	err := a.gov.UpdateOne(ctx, govID, entities.UpdateGovParams{
		Active: func(b bool) *bool { return &b }(false),
	})
	if err != nil {
		return fmt.Errorf("deactivate gov: %w", err)
	}
	return nil
}

func (a App) GetGov(ctx context.Context, govID uuid.UUID) (models.Gov, error) {
	gov, err := a.gov.Get(ctx, govID)
	if err != nil {
		switch {
		case errors.Is(err, errx.ErrorCityGovNotFound):
			return models.Gov{}, err
		default:
			return models.Gov{}, fmt.Errorf("get gov: %w", err)
		}
	}
	return gov, nil
}

func (a App) GetActiveGovForUserInCity(ctx context.Context, cityID, userID uuid.UUID) (models.Gov, error) {
	gov, err := a.gov.GetActiveForUserInCity(ctx, cityID, userID)
	if err != nil {
		return models.Gov{}, err
	}
	return gov, nil
}

func (a App) TransferGovMayor(ctx context.Context, cityID, newMayorUserID uuid.UUID) error {
	_, err := a.gov.GetActiveForUser(ctx, newMayorUserID)
	if err != nil && !errors.Is(err, errx.ErrorCityGovNotFound) {
		return err
	}
	if err == nil {
		return errx.ErrorGovAlreadyExists.Raise(
			fmt.Errorf("active gov already exists for user %s", newMayorUserID),
		)
	}

	currentMayor, err := a.gov.GetActiveMayorForCity(ctx, cityID)
	if err != nil {
		return err
	}

	txErr := a.transaction(func(txCtx context.Context) error {
		err = a.gov.UpdateOne(txCtx, currentMayor.ID, entities.UpdateGovParams{
			Active: func(b bool) *bool { return &b }(false),
		})
		if err != nil {
			return err
		}

		_, err := a.gov.CreateGov(txCtx, entities.CreateGovParams{
			CityID: cityID,
			UserID: newMayorUserID,
			Role:   constant.CityGovRoleMayor,
		})
		if err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		return txErr
	}

	return nil
}

func (a App) SelectGovs(
	ctx context.Context,
	params entities.SelectGovsFilters,
	pagiReq pagi.Request,
	sort []pagi.SortField) ([]models.Gov, pagi.Response, error) {
	return a.gov.SelectGovs(ctx, params, pagiReq, sort)
}
