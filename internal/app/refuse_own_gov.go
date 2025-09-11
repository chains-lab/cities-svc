package app

import (
	"context"

	"github.com/google/uuid"
)

func (a App) RefuseOwnGov(ctx context.Context, userID uuid.UUID) error {
	err := a.gov.RefuseOwnGov(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}
