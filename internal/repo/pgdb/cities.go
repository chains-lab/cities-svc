package pgdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

const citiesTable = "city"

type City struct {
	ID        uuid.UUID
	CountryID string
	Point     orb.Point // [lon, lat]
	Status    string
	Name      string
	Icon      *string
	Slug      *string
	Timezone  string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CitiesQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	updater  sq.UpdateBuilder
	inserter sq.InsertBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewCitiesQ(db *sql.DB) CitiesQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return CitiesQ{
		db: db,
		selector: b.Select(
			"id",
			"country_id",
			"ST_X(point::geometry) AS point_lon",
			"ST_Y(point::geometry) AS point_lat",
			"status",
			"name",
			"icon",
			"slug",
			"timezone",
			"created_at",
			"updated_at",
		).From(citiesTable),
		updater:  b.Update(citiesTable),
		inserter: b.Insert(citiesTable),
		deleter:  b.Delete(citiesTable),
		counter:  b.Select("COUNT(*) AS count").From(citiesTable),
	}
}

func (q CitiesQ) New() CitiesQ { return NewCitiesQ(q.db) }

func scanCityRow(scanner interface{ Scan(dest ...any) error }) (City, error) {
	var (
		c        City
		lon, lat float64
	)
	if err := scanner.Scan(
		&c.ID,
		&c.CountryID,
		&lon,
		&lat,
		&c.Status,
		&c.Name,
		&c.Icon, // sql.NullString
		&c.Slug, // sql.NullString
		&c.Timezone,
		&c.CreatedAt,
		&c.UpdatedAt,
	); err != nil {
		return City{}, err
	}
	c.Point = orb.Point{lon, lat}
	return c, nil
}

func (q CitiesQ) Insert(ctx context.Context, in City) error {
	vals := map[string]any{
		"id":         in.ID,
		"country_id": in.CountryID,
		"point":      sq.Expr("ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography", in.Point[0], in.Point[1]),
		"status":     in.Status,
		"name":       in.Name,
		"timezone":   in.Timezone,
		"created_at": in.CreatedAt,
		"updated_at": in.UpdatedAt,
	}

	if in.Icon != nil {
		vals["icon"] = *in.Icon
	}
	if in.Slug != nil {
		vals["slug"] = *in.Slug
	}

	qry, args, err := q.inserter.SetMap(vals).ToSql()
	if err != nil {
		return fmt.Errorf("build insert %s: %w", citiesTable, err)
	}
	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, qry, args...)
	} else {
		_, err = q.db.ExecContext(ctx, qry, args...)
	}
	return err
}

func (q CitiesQ) Select(ctx context.Context) ([]City, error) {
	qry, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select %s: %w", citiesTable, err)
	}
	var rows *sql.Rows
	if tx, ok := TxFromCtx(ctx); ok {
		rows, err = tx.QueryContext(ctx, qry, args...)
	} else {
		rows, err = q.db.QueryContext(ctx, qry, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []City
	for rows.Next() {
		c, err := scanCityRow(rows)
		if err != nil {
			return nil, fmt.Errorf("scan %s: %w", citiesTable, err)
		}
		out = append(out, c)
	}
	return out, nil
}

func (q CitiesQ) Get(ctx context.Context) (City, error) {
	qry, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return City{}, fmt.Errorf("build select %s: %w", citiesTable, err)
	}
	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
		row = tx.QueryRowContext(ctx, qry, args...)
	} else {
		row = q.db.QueryRowContext(ctx, qry, args...)
	}
	return scanCityRow(row)
}

func (q CitiesQ) Update(ctx context.Context, updatedAt time.Time) error {
	q.updater = q.updater.Set("updated_at", updatedAt)

	query, args, err := q.updater.ToSql()
	if err != nil {
		return fmt.Errorf("building update query for %s: %w", citiesTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q CitiesQ) UpdateCountryID(countryID string) CitiesQ {
	q.updater = q.updater.Set("country_id", countryID)
	return q
}

func (q CitiesQ) UpdatePoint(point orb.Point) CitiesQ {
	q.updater = q.updater.Set(
		"point",
		sq.Expr("ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography", point[0], point[1]),
	)
	return q
}

func (q CitiesQ) UpdateStatus(status string) CitiesQ {
	q.updater = q.updater.Set("status", status)
	return q
}

func (q CitiesQ) UpdateName(name string) CitiesQ {
	q.updater = q.updater.Set("name", name)
	return q
}

func (q CitiesQ) UpdateIcon(icon sql.NullString) CitiesQ {
	if icon.Valid {
		q.updater = q.updater.Set("icon", icon.String)
	} else {
		q.updater = q.updater.Set("icon", nil)
	}
	return q
}

func (q CitiesQ) UpdateSlug(slug sql.NullString) CitiesQ {
	if slug.Valid {
		q.updater = q.updater.Set("slug", slug.String)
	} else {
		q.updater = q.updater.Set("slug", nil)
	}
	return q
}

func (q CitiesQ) UpdateTimezone(timezone string) CitiesQ {
	q.updater = q.updater.Set("timezone", timezone)
	return q
}

func (q CitiesQ) Delete(ctx context.Context) error {
	qry, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("build delete %s: %w", citiesTable, err)
	}
	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, qry, args...)
	} else {
		_, err = q.db.ExecContext(ctx, qry, args...)
	}
	return err
}

