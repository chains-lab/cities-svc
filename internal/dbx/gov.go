package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const cityGovTable = "city_governments"

type CityGov struct {
	ID        uuid.UUID      `db:"id"`
	UserID    uuid.UUID      `db:"user_id"`
	CityID    uuid.UUID      `db:"city_id"`
	Active    bool           `db:"active"`
	Role      string         `db:"role"`
	Label     sql.NullString `db:"label"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

type CityGovQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewCityGovQ(db *sql.DB) CityGovQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// Явный список колонок для стабильного Scan-порядка
	cols := []string{
		"id",
		"user_id",
		"city_id",
		"active",
		"role",
		"label",
		"created_at",
		"updated_at",
	}

	return CityGovQ{
		db:       db,
		selector: b.Select(cols...).From(cityGovTable),
		inserter: b.Insert(cityGovTable),
		updater:  b.Update(cityGovTable),
		deleter:  b.Delete(cityGovTable),
		counter:  b.Select("COUNT(*) AS count").From(cityGovTable),
	}
}

func (q CityGovQ) New() CityGovQ { return NewCityGovQ(q.db) }

// Insert: если ID == uuid.Nil — не ставим его => возьмётся DEFAULT из БД
func (q CityGovQ) Insert(ctx context.Context, in CityGov) error {
	values := map[string]interface{}{
		"user_id":    in.UserID,
		"city_id":    in.CityID,
		"active":     in.Active, // DEFAULT true в БД, но можно и явно
		"role":       in.Role,
		"label":      in.Label,     // NullString корректно передаётся
		"created_at": in.CreatedAt, // в схеме NOT NULL без DEFAULT — передаём сами
		"updated_at": in.UpdatedAt,
	}
	if in.ID != uuid.Nil {
		values["id"] = in.ID
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

// Get: возвращает первый попавшийся по текущим фильтрам
func (q CityGovQ) Get(ctx context.Context) (CityGov, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return CityGov{}, fmt.Errorf("building select query for %s: %w", cityGovTable, err)
	}

	var m CityGov
	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(
		&m.ID,
		&m.UserID,
		&m.CityID,
		&m.Active,
		&m.Role,
		&m.Label,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	return m, err
}

func (q CityGovQ) Select(ctx context.Context) ([]CityGov, error) {
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

	var out []CityGov
	for rows.Next() {
		var m CityGov
		if err := rows.Scan(
			&m.ID,
			&m.UserID,
			&m.CityID,
			&m.Active,
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
	Active    *bool
	Role      *string
	Label     *sql.NullString
	UpdatedAt *time.Time
}

func (q CityGovQ) Update(ctx context.Context, p UpdateCityGovParams) error {
	updates := map[string]interface{}{}

	if p.CityID != nil {
		updates["city_id"] = *p.CityID
	}
	if p.Active != nil {
		updates["active"] = *p.Active
	}
	if p.Role != nil {
		updates["role"] = *p.Role
	}
	if p.Label != nil {
		if p.Label.Valid {
			updates["label"] = p.Label
		} else {
			updates["label"] = nil
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

func (q CityGovQ) Delete(ctx context.Context) error {
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

func (q CityGovQ) FilterID(id uuid.UUID) CityGovQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})
	return q
}

func (q CityGovQ) FilterUserID(userID uuid.UUID) CityGovQ {
	q.selector = q.selector.Where(sq.Eq{"user_id": userID})
	q.deleter = q.deleter.Where(sq.Eq{"user_id": userID})
	q.updater = q.updater.Where(sq.Eq{"user_id": userID})
	q.counter = q.counter.Where(sq.Eq{"user_id": userID})
	return q
}

func (q CityGovQ) FilterCityID(cityID uuid.UUID) CityGovQ {
	q.selector = q.selector.Where(sq.Eq{"city_id": cityID})
	q.deleter = q.deleter.Where(sq.Eq{"city_id": cityID})
	q.updater = q.updater.Where(sq.Eq{"city_id": cityID})
	q.counter = q.counter.Where(sq.Eq{"city_id": cityID})
	return q
}

func (q CityGovQ) FilterRole(role ...string) CityGovQ {
	q.selector = q.selector.Where(sq.Eq{"role": role})
	q.deleter = q.deleter.Where(sq.Eq{"role": role})
	q.updater = q.updater.Where(sq.Eq{"role": role})
	q.counter = q.counter.Where(sq.Eq{"role": role})
	return q
}

func (q CityGovQ) FilterActive(active bool) CityGovQ {
	q.selector = q.selector.Where(sq.Eq{"active": active})
	q.deleter = q.deleter.Where(sq.Eq{"active": active})
	q.updater = q.updater.Where(sq.Eq{"active": active})
	q.counter = q.counter.Where(sq.Eq{"active": active})
	return q
}

func (q CityGovQ) FilterCountryID(countryID uuid.UUID) CityGovQ {
	cond := sq.Expr("city_id IN (SELECT id FROM city WHERE country_id = ?)", countryID)
	q.selector = q.selector.Where(cond)
	q.updater = q.updater.Where(cond)
	q.deleter = q.deleter.Where(cond)
	q.counter = q.counter.Where(cond)
	return q
}

func (q CityGovQ) FilterLabelLike(label string) CityGovQ {
	cond := sq.Expr("label ILIKE ?", fmt.Sprintf("%%%s%%", label))
	q.selector = q.selector.Where(cond)
	q.deleter = q.deleter.Where(cond)
	q.updater = q.updater.Where(cond)
	q.counter = q.counter.Where(cond)
	return q
}

func (q CityGovQ) OrderByRole(asc bool) CityGovQ {
	dir := "ASC"
	if !asc {
		dir = "DESC"
	}
	q.selector = q.selector.OrderBy("role " + dir)
	return q
}

func (q CityGovQ) OrderByCreatedAt(asc bool) CityGovQ {
	dir := "ASC"
	if !asc {
		dir = "DESC"
	}
	q.selector = q.selector.OrderBy("created_at " + dir)
	return q
}

func (q CityGovQ) OrderByUpdatedAt(asc bool) CityGovQ {
	dir := "ASC"
	if !asc {
		dir = "DESC"
	}
	q.selector = q.selector.OrderBy("updated_at " + dir)
	return q
}

func (q CityGovQ) Count(ctx context.Context) (uint64, error) {
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

func (q CityGovQ) Page(limit, offset uint64) CityGovQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	return q
}
