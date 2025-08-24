package dbx

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

const citiesTable = "city"

type City struct {
	ID        uuid.UUID
	CountryID uuid.UUID
	Status    string
	Center    orb.Point        // [lon, lat]
	Boundary  orb.MultiPolygon // многоугольники границы
	Icon      string
	Slug      string
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
			"status",
			"ST_X(center) AS center_lon",
			"ST_Y(center) AS center_lat",
			"ST_AsGeoJSON(boundary) AS boundary_geojson",
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

func scanCityRow(scanner interface{ Scan(dest ...any) error }) (City, error) {
	var (
		c            City
		lon, lat     float64
		boundaryJSON string
	)

	if err := scanner.Scan(
		&c.ID,
		&c.CountryID,
		&c.Status,
		&lon,
		&lat,
		&boundaryJSON,
		&c.Icon,
		&c.Slug,
		&c.Timezone,
		&c.CreatedAt,
		&c.UpdatedAt,
	); err != nil {
		return City{}, err
	}

	c.Center = orb.Point{lon, lat}

	if boundaryJSON == "" {
		return c, nil
	}

	// 1) Ожидаемый путь: чистая Geometry
	if g, err := geojson.UnmarshalGeometry([]byte(boundaryJSON)); err == nil {
		switch geom := g.Geometry().(type) {
		case orb.MultiPolygon:
			c.Boundary = geom
		case orb.Polygon:
			c.Boundary = orb.MultiPolygon{geom}
		default:
			return City{}, fmt.Errorf("unexpected geometry type: %T", geom)
		}
		return c, nil
	}

	// 2) Фоллбэк: вдруг пришёл Feature
	if f, err := geojson.UnmarshalFeature([]byte(boundaryJSON)); err == nil {
		switch geom := f.Geometry.(type) { // NOTE: поле, не метод
		case orb.MultiPolygon:
			c.Boundary = geom
		case orb.Polygon:
			c.Boundary = orb.MultiPolygon{geom}
		default:
			return City{}, fmt.Errorf("unexpected feature geometry type: %T", geom)
		}
		return c, nil
	}

	// 3) Последний фоллбэк: проверить raw JSON на тип
	var raw map[string]any
	if err := json.Unmarshal([]byte(boundaryJSON), &raw); err == nil {
		if raw["type"] == "Polygon" || raw["type"] == "MultiPolygon" {
			if g, err := geojson.UnmarshalGeometry([]byte(boundaryJSON)); err == nil {
				switch geom := g.Geometry().(type) {
				case orb.MultiPolygon:
					c.Boundary = geom
				case orb.Polygon:
					c.Boundary = orb.MultiPolygon{geom}
				default:
					return City{}, fmt.Errorf("unexpected geometry type after fallback: %T", geom)
				}
				return c, nil
			}
		}
	}

	return City{}, fmt.Errorf("failed to decode boundary geojson")
}

func (q CitiesQ) applyConditions(conds ...sq.Sqlizer) CitiesQ {
	q.selector = q.selector.Where(conds)
	q.counter = q.counter.Where(conds)
	q.updater = q.updater.Where(conds)
	q.deleter = q.deleter.Where(conds)
	return q
}

func (q CitiesQ) New() CitiesQ { return NewCitiesQ(q.db) }

func (q CitiesQ) Insert(ctx context.Context, in City) error {
	// boundary: marshaling в GeoJSON geometry
	var boundaryJSON []byte
	{
		g := geojson.NewGeometry(in.Boundary) // допускает также Polygon внутри MultiPolygon
		var err error
		boundaryJSON, err = g.MarshalJSON()
		if err != nil {
			return fmt.Errorf("marshal boundary geojson: %w", err)
		}
	}

	vals := map[string]any{
		"id":         in.ID,
		"country_id": in.CountryID,
		"status":     in.Status,
		"center":     sq.Expr("ST_SetSRID(ST_MakePoint(?, ?), 4326)", in.Center[0], in.Center[1]),
		"boundary":   sq.Expr("ST_SetSRID(ST_GeomFromGeoJSON(?), 4326)", string(boundaryJSON)),
		"icon":       in.Icon,
		"slug":       in.Slug,
		"timezone":   in.Timezone,
		"created_at": in.CreatedAt,
		"updated_at": in.UpdatedAt,
	}

	qry, args, err := q.inserter.SetMap(vals).ToSql()
	if err != nil {
		return fmt.Errorf("build insert %s: %w", citiesTable, err)
	}
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
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
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
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
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, qry, args...)
	} else {
		row = q.db.QueryRowContext(ctx, qry, args...)
	}
	return scanCityRow(row)
}

