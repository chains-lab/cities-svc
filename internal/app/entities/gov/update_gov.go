package gov

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/dbx"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/google/uuid"
)

type UpdateGovParams struct {
	Label *string
}

func (g Gov) UpdateOne(ctx context.Context, userID uuid.UUID, params UpdateGovParams) (models.Gov, error) {
	if params.Label == nil {
		return models.Gov{}, nil
	}

	stmt := dbx.UpdateCityGovParams{
		UpdatedAt: time.Now().UTC(),
	}

	if params.Label != nil {
		if *params.Label != "" {
			stmt.Label.String = *params.Label
		} else {
			stmt.Label.Valid = false
		}
	}

	err := g.govQ.New().FilterUserID(userID).Update(ctx, stmt)
	if err != nil {
		return models.Gov{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update city gov, cause: %w", err),
		)
	}

	return g.GetGov(ctx, GetGovFilters{UserID: &userID})
}
