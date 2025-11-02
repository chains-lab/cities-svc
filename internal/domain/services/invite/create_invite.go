package invite

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/domain/enum"
	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/google/uuid"
)

type CreateInviteParams struct {
	CityID uuid.UUID
	UserID uuid.UUID
	Role   string
}

func (s Service) CreateInvite(ctx context.Context, params CreateInviteParams) (models.Invite, error) {
	err := enum.CheckCityAdminRole(params.Role)
	if err != nil {
		return models.Invite{}, err
	}

	exist, err := s.db.ExistsAdmin(ctx, params.UserID)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to check if city admin exists: %w", err),
		)
	}
	if exist {
		return models.Invite{}, errx.ErrorUserAlreadyCityAdmin.Raise(
			fmt.Errorf("user is already a city admin"),
		)
	}

	city, err := s.db.GetCity(ctx, params.CityID)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get city: %w", err),
		)
	}
	if city.IsNil() {
		return models.Invite{}, errx.ErrorCityNotFound.Raise(
			fmt.Errorf("city not found"),
		)
	}

	cityStatus, err := s.db.GetCityStatus(ctx, params.CityID)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get city status: %w", err),
		)
	}

	if !cityStatus.IsNil() && cityStatus.AllowedAdmin {
		return models.Invite{}, errx.ErrorCityStatusNotAllowedAdmin.Raise(
			fmt.Errorf("city status does not allow adminernment roles"),
		)
	}

	invite, err := s.db.CreateInvite(ctx, models.Invite{
		CityID: params.CityID,
		UserID: params.UserID,
		Role:   params.Role,
	})
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to create invite: %w", err),
		)
	}

	return invite, nil
}
