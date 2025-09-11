package gov

import (
	"database/sql"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/dbx"
)

type Gov struct {
	govQ dbx.GovQ
}

func NewGov(db *sql.DB) Gov {
	return Gov{
		govQ: dbx.NewCityGovQ(db),
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
