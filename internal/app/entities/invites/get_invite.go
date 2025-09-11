package invites

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/google/uuid"
)

func (i Invite) Get(ctx context.Context, id uuid.UUID) (models.Invite, error) {
	inv, err := i.query.New().FilterID(id).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Invite{}, errx.ErrorInviteNotFound.Raise(
				fmt.Errorf("invite not found: %w", err),
			)
		default:
			return models.Invite{}, errx.ErrorInternal.Raise(fmt.Errorf("get invite by id, cause %w", err))
		}
	}

	return modelsFromDB(inv), nil
}
