package app

import (
	"context"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/google/uuid"
)

func (a App) AnswerToInvite(ctx context.Context, initiatorID uuid.UUID, status, token string) (models.Invite, error) {
	var invite models.Invite
	var err error

	txErr := a.transaction(func(ctx context.Context) error {
		invite, err = a.gov.AnsweredInvite(ctx, initiatorID, token, status)
		if err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		return models.Invite{}, txErr
	}

	return invite, err
}
