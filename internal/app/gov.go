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

func (a App) CreateGovMayor(ctx context.Context, cityID, userID uuid.UUID) (models.Gov, error) {
	_, err := a.gov.Get(ctx, entities.GetGovFilters{
		CityID: &cityID,
		Role:   func(s string) *string { return &s }(constant.CityGovRoleMayor),
	})
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

func (a App) GetInitiatorGov(ctx context.Context, initiatorID uuid.UUID) (models.Gov, error) {
	initiator, err := a.Get(ctx, initiatorID)
	if err != nil {
		switch {
		case errors.Is(err, errx.ErrorCityGovNotFound):
			return models.Gov{}, errx.ErrorNotActiveCityGovInitiator.Raise(
				fmt.Errorf("initiator %s is not an active city gov", initiatorID),
			)
		default:
			return models.Gov{}, err
		}
	}

	return initiator, nil
}

func (a App) Get(ctx context.Context, userID uuid.UUID) (models.Gov, error) {
	return a.gov.Get(ctx, entities.GetGovFilters{
		UserID: &userID,
	})

}

func (a App) GetForCity(ctx context.Context, cityID, userID uuid.UUID) (models.Gov, error) {
	return a.gov.Get(ctx, entities.GetGovFilters{
		CityID: &cityID,
		UserID: &userID,
	})
}

func (a App) GetForCityAndRole(ctx context.Context, userID, cityID uuid.UUID, role string) (models.Gov, error) {
	return a.gov.Get(ctx, entities.GetGovFilters{
		UserID: &userID,
		CityID: &cityID,
		Role:   &role,
	})
}

type SearchGovsFilters struct {
	CityID *uuid.UUID
	Roles  []string
}

func (a App) SearchGovs(
	ctx context.Context,
	filters SearchGovsFilters,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.Gov, pagi.Response, error) {
	input := entities.SelectGovsFilters{}
	if filters.CityID != nil {
		input.CityID = filters.CityID
	}
	if len(filters.Roles) > 0 && filters.Roles != nil {
		input.Role = filters.Roles
	}

	return a.gov.SelectGovs(ctx, input, pag, sort)
}

type UpdateOwnGovParams struct {
	Label *string
}

func (a App) UpdateOwnActiveGov(ctx context.Context, userID uuid.UUID, params UpdateOwnGovParams) (models.Gov, error) {
	gov, err := a.GetInitiatorGov(ctx, userID)
	if err != nil {
		return models.Gov{}, err
	}

	entitiesParams := entities.UpdateGovParams{}
	if params.Label != nil {
		entitiesParams.Label = params.Label
	}
	res, err := a.gov.UpdateOne(ctx, gov.UserID, entitiesParams)
	if err != nil {
		return models.Gov{}, err
	}

	return res, nil
}

func (a App) RefuseOwnGov(ctx context.Context, userID uuid.UUID) error {
	gov, err := a.GetInitiatorGov(ctx, userID)
	if err != nil {
		return err
	}

	if gov.Role == constant.CityGovRoleMayor {
		return errx.ErrorCannotRefuseMayor.Raise(
			fmt.Errorf("mayor %s cannot refuse self gov", userID),
		)
	}

	err = a.gov.DeleteOne(ctx, gov.UserID)
	if err != nil {
		return fmt.Errorf("refuse own gov: %w", err)
	}

	return nil
}

func (a App) DeleteGov(ctx context.Context, initiatorID, userID uuid.UUID) (models.Gov, error) {
	initiator, err := a.GetInitiatorGov(ctx, initiatorID)
	if err != nil {
		return models.Gov{}, err
	}

	gov, err := a.gov.Get(ctx, entities.GetGovFilters{
		UserID: &userID,
	})
	if err != nil {
		return models.Gov{}, err
	}

	access, err := constant.CompareCityGovRoles(initiator.Role, gov.Role)
	if err != nil {
		return models.Gov{}, errx.ErrorInvalidGovRole.Raise(
			fmt.Errorf("invalid city gov role, cause: %w", err),
		)
	}
	if access != 1 {
		return models.Gov{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("no access to delete gov %s by initiator %s", userID, initiatorID),
		)
	}

	if gov.Role == constant.CityGovRoleMayor {
		return models.Gov{}, errx.ErrorCannotRefuseMayor.Raise(
			fmt.Errorf("mayor %s cannot refuse self gov", userID),
		)
	}

	err = a.gov.DeleteOne(ctx, gov.UserID)
	if err != nil {
		return models.Gov{}, fmt.Errorf("refuse own gov: %w", err)
	}

	return gov, nil
}

func (a App) TransferGovMayor(ctx context.Context, initiatorID, UserID uuid.UUID) error {
	initiator, err := a.GetInitiatorGov(ctx, initiatorID)
	if err != nil {
		return err
	}

	_, err = a.Get(ctx, UserID)
	if err != nil && !errors.Is(err, errx.ErrorCityGovNotFound) {
		return err
	}
	if err == nil {
		return errx.ErrorGovAlreadyExists.Raise(
			fmt.Errorf("active mayor already exists in city %s", initiator.CityID),
		)
	}

	txErr := a.transaction(func(txCtx context.Context) error {
		err = a.gov.DeleteOne(txCtx, initiator.CityID)
		if err != nil {
			return err
		}

		_, err = a.gov.CreateGov(txCtx, entities.CreateGovParams{
			CityID: initiator.CityID,
			UserID: UserID,
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
