package app

import (
	"context"

	"github.com/chains-lab/cities-svc/internal/app/entities/gov"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/google/uuid"
)

type UpdateOwnGovParams struct {
	Label *string
}

func (a App) UpdateOwnActiveGov(ctx context.Context, userID uuid.UUID, params UpdateOwnGovParams) (models.Gov, error) {
	g, err := a.GetInitiatorGov(ctx, userID)
	if err != nil {
		return models.Gov{}, err
	}

	entitiesParams := gov.UpdateGovParams{}
	if params.Label != nil {
		entitiesParams.Label = params.Label
	}

	return a.gov.UpdateOne(ctx, g.UserID, entitiesParams)
}
