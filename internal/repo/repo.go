package repo

import (
	"context"
	"database/sql"

	"github.com/chains-lab/cities-svc/internal/repo/pgdb"
)

type Repo struct {
	sql SqlDB
}

func (r *Repo) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.sql.cities.New().Transaction(ctx, fn)
}

type SqlDB struct {
	cities    pgdb.CitiesQ
	invites   pgdb.InvitesQ
	cityAdmin pgdb.CityAdminsQ
}

func NewDatabase(db *sql.DB) *Repo {
	return &Repo{
		sql: SqlDB{
			cities:    pgdb.NewCitiesQ(db),
			invites:   pgdb.NewInvitesQ(db),
			cityAdmin: pgdb.NewCityAdminsQ(db),
		},
	}
}
