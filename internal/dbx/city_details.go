package dbx

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const cityDetailsTable = "city_details"

type CityDetail struct {
	CityID      uuid.UUID `db:"city_id"`
	Language    string    `db:"language"`
	Name        string    `db:"name"`
	Description *string   `db:"description"`
}

type CityDetailsQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	updater  sq.UpdateBuilder
	inserter sq.InsertBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewCityDetailsQ(db *sql.DB) CityDetailsQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return CityDetailsQ{
		db: db,
		selector: b.Select(
			"city_id",
			"language",
			"name",
			"description",
		).From(cityDetailsTable),
		updater:  b.Update(cityDetailsTable),
		inserter: b.Insert(cityDetailsTable),
		deleter:  b.Delete(cityDetailsTable),
		counter:  b.Select("COUNT(*) AS count").From(cityDetailsTable),
	}
}

func scanCityDetailRow(scanner interface{ Scan(dest ...any) error }) (CityDetail, error) {
	var d CityDetail
	if err := scanner.Scan(
		&d.CityID,
		&d.Language,
		&d.Name,
		&d.Description,
	); err != nil {
		return CityDetail{}, err
	}
	return d, nil
}

func (q CityDetailsQ) applyConditions(conds ...sq.Sqlizer) CityDetailsQ {
	q.selector = q.selector.Where(conds)
	q.counter = q.counter.Where(conds)
	q.updater = q.updater.Where(conds)
	q.deleter = q.deleter.Where(conds)
	return q
}

func (q CityDetailsQ) New() CityDetailsQ { return NewCityDetailsQ(q.db) }

func (q CityDetailsQ) Insert(ctx context.Context, in CityDetail) error {
	vals := map[string]any{
		"city_id":     in.CityID,
		"language":    in.Language,
		"name":        in.Name,
		"description": in.Description,
	}
	qry, args, err := q.inserter.SetMap(vals).ToSql()
	if err != nil {
		return fmt.Errorf("build insert %s: %w", cityDetailsTable, err)
	}
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, qry, args...)
	} else {
		_, err = q.db.ExecContext(ctx, qry, args...)
	}
	return err
}

func (q CityDetailsQ) Get(ctx context.Context) (CityDetail, error) {
	qry, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return CityDetail{}, fmt.Errorf("build select %s: %w", cityDetailsTable, err)
	}
	var row *sql.Row
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, qry, args...)
	} else {
		row = q.db.QueryRowContext(ctx, qry, args...)
	}
	return scanCityDetailRow(row)
}

func (q CityDetailsQ) Select(ctx context.Context) ([]CityDetail, error) {
	qry, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select %s: %w", cityDetailsTable, err)
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

	var out []CityDetail
	for rows.Next() {
		d, err := scanCityDetailRow(rows)
		if err != nil {
			return nil, fmt.Errorf("scan %s: %w", cityDetailsTable, err)
		}
		out = append(out, d)
	}
	return out, nil
}

func (q CityDetailsQ) Update(ctx context.Context, in map[string]any) error {
	vals := map[string]any{}
	if v, ok := in["name"]; ok {
		vals["name"] = v
	}
	if v, ok := in["description"]; ok {
		vals["description"] = v
	}
	qry, args, err := q.updater.SetMap(vals).ToSql()
	if err != nil {
		return fmt.Errorf("build update %s: %w", cityDetailsTable, err)
	}
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, qry, args...)
	} else {
		_, err = q.db.ExecContext(ctx, qry, args...)
	}
	return err
}

func (q CityDetailsQ) Delete(ctx context.Context) error {
	qry, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("build delete %s: %w", cityDetailsTable, err)
	}
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, qry, args...)
	} else {
		_, err = q.db.ExecContext(ctx, qry, args...)
	}
	return err
}

func (q CityDetailsQ) FilterCityID(id uuid.UUID) CityDetailsQ {
	return q.applyConditions(sq.Eq{"city_id": id})
}

func (q CityDetailsQ) FilterLanguage(lang string) CityDetailsQ {
	return q.applyConditions(sq.Eq{"language": lang})
}

func (q CityDetailsQ) SearchName(name string) CityDetailsQ {
	return q.applyConditions(sq.Expr("name ILIKE ?", fmt.Sprintf("%%%s%%", name)))
}

func (q CityDetailsQ) Count(ctx context.Context) (uint64, error) {
	qry, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build count %s: %w", cityDetailsTable, err)
	}
	var n uint64
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		err = tx.QueryRowContext(ctx, qry, args...).Scan(&n)
	} else {
		err = q.db.QueryRowContext(ctx, qry, args...).Scan(&n)
	}
	return n, err
}

func (q CityDetailsQ) Page(limit, offset uint64) CityDetailsQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	return q
}
