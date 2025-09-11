package invites

import (
	"database/sql"

	"github.com/chains-lab/cities-svc/internal/app/jwtmanager"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/config"
	"github.com/chains-lab/cities-svc/internal/dbx"
)

type Invite struct {
	jwt   jwtmanager.Manager
	query dbx.InviteQ
}

func NewInvites(pg *sql.DB, cfg config.Config) Invite {
	return Invite{
		jwt:   jwtmanager.NewManager(cfg),
		query: dbx.NewInviteQ(pg),
	}
}

func modelsFromDB(inv dbx.Invite) models.Invite {
	res := models.Invite{
		ID:          inv.ID,
		Status:      inv.Status,
		Role:        inv.Role,
		CityID:      inv.CityID,
		InitiatorID: inv.InitiatorID,
		ExpiresAt:   inv.ExpiresAt,
		CreatedAt:   inv.CreatedAt,
	}
	if inv.UserID.Valid {
		res.UserID = &inv.UserID.UUID
	}
	if inv.AnsweredAt.Valid {
		res.AnsweredAt = &inv.AnsweredAt.Time
	}

	return res
}
