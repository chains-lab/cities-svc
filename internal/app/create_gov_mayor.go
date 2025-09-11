package app

import (
	"context"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/google/uuid"
)

func (a App) CreateGovMayor(ctx context.Context, cityID, userID uuid.UUID) (models.Gov, error) {
	newGov, err := a.gov.CreateMayor(ctx, userID, cityID)
	if err != nil {
		return models.Gov{}, err
	}

	return newGov, nil
}
