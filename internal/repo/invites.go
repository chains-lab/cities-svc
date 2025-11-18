package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/cities-svc/internal/repo/pgdb"
	"github.com/google/uuid"
)

func (r *Repo) CreateInvite(ctx context.Context, input models.Invite) error {
	schema := modelToInviteSchema(input)

	return r.sql.invites.New().Insert(ctx, schema)
}

func (r *Repo) GetInvite(ctx context.Context, ID uuid.UUID) (models.Invite, error) {
	row, err := r.sql.invites.New().FilterID(ID).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.Invite{}, err
	case err != nil:
		return models.Invite{}, err
	}

	return inviteSchemaToModel(row), nil
}

func (r *Repo) UpdateInviteStatus(ctx context.Context, inviteID uuid.UUID, status string) error {
	err := r.sql.invites.New().
		FilterID(inviteID).
		UpdateStatus(status).
		Update(ctx)
	return err
}

func inviteSchemaToModel(s pgdb.Invite) models.Invite {
	res := models.Invite{
		ID:        s.ID,
		Status:    s.Status,
		Role:      s.Role,
		CityID:    s.CityID,
		UserID:    s.UserID,
		CreatedAt: s.CreatedAt,
		ExpiresAt: s.ExpiresAt,
	}

	return res
}

func modelToInviteSchema(m models.Invite) pgdb.Invite {
	res := pgdb.Invite{
		ID:        m.ID,
		Status:    m.Status,
		Role:      m.Role,
		CityID:    m.CityID,
		UserID:    m.UserID,
		ExpiresAt: m.ExpiresAt,
		CreatedAt: m.CreatedAt,
	}

	return res
}
