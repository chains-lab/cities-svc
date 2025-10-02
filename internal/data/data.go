package data

import (
	"context"
	"database/sql"

	"github.com/chains-lab/cities-svc/internal/data/pgdb"
)

type Database struct {
	sql SqlDB
}

func (d *Database) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.sql.cities.New().Transaction(ctx, fn)
}

type SqlDB struct {
	cities    pgdb.CitiesQ
	countries pgdb.CountriesQ
	invites   pgdb.InvitesQ
	cityMod   pgdb.CityModersQ
}

func NewDatabase(db *sql.DB) *Database {
	citySql := pgdb.NewCitiesQ(db)
	countrySql := pgdb.NewCountriesQ(db)
	inviteSql := pgdb.NewInvitesQ(db)
	cityModSql := pgdb.NewCityModersQ(db)

	return &Database{
		sql: SqlDB{
			cities:    citySql,
			countries: countrySql,
			invites:   inviteSql,
			cityMod:   cityModSql,
		},
	}
}
