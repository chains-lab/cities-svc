package app

import (
	"context"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/google/uuid"
)

func (a App) CreateInviteMayor(ctx context.Context, cityID uuid.UUID) (models.Invite, error) {
	newInvite, err := a.gov.CreateMayorInvite(ctx, cityID)
	if err != nil {
		return models.Invite{}, err
	}

	return newInvite, nil
}
