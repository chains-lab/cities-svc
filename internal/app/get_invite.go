package app

import (
	"context"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/google/uuid"
)

func (a App) GetInvite(ctx context.Context, initiatorID uuid.UUID) (models.Invite, error) {
	return a.gov.GetInvite(ctx, initiatorID)
}
