package gov

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/dbx"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

func (g Gov) CreateMayor(ctx context.Context, userID, cityID uuid.UUID) (models.Gov, error) {
	_, err := g.GetGov(ctx, GetGovFilters{
		CityID: &cityID,
		Role:   func(s string) *string { return &s }(enum.CityGovRoleMayor),
	})
	if err != nil && !errors.Is(err, errx.ErrorCityGovNotFound) {
		return models.Gov{}, err
	}
	if err == nil {
		return models.Gov{}, errx.ErrorGovAlreadyExists.Raise(
			fmt.Errorf("active mayor already exists in city %s", cityID),
		)
	}

	now := time.Now().UTC()

	stmt := dbx.Gov{
		UserID:    userID,
		CityID:    cityID,
		Role:      enum.CityGovRoleMayor,
		UpdatedAt: now,
		CreatedAt: now,
	}

	resp := models.Gov{
		UserID:    userID,
		CityID:    cityID,
		Role:      enum.CityGovRoleMayor,
		UpdatedAt: now,
		CreatedAt: now,
	}

	err = g.govQ.New().Insert(ctx, stmt)
	if err != nil {
		return models.Gov{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to creating city gov: %w", err),
		)
	}

	return resp, nil
}
