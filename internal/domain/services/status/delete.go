package status

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
)

func (s Service) DeleteStatus(ctx context.Context, ID string) error {
	exist, err := s.db.ExistsCitiesWitStatus(ctx, ID)
	if err != nil {
		return err
	}
	if exist {
		return errx.ErrorStatusInUse.Raise(
			fmt.Errorf("cannot delete status with ID %s: it is in use by some cities", ID),
		)
	}

	err = s.db.DeleteStatus(ctx, ID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete status with ID %s: %w", ID, err),
		)
	}

	return nil
}
