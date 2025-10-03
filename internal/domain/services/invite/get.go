package invite

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) GetInvite(ctx context.Context, ID uuid.UUID) (models.Invite, error) {
	inv, err := s.db.GetInvite(ctx, ID)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(fmt.Errorf("failed to get invite by ID, cause %w", err))
	}

	if inv.IsNil() {
		return models.Invite{}, errx.ErrorInviteNotFound.Raise(fmt.Errorf("invite %s is not found", ID))
	}

	return inv, err
}
