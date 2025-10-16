package admin

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/google/uuid"
)

type FilterParams struct {
	CityID *uuid.UUID
	Roles  []string
}

func (s Service) Filter(
	ctx context.Context,
	filters FilterParams,
	page, size uint64,
) (models.CityAdminsWithUserDataCollection, error) {
	res, err := s.db.FilterCityAdmins(ctx, filters, page, size)
	if err != nil {
		return models.CityAdminsWithUserDataCollection{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to filter city admin, cause: %w", err),
		)
	}

	userIDs := make([]uuid.UUID, 0, len(res.Data))
	for _, emp := range res.Data {
		userIDs = append(userIDs, emp.UserID)
	}

	profiles, err := s.userGuesser.Guess(ctx, userIDs...)
	if err != nil {
		return models.CityAdminsWithUserDataCollection{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to guess city admins profiles, cause: %w", err),
		)
	}

	return res.AddProfileData(profiles), err
}
