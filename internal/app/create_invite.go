package app

import (
	"context"

	"github.com/chains-lab/cities-svc/internal/app/entities/gov"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/google/uuid"
)

type CreateInviteParams struct {
	InitiatorID uuid.UUID
	Role        string
}

func (a App) CreateInvite(ctx context.Context, params CreateInviteParams) (models.Invite, error) {
	p := gov.CreateInviteParams{
		InitiatorID: params.InitiatorID,
		Role:        params.Role,
	}

	newInvite, err := a.gov.CreateInvite(ctx, p)
	if err != nil {
		return models.Invite{}, err
	}

	return newInvite, nil
}
