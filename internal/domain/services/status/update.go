package status

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
)

type UpdateParams struct {
	Name         *string
	Description  *string
	AllowedAdmin *bool
}

func (s Service) Update(ctx context.Context, ID string, params UpdateParams) (models.CityStatus, error) {
	status, err := s.Get(ctx, ID)
	if err != nil {
		return models.CityStatus{}, err
	}

	if params.Name != nil {
		status.Name = *params.Name
	}
	if params.Description != nil {
		status.Description = *params.Description
	}
	if params.AllowedAdmin != nil {
		status.AllowedAdmin = *params.AllowedAdmin
	}

	updatedStatus, err := s.db.UpdateStatus(ctx, ID, params)
	if err != nil {
		return models.CityStatus{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update status with ID %s: %w", ID, err),
		)
	}

	return updatedStatus, nil
}
