package admin

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/google/uuid"
)

func (s Service) RefuseOwn(ctx context.Context, userID uuid.UUID) error {
	mod, err := s.GetInitiator(ctx, userID)
	if err != nil {
		return err
	}

	err = s.db.DeleteAdmin(ctx, userID, mod.Admin.CityID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete city admin, cause: %w", err),
		)
	}

	err = s.event.CityAdminDeleted(ctx, userID, mod.Admin.CityID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to emit city admin deleted event, cause: %w", err),
		)
	}

	return nil
}
