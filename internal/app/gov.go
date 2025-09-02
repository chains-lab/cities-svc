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

func (a App) CheckNoActiveForUser(ctx context.Context, userID uuid.UUID) (models.Gov, error) {
	st := constant.GovStatusActive
	gov, err := a.gov.Get(ctx, entities.GetGovFilters{
		UserID: &userID,
		Status: &st,
	})
	if err != nil && !errors.Is(err, errx.ErrorCityGovNotFound) {
		return models.Gov{}, err
	}
	if err == nil {
		return models.Gov{}, errx.ErrorGovAlreadyExists.Raise(
			fmt.Errorf("active gov already exists for user %s", userID),
		)
	}

	return gov, nil
}

func (a App) CreateGovMayor(ctx context.Context, cityID, userID uuid.UUID) (models.Gov, error) {
	_, err := a.CheckNoActiveForUser(ctx, userID)
	if err != nil {
		return models.Gov{}, err
	}

	_, err = a.gov.Get(ctx, entities.GetGovFilters{
		CityID: &cityID,
		Role:   func(s string) *string { return &s }(constant.CityGovRoleMayor),
		Status: func(s string) *string { return &s }(constant.GovStatusActive),
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

type CreateGovParams struct {
	CityID uuid.UUID
	UserID uuid.UUID
	Label  string
	Role   string
}

func (a App) CreateGov(ctx context.Context, initiatorID uuid.UUID, param CreateGovParams) (models.Gov, error) {
	_, err := a.CheckNoActiveForUser(ctx, param.UserID)
	if err != nil {
		return models.Gov{}, err
	}

	err = constant.CheckCityGovRole(param.Role)
	if err != nil {
		return models.Gov{}, errx.ErrorInvalidCityGovRole.Raise(
			fmt.Errorf("failed to parse city gov role: %w", err),
		)
	}

	var newGov models.Gov
	var pagination pagi.Response

	switch param.Role {
	case constant.CityGovRoleMayor:
		_, err = a.gov.Get(ctx, entities.GetGovFilters{
			CityID: &param.CityID,
			Role:   func(s string) *string { return &s }(constant.CityGovRoleMayor),
			Status: func(s string) *string { return &s }(constant.GovStatusActive),
		})
		if err != nil && !errors.Is(err, errx.ErrorCityGovNotFound) {
			return models.Gov{}, err
		}
		if err == nil {
			return models.Gov{}, errx.ErrorGovAlreadyExists.Raise(
				fmt.Errorf("active mayor already exists in city %s", param.CityID),
			)
		}

		newGov, err = a.gov.CreateGov(ctx, entities.CreateGovParams{
			CityID: param.CityID,
			UserID: param.UserID,
			Role:   constant.CityGovRoleMayor,
		})
	case constant.CityGovRoleAdvisor:
		_, pagination, err = a.gov.SelectGovs(ctx, entities.SelectGovsFilters{
			CityID: &param.CityID,
			Role:   []string{constant.CityGovRoleAdvisor},
			Status: []string{constant.GovStatusActive},
		}, pagi.Request{
			Page: 1,
			Size: 1,
		}, nil)
		if err != nil {
			return models.Gov{}, err
		}
		if pagination.Total >= 10 {
			return models.Gov{}, errx.ErrorAdvisorMaxNumberReached.Raise(
				fmt.Errorf("city %s already has maximum number of advisors", param.CityID),
			)
		}

		newGov, err = a.gov.CreateGov(ctx, entities.CreateGovParams{
			CityID: param.CityID,
			UserID: param.UserID,
			Role:   constant.CityGovRoleAdvisor,
			Label:  param.Label,
		})
		if err != nil {
			return models.Gov{}, fmt.Errorf("create advisor gov: %w", err)
		}
	case constant.CityGovRoleMember, constant.CityGovRoleModerator:
		newGov, err = a.gov.CreateGov(ctx, entities.CreateGovParams{
			CityID: param.CityID,
			UserID: param.UserID,
			Role:   param.Role,
			Label:  param.Label,
		})
	}

	return newGov, nil
}

func (a App) GetGov(ctx context.Context, govID uuid.UUID) (models.Gov, error) {
	gov, err := a.gov.Get(ctx, entities.GetGovFilters{
		ID: &govID,
	})
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
	gov, err := a.gov.Get(ctx, entities.GetGovFilters{
		CityID: &cityID,
		UserID: &userID,
		Status: func(s string) *string { return &s }(constant.GovStatusActive),
	})
	if err != nil {
		return models.Gov{}, err
	}
	return gov, nil
}

func (a App) SelectGovs(
	ctx context.Context,
	params entities.SelectGovsFilters,
	pagiReq pagi.Request,
	sort []pagi.SortField) ([]models.Gov, pagi.Response, error) {
	return a.gov.SelectGovs(ctx, params, pagiReq, sort)
}

type UpdateGovParams struct {
	Label  *string
	Active *bool
}

func (a App) UpdateGov(ctx context.Context, initiatorUserID, govID uuid.UUID, params UpdateGovParams) (models.Gov, error) {
	gov, err := a.gov.Get(ctx, entities.GetGovFilters{
		ID: &govID,
	})
	if err != nil {
		return models.Gov{}, err
	}
	if gov.Active == false {
		return models.Gov{}, errx.ErrorCannotUpdateInactiveGov.Raise(
			fmt.Errorf("cannot update inactive gov %s", govID),
		)
	}
	if gov.UserID == initiatorUserID {
		return models.Gov{}, errx.ErrorCannotUpdateSelfGov.Raise(
			fmt.Errorf("cannot update self gov %s", govID),
		)
	}

	initiator, err := a.gov.Get(ctx, entities.GetGovFilters{
		CityID: &gov.CityID,
		UserID: &initiatorUserID,
		Active: func(b bool) *bool { return &b }(true),
	})
	if err != nil {
		switch {
		case errors.Is(err, errx.ErrorCityGovNotFound):
			return models.Gov{}, errx.ErrorNotActiveCityGovInitiator.Raise(
				fmt.Errorf("no active gov found for user %s in city %s", initiatorUserID, gov.CityID),
			)
		default:
			return models.Gov{}, fmt.Errorf("get initiator gov: %w", err)
		}
	}

	if constant.CityGovRoleMayor == gov.Role {
		return models.Gov{}, errx.ErrorCannotUpdateMayorGovByOther.Raise(
			fmt.Errorf("cannot update mayor gov %s by other", govID),
		)
	}
	access, err := constant.CompareCityGovRoles(initiator.Role, gov.Role)
	if err != nil {
		return models.Gov{}, errx.ErrorInvalidCityGovRole.Raise(err)
	}
	if access != 1 {
		return models.Gov{}, errx.ErrorEnoughRights.Raise(
			fmt.Errorf("initiator gov %s with role %s does not have enough rights to update gov %s with role %s",
				initiator.ID, initiator.Role, gov.ID, gov.Role),
		)
	}

	if params.Role != nil {
		access, err = constant.CompareCityGovRoles(initiator.Role, *params.Role)
		if err != nil {
			return models.Gov{}, errx.ErrorInvalidCityGovRole.Raise(err)
		}
		if access < 0 {
			return models.Gov{}, errx.ErrorEnoughRights.Raise(
				fmt.Errorf("initiator gov %s with role %s does not have enough rights to update gov %s to role %s",
					initiator.ID, initiator.Role, gov.ID, *params.Role),
			)
		}

		gov, err = a.UpdateGovRole(ctx, govID, *params.Role)
		if err != nil {
			return models.Gov{}, err
		}
	}

	if params.Label != nil {
		gov, err = a.UpdateGovLabel(ctx, govID, params.Label)
		if err != nil {
			return models.Gov{}, err
		}
	}

	if params.Active != nil {
		gov, err = a.UpdateGovActiveStatus(ctx, govID, *params.Active)
		if err != nil {
			return models.Gov{}, err
		}
	}

	return gov, nil
}

func (a App) UpdateGovLabel(ctx context.Context, govID uuid.UUID, newLabel *string) (models.Gov, error) {
	err := a.gov.UpdateOne(ctx, govID, entities.UpdateGovParams{
		Label: newLabel,
	})
	if err != nil {
		return models.Gov{}, fmt.Errorf("update gov label: %w", err)
	}

	updatedGov, err := a.gov.Get(ctx, entities.GetGovFilters{
		ID: &govID,
	})
	if err != nil {
		return models.Gov{}, fmt.Errorf("get updated gov: %w", err)
	}

	return updatedGov, nil
}

func (a App) UpdateGovRole(ctx context.Context, govID uuid.UUID, newRole string) (models.Gov, error) {
	err := constant.CheckCityGovRole(newRole)
	if err != nil {
		return models.Gov{}, errx.ErrorInvalidCityGovRole.Raise(
			fmt.Errorf("failed to parse city gov role: %w", err),
		)
	}

	gov, err := a.gov.Get(ctx, entities.GetGovFilters{
		ID: &govID,
	})
	if err != nil {
		return models.Gov{}, err
	}

	if gov.Role == newRole {
		return gov, nil
	}

	if newRole == constant.CityGovRoleMayor {
		return models.Gov{}, errx.ErrorEnoughRights.Raise(
			fmt.Errorf("only transfer to new mayor is allowed"),
		)
	}
	if gov.Role == constant.CityGovRoleMayor {
		return models.Gov{}, errx.ErrorEnoughRights.Raise(
			fmt.Errorf("cannot update mayor gov %s by other", govID),
		)
	}

	err = a.gov.UpdateOne(ctx, govID, entities.UpdateGovParams{
		Role: &newRole,
	})
	if err != nil {
		return models.Gov{}, fmt.Errorf("update gov role: %w", err)
	}

	updatedGov, err := a.gov.Get(ctx, entities.GetGovFilters{
		ID: &govID,
	})
	if err != nil {
		return models.Gov{}, fmt.Errorf("get updated gov: %w", err)
	}

	return updatedGov, nil
}

func (a App) UpdateGovActiveStatus(ctx context.Context, govID uuid.UUID, active bool) (models.Gov, error) {
	gov, err := a.gov.Get(ctx, entities.GetGovFilters{
		ID: &govID,
	})
	if err != nil {
		return models.Gov{}, err
	}

	if gov.Active == active {
		return gov, nil
	}

	if active {
		_, err := a.CheckNoActiveForUser(ctx, gov.UserID)
		if err != nil {
			return models.Gov{}, err
		}

		if gov.Role == constant.CityGovRoleMayor {
			_, err := a.gov.Get(ctx, entities.GetGovFilters{
				CityID: &gov.CityID,
				Role:   func(s string) *string { return &s }(constant.CityGovRoleMayor),
				Active: func(b bool) *bool { return &b }(true),
			})
			if err != nil && !errors.Is(err, errx.ErrorCityGovNotFound) {
				return models.Gov{}, err
			}
			if err == nil {
				return models.Gov{}, errx.ErrorGovAlreadyExists.Raise(
					fmt.Errorf("active mayor already exists in city %s", gov.CityID),
				)
			}
		}
	}

	err = a.gov.UpdateOne(ctx, govID, entities.UpdateGovParams{
		Active: &active,
	})
	if err != nil {
		return models.Gov{}, fmt.Errorf("update gov active status: %w", err)
	}

	updatedGov, err := a.gov.Get(ctx, entities.GetGovFilters{
		ID: &govID,
	})
	if err != nil {
		return models.Gov{}, fmt.Errorf("get updated gov: %w", err)
	}

	return updatedGov, nil
}

func (a App) TransferGovMayor(ctx context.Context, cityID, newMayorUserID uuid.UUID) error {
	_, err := a.CheckNoActiveForUser(ctx, newMayorUserID)
	if err != nil {
		return err
	}

	currentMayor, err := a.gov.Get(ctx, entities.GetGovFilters{
		CityID: &cityID,
		Role:   func(s string) *string { return &s }(constant.CityGovRoleMayor),
		Active: func(b bool) *bool { return &b }(true),
	})
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

func (a App) DeactivateGov(ctx context.Context, govID uuid.UUID) error {
	err := a.gov.UpdateOne(ctx, govID, entities.UpdateGovParams{
		Status: func(s string) *string { return &s }(constant.GovStatusInactive),
	})
	if err != nil {
		return fmt.Errorf("deactivate gov: %w", err)
	}
	return nil
}
