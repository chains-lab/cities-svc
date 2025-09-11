package app

import (
	"context"

	"github.com/google/uuid"
)

func (a App) DeleteGov(ctx context.Context, initiatorID, userID uuid.UUID) error {
	return a.gov.Delete(ctx, initiatorID, userID)
}
