package pgdb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const cityModersTable = "city_moders"

type CityMod struct {
	UserID    uuid.UUID      `db:"user_id"`
	CityID    uuid.UUID      `db:"city_id"`
	Role      string         `db:"role"`
	Label     sql.NullString `db:"label"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

type CityModersQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewCityModersQ(db *sql.DB) CityModersQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	cols := []string{
		"user_id",
		"city_id",
		"role",
		"label",
		"created_at",
		"updated_at",
	}

	return CityModersQ{
		db:       db,
		selector: b.Select(cols...).From(cityModersTable),
		inserter: b.Insert(cityModersTable),
		updater:  b.Update(cityModersTable),
		deleter:  b.Delete(cityModersTable),
		counter:  b.Select("COUNT(*) AS count").From(cityModersTable),
	}
}

func (q CityModersQ) New() CityModersQ { return NewCityModersQ(q.db) }

func (q CityModersQ) Insert(ctx context.Context, in CityMod) error {
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
		return fmt.Errorf("building insert query for %s: %w", cityModersTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q CityModersQ) Get(ctx context.Context) (CityMod, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return CityMod{}, fmt.Errorf("building select query for %s: %w", cityModersTable, err)
	}

	var m CityMod
	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
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

func (q CityModersQ) Select(ctx context.Context) ([]CityMod, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", cityModersTable, err)
	}

	var rows *sql.Rows
	if tx, ok := TxFromCtx(ctx); ok {
		rows, err = tx.QueryContext(ctx, query, args...)
	} else {
		rows, err = q.db.QueryContext(ctx, query, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []CityMod
	for rows.Next() {
		var m CityMod
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

func (q CityModersQ) Update(ctx context.Context, updatedAt time.Time) error {
	q.updater = q.updater.Set("updated_at", updatedAt)

	query, args, err := q.updater.ToSql()
	if err != nil {
		return fmt.Errorf("building update query for %s: %w", cityModersTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q CityModersQ) UpdateCityID(cityID uuid.UUID) CityModersQ {
	q.updater = q.updater.Set("city_id", cityID)
	return q
}

func (q CityModersQ) UpdateStatus(status string) CityModersQ {
	q.updater = q.updater.Set("status", status)
	return q
}

func (q CityModersQ) UpdateRole(role string) CityModersQ {
	q.updater = q.updater.Set("role", role)
	return q
}

func (q CityModersQ) UpdateLabel(label sql.NullString) CityModersQ {
	q.updater = q.updater.Set("label", label)
	return q
}

func (q CityModersQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", cityModersTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q CityModersQ) FilterUserID(userID uuid.UUID) CityModersQ {
	q.selector = q.selector.Where(sq.Eq{"user_id": userID})
	q.deleter = q.deleter.Where(sq.Eq{"user_id": userID})
	q.updater = q.updater.Where(sq.Eq{"user_id": userID})
	q.counter = q.counter.Where(sq.Eq{"user_id": userID})
	return q
}

func (q CityModersQ) FilterCityID(cityID uuid.UUID) CityModersQ {
	q.selector = q.selector.Where(sq.Eq{"city_id": cityID})
	q.deleter = q.deleter.Where(sq.Eq{"city_id": cityID})
	q.updater = q.updater.Where(sq.Eq{"city_id": cityID})
	q.counter = q.counter.Where(sq.Eq{"city_id": cityID})
	return q
}

func (q CityModersQ) FilterRole(role ...string) CityModersQ {
	q.selector = q.selector.Where(sq.Eq{"role": role})
	q.deleter = q.deleter.Where(sq.Eq{"role": role})
	q.updater = q.updater.Where(sq.Eq{"role": role})
	q.counter = q.counter.Where(sq.Eq{"role": role})
	return q
}

func (q CityModersQ) FilterCountryID(countryID uuid.UUID) CityModersQ {
	join := fmt.Sprintf("LEFT JOIN %s c ON c.id = cg.city_id", citiesTable)
	q.selector = q.selector.LeftJoin(join).Where(sq.Eq{"c.country_id": countryID})
	q.counter = q.counter.LeftJoin(join).Where(sq.Eq{"c.country_id": countryID})

	sub := sq.
		Select("1").
		From(citiesTable + " c").
		Where("c.id = " + cityModersTable + ".city_id").
		Where(sq.Eq{"c.country_id": countryID})

	subSQL, subArgs, _ := sub.ToSql()

	q.updater = q.updater.Where(sq.Expr("EXISTS ("+subSQL+")", subArgs...))
	q.deleter = q.deleter.Where(sq.Expr("EXISTS ("+subSQL+")", subArgs...))

	return q
}

func (q CityModersQ) FilterLabelLike(label string) CityModersQ {
	q.selector = q.selector.Where("label ILIKE ?", "%"+label+"%")
	q.deleter = q.deleter.Where("label ILIKE ?", "%"+label+"%")
	q.updater = q.updater.Where("label ILIKE ?", "%"+label+"%")
	q.counter = q.counter.Where("label ILIKE ?", "%"+label+"%")
	return q
}

func (q CityModersQ) OrderByRole(asc bool) CityModersQ {
	dir := "ASC"
	if !asc {
		dir = "DESC"
	}
	q.selector = q.selector.OrderBy("role " + dir)
	return q
}

func (q CityModersQ) OrderByCreatedAt(asc bool) CityModersQ {
	dir := "ASC"
	if !asc {
		dir = "DESC"
	}
	q.selector = q.selector.OrderBy("created_at " + dir)
	return q
}

func (q CityModersQ) OrderByUpdatedAt(asc bool) CityModersQ {
	dir := "ASC"
	if !asc {
		dir = "DESC"
	}
	q.selector = q.selector.OrderBy("updated_at " + dir)
	return q
}

func (q CityModersQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", cityModersTable, err)
	}

	var n uint64
	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	if err := row.Scan(&n); err != nil {
		return 0, fmt.Errorf("scanning count result: %w", err)
	}
	return n, nil
}

func (q CityModersQ) Page(limit, offset uint64) CityModersQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	return q
}
