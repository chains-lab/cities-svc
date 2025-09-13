package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

const citiesTable = "city"

type City struct {
	ID        uuid.UUID
	CountryID uuid.UUID
	Point     orb.Point // [lon, lat]
	Status    string
	Name      string
	Icon      sql.NullString // was string, now nullable
	Slug      sql.NullString // was string, now nullable
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
	var icon any
	if in.Icon.Valid {
		icon = in.Icon.String
	} else {
		icon = nil
	}
	var slug any
	if in.Slug.Valid {
		slug = in.Slug.String
	} else {
		slug = nil
	}

	vals := map[string]any{
		"id":         in.ID,
		"country_id": in.CountryID,
		"point":      sq.Expr("ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography", in.Point[0], in.Point[1]),
		"status":     in.Status,
		"name":       in.Name,
		"icon":       icon, // may be NULL
		"slug":       slug, // may be NULL
		"timezone":   in.Timezone,
		"created_at": in.CreatedAt,
		"updated_at": in.UpdatedAt,
	}

	qry, args, err := q.inserter.SetMap(vals).ToSql()
	if err != nil {
		return fmt.Errorf("build insert %s: %w", citiesTable, err)
	}
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
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
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
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
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, qry, args...)
	} else {
		row = q.db.QueryRowContext(ctx, qry, args...)
	}
	return scanCityRow(row)
}

type UpdateCityParams struct {
	CountryID *uuid.UUID
	Point     *orb.Point // [lon, lat]
	Status    *string
	Name      *string
	Icon      *sql.NullString // nullable
	Slug      *sql.NullString // nullable
	Timezone  *string
	UpdatedAt time.Time
}

func (q CitiesQ) Update(ctx context.Context, p UpdateCityParams) error {
	updates := map[string]any{}

	if p.CountryID != nil {
		updates["country_id"] = *p.CountryID
	}
	if p.Point != nil {
		pt := *p.Point
		updates["point"] = sq.Expr(
			"ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography",
			pt[0], pt[1],
		)
	}
	if p.Status != nil {
		updates["status"] = *p.Status
	}
	if p.Name != nil {
		updates["name"] = *p.Name
	}
	if p.Icon != nil {
		if p.Icon.Valid {
			updates["icon"] = p.Icon.String
		} else {
			updates["icon"] = nil // set NULL
		}
	}
	if p.Slug != nil {
		if p.Slug.Valid {
			updates["slug"] = p.Slug.String
		} else {
			updates["slug"] = nil // set NULL
		}
	}
	if p.Timezone != nil {
		updates["timezone"] = *p.Timezone
	}

	updates["updated_at"] = p.UpdatedAt

	if len(updates) == 1 { // только updated_at
		return nil
	}

	qry, args, err := q.updater.SetMap(updates).ToSql()
	if err != nil {
		return fmt.Errorf("build update %s: %w", citiesTable, err)
	}
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, qry, args...)
	} else {
		_, err = q.db.ExecContext(ctx, qry, args...)
	}
	return err
}

func (q CitiesQ) Delete(ctx context.Context) error {
	qry, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("build delete %s: %w", citiesTable, err)
	}
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, qry, args...)
	} else {
		_, err = q.db.ExecContext(ctx, qry, args...)
	}
	return err
}

// -------- Фильтры/сортировки --------

func (q CitiesQ) FilterID(id uuid.UUID) CitiesQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	return q
}

func (q CitiesQ) FilterCountryID(countryID ...uuid.UUID) CitiesQ {
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
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		err = tx.QueryRowContext(ctx, qry, args...).Scan(&n)
	} else {
		err = q.db.QueryRowContext(ctx, qry, args...).Scan(&n)
	}
	if err != nil {
		return 0, fmt.Errorf("scan count %s: %w", citiesTable, err)
	}
	return n, nil
}
