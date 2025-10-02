package citymod

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

func (s Service) RefuseOwn(ctx context.Context, userID uuid.UUID) error {
	mod, err := s.GetInitiator(ctx, userID)
	if err != nil {
		return err
	}

	if mod.Role == enum.CityGovRoleMayor {
		return errx.ErrorCannotRefuseMayor.Raise(
			fmt.Errorf("mayor %s cannot refuse self citymod", userID),
		)
	}

	err = s.db.DeleteCityModer(ctx, userID, mod.CityID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("delete citymod: %w", err),
		)
	}

	return nil
}
