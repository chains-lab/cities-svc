package gov

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

type DeleteGovsFilters struct {
	UserID    *uuid.UUID
	CityID    *uuid.UUID
	CountryID *uuid.UUID
	Role      *string
}

func (g Gov) DeleteMany(ctx context.Context, filters DeleteGovsFilters) error {
	query := g.gov.New()

	if filters.UserID != nil {
		query = query.FilterUserID(*filters.UserID)
	}
	if filters.CityID != nil {
		query = query.FilterCityID(*filters.CityID)
	}
	if filters.CountryID != nil {
		query = query.FilterCountryID(*filters.CountryID)
	}
	if filters.Role != nil {
		err := enum.CheckCityGovRole(*filters.Role)
		if err != nil {
			return errx.ErrorInvalidGovRole.Raise(
				fmt.Errorf("invalid city gov role, cause: %w", err),
			)
		}
		query = query.FilterRole(*filters.Role)
	}

	err := query.Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete city govs, cause: %w", err),
		)
	}

	return nil
}

func (g Gov) Delete(ctx context.Context, InitiatorID, userID, cityID uuid.UUID) error {
	initiator, err := g.GetInitiatorGov(ctx, InitiatorID)
	if err != nil {
		return err
	}

	gov, err := g.GetGov(ctx, GetGovFilters{
		UserID: &userID,
	})
	if err != nil {
		return err
	}

	if initiator.CityID != cityID {
		return errx.ErrorInitiatorAndUserHaveDifferentCity.Raise(
			fmt.Errorf("initiator is not gov of city %s", cityID),
		)
	}

	if initiator.CityID != gov.CityID {
		return errx.ErrorInitiatorIsNotThisCityGov.Raise(
			fmt.Errorf("initiator is not gov of city %s", gov.CityID),
		)
	}

	access, err := enum.CompareCityGovRoles(gov.Role, initiator.Role)
	if err != nil {
		return errx.ErrorInvalidGovRole.Raise(
			fmt.Errorf("compare city gov roles: %w", err),
		)
	}
	if access >= 0 {
		return errx.ErrorInitiatorGovRoleHaveNotEnoughRights.Raise(
			fmt.Errorf("initiator have not enough rights to delete role %s", gov.Role),
		)
	}

	return g.delete(ctx, userID)
}

func (g Gov) RefuseOwnGov(ctx context.Context, userID uuid.UUID) error {
	gov, err := g.GetInitiatorGov(ctx, userID)
	if err != nil {
		return err
	}

	if gov.Role == enum.CityGovRoleMayor {
		return errx.ErrorCannotRefuseMayor.Raise(
			fmt.Errorf("mayor %s cannot refuse self gov", userID),
		)
	}

	return g.delete(ctx, userID)
}

func (g Gov) delete(ctx context.Context, userID uuid.UUID) error {
	err := g.gov.New().FilterUserID(userID).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete city gov, cause: %w", err),
		)
	}

	return nil
}

func (g Gov) deleteCityMayor(ctx context.Context, cityID uuid.UUID) error {
	err := g.gov.New().FilterCityID(cityID).FilterRole(enum.CityGovRoleMayor).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete city mayor, cause: %w", err),
		)
	}

	return nil
}