func (q CitiesQ) Update(ctx context.Context, in map[string]any) error {
	vals := map[string]any{}

	if v, ok := in["country_id"]; ok {
		vals["country_id"] = v
	}
	if v, ok := in["status"]; ok {
		vals["status"] = v
	}
	// центр: ждём orb.Point в in["center"]
	if v, ok := in["center"]; ok {
		if p, ok2 := v.(orb.Point); ok2 {
			vals["center"] = sq.Expr("ST_SetSRID(ST_MakePoint(?, ?), 4326)", p[0], p[1])
		} else {
			return fmt.Errorf("center must be orb.Point")
		}
	}
	// boundary: ждём orb.MultiPolygon (или Polygon)
	if v, ok := in["boundary"]; ok {
		switch g := v.(type) {
		case orb.MultiPolygon:
			j, err := geojson.NewGeometry(g).MarshalJSON()
			if err != nil {
				return fmt.Errorf("marshal boundary: %w", err)
			}
			vals["boundary"] = sq.Expr("ST_SetSRID(ST_GeomFromGeoJSON(?), 4326)", string(j))
		case orb.Polygon:
			j, err := geojson.NewGeometry(g).MarshalJSON()
			if err != nil {
				return fmt.Errorf("marshal boundary: %w", err)
			}
			vals["boundary"] = sq.Expr("ST_SetSRID(ST_GeomFromGeoJSON(?), 4326)", string(j))
		default:
			return fmt.Errorf("boundary must be orb.MultiPolygon or orb.Polygon")
		}
	}
	if v, ok := in["icon"]; ok {
		vals["icon"] = v
	}
	if v, ok := in["slug"]; ok {
		vals["slug"] = v
	}
	if v, ok := in["timezone"]; ok {
		vals["timezone"] = v
	}
	if _, ok := in["updated_at"]; ok {
		vals["updated_at"] = in["updated_at"]
	} else {
		vals["updated_at"] = time.Now().UTC()
	}

	if len(vals) == 0 {
		return nil
	}

	qry, args, err := q.updater.SetMap(vals).ToSql()
	if err != nil {
		return fmt.Errorf("build update %s: %w", citiesTable, err)
	}
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
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
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, qry, args...)
	} else {
		_, err = q.db.ExecContext(ctx, qry, args...)
	}
	return err
}

func (q CitiesQ) FilterID(id uuid.UUID) CitiesQ {
	return q.applyConditions(sq.Eq{"id": id})
}

func (q CitiesQ) FilterCountryID(countryID uuid.UUID) CitiesQ {
	return q.applyConditions(sq.Eq{"country_id": countryID})
}

func (q CitiesQ) FilterStatus(status string) CitiesQ {
	return q.applyConditions(sq.Eq{"status": status})
}

func (q CitiesQ) FilterSlug(slug string) CitiesQ {
	return q.applyConditions(sq.Eq{"slug": slug})
}

func (q CitiesQ) FilterWithinRadiusMeters(lon, lat float64, radiusM uint64) CitiesQ {
	point := sq.Expr("ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography", lon, lat)
	cond := sq.Expr("ST_DWithin(center::geography, ?, ?)", point, radiusM)
	return q.applyConditions(cond)
}

func (q CitiesQ) OrderByDistance(lon, lat float64, asc bool) CitiesQ {
	asvVlue := "DESC"
	if asc {
		asvVlue = "ASC"
	}

	cond := sq.Expr("ST_Distance(center::geography, ST_SetSRID(ST_MakePoint(?, ?),4326)::geography) "+asvVlue, lon, lat)

	q.selector = q.selector.OrderByClause(cond)
	q.counter = q.counter.OrderByClause(cond)
	return q
}

func (q CitiesQ) Count(ctx context.Context) (uint64, error) {
	qry, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build count %s: %w", citiesTable, err)
	}
	var n uint64
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		err = tx.QueryRowContext(ctx, qry, args...).Scan(&n)
	} else {
		err = q.db.QueryRowContext(ctx, qry, args...).Scan(&n)
	}
	if err != nil {
		return 0, fmt.Errorf("scan count %s: %w", citiesTable, err)
	}
	return n, nil
}

func (q CitiesQ) Page(limit, offset uint64) CitiesQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	return q
}

func (q CitiesQ) Transaction(fn func(ctx context.Context) error) error {
	ctx := context.Background()

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	ctxWithTx := context.WithValue(ctx, txKey, tx)

	if err := fn(ctxWithTx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction failed: %v, rollback error: %v", err, rbErr)
		}
		return fmt.Errorf("transaction failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
