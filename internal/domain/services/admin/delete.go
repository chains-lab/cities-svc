package admin

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/google/uuid"
)

func (s Service) Delete(ctx context.Context, UserID, cityID uuid.UUID) error {
	_, err := s.Get(ctx, GetFilters{
		UserID: &UserID,
		CityID: &cityID,
	})
	if err != nil {
		return err
	}

	err = s.db.DeleteCityAdmin(ctx, UserID, cityID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete city admin, cause: %w", err),
		)
	}

	return nil
}
