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

	err = s.db.DeleteCityAdmin(ctx, userID, mod.CityID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("delete admin: %w", err),
		)
	}

	return nil
}
