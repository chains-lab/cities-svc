package country

import (
	"database/sql"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/dbx"
)

type Country struct {
	countryQ dbx.CountryQ
}

func NewCountry(db *sql.DB) Country {
	return Country{
		countryQ: dbx.NewCountryQ(db),
	}
}

func countryFromDb(c dbx.Country) models.Country {
	return models.Country{
		ID:        c.ID,
		Name:      c.Name,
		Status:    c.Status,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