func (q CitiesQ) FilterID(id uuid.UUID) CitiesQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	return q
}

func (q CitiesQ) FilterCountryID(countryID ...string) CitiesQ {
	q.selector = q.selector.Where(sq.Eq{"country_id": countryID})
	q.counter = q.counter.Where(sq.Eq{"country_id": countryID})
	q.updater = q.updater.Where(sq.Eq{"country_id": countryID})
	q.deleter = q.deleter.Where(sq.Eq{"country_id": countryID})
	return q
}

func (q CitiesQ) FilterStatus(status ...string) CitiesQ {
	q.selector = q.selector.Where(sq.Eq{"status": status})
	q.counter = q.counter.Where(sq.Eq{"status": status})
	q.updater = q.updater.Where(sq.Eq{"status": status})
	q.deleter = q.deleter.Where(sq.Eq{"status": status})
	return q
}

func (q CitiesQ) FilterSlug(slug string) CitiesQ {
	q.selector = q.selector.Where(sq.Eq{"slug": slug})
	q.counter = q.counter.Where(sq.Eq{"slug": slug})
	q.updater = q.updater.Where(sq.Eq{"slug": slug})
	q.deleter = q.deleter.Where(sq.Eq{"slug": slug})
	return q
}

func (q CitiesQ) FilterNameLike(substr string) CitiesQ {
	cond := sq.Expr("name ILIKE ?", fmt.Sprintf("%%%s%%", substr))
	q.selector = q.selector.Where(cond)
	q.counter = q.counter.Where(cond)
	q.updater = q.updater.Where(cond)
	q.deleter = q.deleter.Where(cond)
	return q
}

func (q CitiesQ) FilterWithinRadiusMeters(point orb.Point, radiusM uint64) CitiesQ {
	p := sq.Expr("ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography", point[0], point[1])
	cond := sq.Expr("ST_DWithin(point, ?, ?)", p, radiusM)
	q.selector = q.selector.Where(cond)
	q.counter = q.counter.Where(cond)
	q.updater = q.updater.Where(cond)
	q.deleter = q.deleter.Where(cond)
	return q
}

func (q CitiesQ) OrderByAlphabetical(asc bool) CitiesQ {
	dir := "DESC"
	if asc {
		dir = "ASC"
	}
	orderExpr := fmt.Sprintf("name %s", dir)
	q.selector = q.selector.OrderBy(orderExpr)
	return q
}

func (q CitiesQ) OrderByNearest(point orb.Point, asc bool) CitiesQ {
	dir := "DESC"
	if asc {
		dir = "ASC"
	}

	q.selector = q.selector.SuffixExpr(
		sq.Expr(
			fmt.Sprintf("ORDER BY point <-> ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography %s", dir),
			point[0], point[1],
		),
	)
	return q
}

func (q CitiesQ) Page(limit, offset uint64) CitiesQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	return q
}

func (q CitiesQ) Count(ctx context.Context) (uint64, error) {
	qry, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build count %s: %w", citiesTable, err)
	}
	var n uint64
	if tx, ok := TxFromCtx(ctx); ok {
		err = tx.QueryRowContext(ctx, qry, args...).Scan(&n)
	} else {
		err = q.db.QueryRowContext(ctx, qry, args...).Scan(&n)
	}
	if err != nil {
		return 0, fmt.Errorf("scan count %s: %w", citiesTable, err)
	}
	return n, nil
}

func (q CitiesQ) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	_, ok := TxFromCtx(ctx)
	if ok {
		return fn(ctx)
	}

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		if err != nil {
			rbErr := tx.Rollback()
			if rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
				err = fmt.Errorf("tx err: %v; rollback err: %v", err, rbErr)
			}
		}
	}()

	ctxWithTx := context.WithValue(ctx, TxKey, tx)

	if err = fn(ctxWithTx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction failed: %v, rollback error: %v", err, rbErr)
		}
		return fmt.Errorf("transaction failed: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
