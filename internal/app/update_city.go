package app

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/app/entities/city"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type UpdateCityParams struct {
	Name     *string
	Point    *orb.Point
	Icon     *string
	Slug     *string
	Timezone *string
}

func (a App) UpdateCity(ctx context.Context, cityID, initiatorID uuid.UUID, role string, params UpdateCityParams) (models.City, error) {
	if role == roles.User {
		gov, err := a.GetInitiatorGov(ctx, initiatorID)
		if err != nil {
			return models.City{}, err
		}
		if gov.CityID != cityID {
			return models.City{}, errx.ErrorInitiatorIsNotThisCityGov.Raise(
				fmt.Errorf("initiator %s is not the city %s", initiatorID.String(), cityID.String()),
			)
		}
		if gov.Role == enum.CityGovRoleMayor || gov.Role == enum.CityGovRoleModerator || gov.Role == enum.CityGovRoleAdvisor {
			return models.City{}, errx.ErrorInitiatorGovRoleHaveNotEnoughRights.Raise(
				fmt.Errorf("initiator %s have not enough rights", initiatorID.String()),
			)
		}
	}

	update := city.UpdateCityParams{}

	if params.Point != nil {
		update.Point = params.Point
	}
	if params.Slug != nil {
		update.Slug = params.Slug
	}
	if params.Name != nil {
		update.Name = params.Name
	}
	if params.Timezone != nil {
		update.Timezone = params.Timezone
	}
	if params.Icon != nil {
		update.Icon = params.Icon
	}

	return a.cities.UpdateOne(ctx, cityID, update)
}
