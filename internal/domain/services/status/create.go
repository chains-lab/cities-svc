package status

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
)

type CreateParams struct {
	ID           string
	Name         string
	Description  string
	AllowedAdmin bool
	AllowedMod   bool
}

func (s Service) CreateStatus(ctx context.Context, params CreateParams) (models.CityStatus, error) {
	exist, err := s.Exists(ctx, params.ID)
	if err != nil {
		return models.CityStatus{}, err
	}
	if exist {
		return models.CityStatus{}, errx.ErrorStatusAlreadyExists.Raise(
			fmt.Errorf("status with id %s already exists", params.ID),
		)
	}

	now := time.Now().UTC()
	status := models.CityStatus{
		ID:           params.ID,
		Name:         params.Name,
		Description:  params.Description,
		AllowedAdmin: params.AllowedAdmin,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err = s.db.CreateStatus(ctx, status)
	if err != nil {
		return models.CityStatus{}, err
	}

	return status, nil
}
