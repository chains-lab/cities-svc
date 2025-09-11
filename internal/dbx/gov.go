package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const cityGovTable = "city_govs"

type Gov struct {
	UserID    uuid.UUID      `db:"user_id"`
	CityID    uuid.UUID      `db:"city_id"`
	Role      string         `db:"role"`
	Label     sql.NullString `db:"label"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

type GovQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewCityGovQ(db *sql.DB) GovQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	cols := []string{
		"user_id",
		"city_id",
		"role",
		"label",
		"created_at",
		"updated_at",
	}

	return GovQ{
		db:       db,
		selector: b.Select(cols...).From(cityGovTable),
		inserter: b.Insert(cityGovTable),
		updater:  b.Update(cityGovTable),
		deleter:  b.Delete(cityGovTable),
		counter:  b.Select("COUNT(*) AS count").From(cityGovTable),
	}
}

func (q GovQ) New() GovQ { return NewCityGovQ(q.db) }

func (q GovQ) Insert(ctx context.Context, in Gov) error {
	values := map[string]interface{}{
		"user_id": in.UserID,
		"city_id": in.CityID,
		"role":    in.Role,
	}
	if in.Label.Valid {
		values["label"] = in.Label
	}
	if !in.CreatedAt.IsZero() {
		values["created_at"] = in.CreatedAt
	}
	if !in.UpdatedAt.IsZero() {
		values["updated_at"] = in.UpdatedAt
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for %s: %w", cityGovTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q GovQ) Get(ctx context.Context) (Gov, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return Gov{}, fmt.Errorf("building select query for %s: %w", cityGovTable, err)
	}

	var m Gov
	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(
		&m.UserID,
		&m.CityID,
		&m.Role,
		&m.Label,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	return m, err
}

func (q GovQ) Select(ctx context.Context) ([]Gov, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", cityGovTable, err)
	}

	var rows *sql.Rows
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		rows, err = tx.QueryContext(ctx, query, args...)
	} else {
		rows, err = q.db.QueryContext(ctx, query, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Gov
	for rows.Next() {
		var m Gov
		if err := rows.Scan(
			&m.UserID,
			&m.CityID,
			&m.Role,
			&m.Label,
			&m.CreatedAt,
			&m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, nil
}

type UpdateCityGovParams struct {
	CityID    *uuid.UUID
	Status    *string
	Role      *string
	Label     *sql.NullString
	UpdatedAt *time.Time
}

func (q GovQ) Update(ctx context.Context, p UpdateCityGovParams) error {
	updates := map[string]interface{}{}

	if p.CityID != nil {
		updates["city_id"] = *p.CityID
	}
	if p.Role != nil {
		updates["role"] = *p.Role
	}
	if p.Label.Valid {
		if p.Label.String == "" {
			updates["label"] = nil
		} else {
			updates["label"] = p.Label
		}
	}
	if p.UpdatedAt != nil {
		updates["updated_at"] = *p.UpdatedAt
	} else {
		updates["updated_at"] = time.Now().UTC()
	}

	if len(updates) == 0 {
		return nil
	}

	query, args, err := q.updater.SetMap(updates).ToSql()
	if err != nil {
		return fmt.Errorf("building update query for %s: %w", cityGovTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q GovQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", cityGovTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q GovQ) FilterUserID(userID uuid.UUID) GovQ {
	q.selector = q.selector.Where(sq.Eq{"user_id": userID})
	q.deleter = q.deleter.Where(sq.Eq{"user_id": userID})
	q.updater = q.updater.Where(sq.Eq{"user_id": userID})
	q.counter = q.counter.Where(sq.Eq{"user_id": userID})
	return q
}

func (q GovQ) FilterCityID(cityID uuid.UUID) GovQ {
	q.selector = q.selector.Where(sq.Eq{"city_id": cityID})
	q.deleter = q.deleter.Where(sq.Eq{"city_id": cityID})
	q.updater = q.updater.Where(sq.Eq{"city_id": cityID})
	q.counter = q.counter.Where(sq.Eq{"city_id": cityID})
	return q
}

func (q GovQ) FilterRole(role ...string) GovQ {
	q.selector = q.selector.Where(sq.Eq{"role": role})
	q.deleter = q.deleter.Where(sq.Eq{"role": role})
	q.updater = q.updater.Where(sq.Eq{"role": role})
	q.counter = q.counter.Where(sq.Eq{"role": role})
	return q
}

func (q GovQ) FilterCountryID(countryID uuid.UUID) GovQ {
	join := fmt.Sprintf("LEFT JOIN %s c ON c.id = cg.city_id", citiesTable)
	q.selector = q.selector.LeftJoin(join).Where(sq.Eq{"c.country_id": countryID})
	q.counter = q.counter.LeftJoin(join).Where(sq.Eq{"c.country_id": countryID})

	sub := sq.
		Select("1").
		From(citiesTable + " c").
		Where("c.id = city_govs.city_id").
		Where(sq.Eq{"c.country_id": countryID})

	subSQL, subArgs, _ := sub.ToSql()

	q.updater = q.updater.Where(sq.Expr("EXISTS ("+subSQL+")", subArgs...))
	q.deleter = q.deleter.Where(sq.Expr("EXISTS ("+subSQL+")", subArgs...))

	return q
}

func (q GovQ) FilterLabelLike(label string) GovQ {
	q.selector = q.selector.Where("label ILIKE ?", "%"+label+"%")
	q.deleter = q.deleter.Where("label ILIKE ?", "%"+label+"%")
	q.updater = q.updater.Where("label ILIKE ?", "%"+label+"%")
	q.counter = q.counter.Where("label ILIKE ?", "%"+label+"%")
	return q
}

func (q GovQ) OrderByRole(asc bool) GovQ {
	dir := "ASC"
	if !asc {
		dir = "DESC"
	}
	q.selector = q.selector.OrderBy("role " + dir)
	return q
}

func (q GovQ) OrderByCreatedAt(asc bool) GovQ {
	dir := "ASC"
	if !asc {
		dir = "DESC"
	}
	q.selector = q.selector.OrderBy("created_at " + dir)
	return q
}

func (q GovQ) OrderByUpdatedAt(asc bool) GovQ {
	dir := "ASC"
	if !asc {
		dir = "DESC"
	}
	q.selector = q.selector.OrderBy("updated_at " + dir)
	return q
}

func (q GovQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", cityGovTable, err)
	}

	var n uint64
	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	if err := row.Scan(&n); err != nil {
		return 0, fmt.Errorf("scanning count result: %w", err)
	}
	return n, nil
}

func (q GovQ) Page(limit, offset uint64) GovQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	return q
}
