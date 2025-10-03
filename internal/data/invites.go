package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/chains-lab/cities-svc/internal/data/pgdb"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (d *Database) CreateInvite(ctx context.Context, input models.Invite) error {
	schema := modelToInviteSchema(input)

	return d.sql.invites.New().Insert(ctx, schema)
}

func (d *Database) GetInvite(ctx context.Context, ID uuid.UUID) (models.Invite, error) {
	row, err := d.sql.invites.New().FilterID(ID).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.Invite{}, err
	case err != nil:
		return models.Invite{}, err
	}

	return inviteSchemaToModel(row), nil
}

func (d *Database) UpdateInviteStatus(ctx context.Context, inviteID uuid.UUID, userID uuid.UUID, status string, answeredAt time.Time) error {
	err := d.sql.invites.New().
		FilterID(inviteID).
		UpdateStatus(status).
		UpdateUserID(userID).
		UpdateAnsweredAt(answeredAt).
		Update(ctx)
	return err
}

func inviteSchemaToModel(s pgdb.Invite) models.Invite {
	res := models.Invite{
		ID:        s.ID,
		Status:    s.Status,
		Role:      s.Role,
		CityID:    s.CityID,
		Token:     s.Token,
		CreatedAt: s.CreatedAt,
		ExpiresAt: s.ExpiresAt,
	}
	if s.UserID.Valid {
		res.UserID = &s.UserID.UUID
	}
	if s.AnsweredAt.Valid {
		res.AnsweredAt = &s.AnsweredAt.Time
	}

	return res
}

func modelToInviteSchema(m models.Invite) pgdb.Invite {
	res := pgdb.Invite{
		ID:        m.ID,
		Status:    m.Status,
		Role:      m.Role,
		CityID:    m.CityID,
		Token:     m.Token,
		ExpiresAt: m.ExpiresAt,
		CreatedAt: m.CreatedAt,
	}
	if m.UserID != nil {
		res.UserID = uuid.NullUUID{UUID: *m.UserID, Valid: true}
	}
	if m.AnsweredAt != nil {
		res.AnsweredAt = sql.NullTime{Time: *m.AnsweredAt, Valid: true}
	}
	return res
}
