package status

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
)

func (s Service) Get(ctx context.Context, ID string) (models.CityStatus, error) {
	status, err := s.db.GetStatusByID(ctx, ID)
	if err != nil {
		return models.CityStatus{}, errx.ErrorInternal.Raise(
			fmt.Errorf("error getting status by ID: %s", ID),
		)
	}

	if status.IsNil() {
		return models.CityStatus{}, errx.ErrorCityStatusNotFound.Raise(
			fmt.Errorf("status not found with ID: %s", ID),
		)
	}

	return status, nil
}

func (s Service) Exists(ctx context.Context, ID string) (bool, error) {
	status, err := s.db.GetStatusByID(ctx, ID)
	if err != nil {
		return false, errx.ErrorInternal.Raise(
			fmt.Errorf("error checking existence of status by ID: %s", ID),
		)
	}

	return !status.IsNil(), nil
}
