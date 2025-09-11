package gov

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/constant"
	"github.com/chains-lab/cities-svc/internal/dbx"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/google/uuid"
)

type UpdateGovParams struct {
	Role      *string
	Label     *string
	UpdatedAt time.Time
}

func (g Gov) UpdateOne(ctx context.Context, userID uuid.UUID, params UpdateGovParams) (models.Gov, error) {
	if (params.Role == nil) && (params.Label == nil) {
		return models.Gov{}, nil
	}

	stmt := dbx.UpdateCityGovParams{}

	if params.Role != nil {
		err := constant.CheckCityGovRole(*params.Role)
		if err != nil {
			return models.Gov{}, errx.ErrorInvalidGovRole.Raise(
				fmt.Errorf("invalid city gov role, cause: %w", err),
			)
		}
		stmt.Role = params.Role
	}

	if params.Label != nil {
		if *params.Label != "" {
			stmt.Label.String = *params.Label
		} else {
			stmt.Label.Valid = false
		}
	}

	stmt.UpdatedAt = &params.UpdatedAt

	err := g.govQ.New().FilterUserID(userID).Update(ctx, stmt)
	if err != nil {
		return models.Gov{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update city gov, cause: %w", err),
		)
	}

	return g.Get(ctx, GetGovFilters{UserID: &userID})
}

func (g Gov) DeleteOne(ctx context.Context, userID uuid.UUID) error {
	err := g.govQ.New().FilterUserID(userID).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete city gov, cause: %w", err),
		)
	}

	return nil
}
