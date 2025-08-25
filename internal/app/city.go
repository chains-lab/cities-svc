package app

import (
	"context"

	"github.com/chains-lab/cities-svc/internal/app/entities"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

//tx, err := c.db.BeginTx(ctx, pgx.TxOptions{})
//if err != nil {
//	return err
//}
//defer tx.Rollback(ctx)
//
//qtx := c.queries.WithTx(tx)

//if err := tx.Commit(ctx); err != nil {
//	return problems.RaiseInternal(ctx, err)
//}

func CreateCity(
	ctx context.Context,
	CountryID uuid.UUID,
	Icon, Slug, TimeZone, Name, Language string,
	Zone orb.MultiPolygon,
) (models.City, error) {
	return models.City{}, nil
}
