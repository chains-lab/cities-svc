package gov

import (
	"database/sql"

	"github.com/chains-lab/cities-svc/internal/app/jwtmanager"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/config"
	"github.com/chains-lab/cities-svc/internal/dbx"
)

type Gov struct {
	gov dbx.GovQ
	inv dbx.InviteQ
	jwt jwtmanager.Manager
}

func NewGov(db *sql.DB, cfg config.Config) Gov {
	return Gov{
		gov: dbx.NewCityGovQ(db),
		inv: dbx.NewInviteQ(db),
		jwt: jwtmanager.NewManager(cfg),
	}
}

func govFromDb(g dbx.Gov) models.Gov {
	res := models.Gov{
		UserID:    g.UserID,
		CityID:    g.CityID,
		Role:      g.Role,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
	}
	if g.Label.Valid {
		res.Label = &g.Label.String
	}

	return res
}

func inviteFromDB(inv dbx.Invite, token string) models.Invite {
	res := models.Invite{
		ID:        inv.ID,
		Status:    inv.Status,
		Role:      inv.Role,
		CityID:    inv.CityID,
		Token:     token,
		ExpiresAt: inv.ExpiresAt,
		CreatedAt: inv.CreatedAt,
	}
	if inv.UserID.Valid {
		res.UserID = &inv.UserID.UUID
	}
	if inv.AnsweredAt.Valid {
		res.AnsweredAt = &inv.AnsweredAt.Time
	}

	return res
}
