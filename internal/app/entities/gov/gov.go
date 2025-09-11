package gov

import (
	"database/sql"

	"github.com/chains-lab/cities-svc/internal/app/jwtmanager"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/config"
	"github.com/chains-lab/cities-svc/internal/dbx"
)

type Gov struct {
	govQ  dbx.GovQ
	jwt   jwtmanager.Manager
	query dbx.InviteQ
}

func NewGov(db *sql.DB, cfg config.Config) Gov {
	return Gov{
		govQ:  dbx.NewCityGovQ(db),
		jwt:   jwtmanager.NewManager(cfg),
		query: dbx.NewInviteQ(db),
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

func modelsFromDB(inv dbx.Invite, token string) models.Invite {
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
