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

type CreateParams struct {
	UserID uuid.UUID
	CityID uuid.UUID
	Role   string
}

func (g Gov) Create(ctx context.Context, params CreateParams) (models.Gov, error) {
	err := constant.CheckCityGovRole(params.Role)
	if err != nil {
		return models.Gov{}, errx.ErrorInvalidGovRole.Raise(
			fmt.Errorf("invalid city gov role, cause: %w", err),
		)
	}

	now := time.Now().UTC()

	stmt := dbx.Gov{
		UserID:    params.UserID,
		CityID:    params.CityID,
		Role:      params.Role,
		UpdatedAt: now,
		CreatedAt: now,
	}

	resp := models.Gov{
		UserID:    params.UserID,
		CityID:    params.CityID,
		Role:      params.Role,
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
