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
	UserID    uuid.UUID      `db:"user_id"`
	CityID    uuid.UUID      `db:"city_id"`
	Role      string         `db:"role"`
	Label     sql.NullString `db:"label"`
	UpdatedAt time.Time      `db:"updated_at"`
	CreatedAt time.Time      `db:"created_at"`
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
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return CityGovQ{
		db:       db,
		selector: builder.Select("*").From(cityGovTable),
		inserter: builder.Insert(cityGovTable),
		updater:  builder.Update(cityGovTable),
		deleter:  builder.Delete(cityGovTable),
		counter:  builder.Select("COUNT(*) AS count").From(cityGovTable),
	}
}

func (q CityGovQ) New() CityGovQ {
	return NewCityGovQ(q.db)
}

func (q CityGovQ) Insert(ctx context.Context, input CityGov) error {
	values := map[string]interface{}{
		"user_id":    input.UserID,
		"city_id":    input.CityID,
		"role":       input.Role,
		"label":      input.Label,
		"updated_at": input.UpdatedAt,
		"created_at": input.CreatedAt,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building insertor query for table: %s: %w", cityGovTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q CityGovQ) Get(ctx context.Context) (CityGov, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return CityGov{}, fmt.Errorf("building selector query for table: %s: %w", cityGovTable, err)
	}

	var model CityGov
	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(
		&model.UserID,
		&model.CityID,
		&model.Role,
		&model.Label,
		&model.UpdatedAt,
		&model.CreatedAt,
	)

	return model, err
}

func (q CityGovQ) Select(ctx context.Context) ([]CityGov, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building selector query for table: %s: %w", cityGovTable, err)
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

	var models []CityGov
	for rows.Next() {
		var model CityGov
		if err := rows.Scan(
			&model.UserID,
			&model.CityID,
			&model.Role,
			&model.Label,
			&model.UpdatedAt,
			&model.CreatedAt,
		); err != nil {
			return nil, err
		}
		models = append(models, model)
	}

	return models, nil
}

type UpdateCityGovParams struct {
	CityID    *uuid.UUID      `db:"city_id"`
	Role      *string         `db:"role"`
	Label     *sql.NullString `db:"label"`
	UpdatedAt *time.Time      `db:"updated_at"`
}

func (q CityGovQ) Update(ctx context.Context, params UpdateCityGovParams) error {
	updates := map[string]interface{}{}

	if params.CityID != nil {
		updates["city_id"] = *params.CityID
	}
	if params.Role != nil {
		updates["role"] = *params.Role
	}
	if params.Label != nil {
		if params.Label.Valid {
			updates["label"] = params.Label
		} else {
			updates["label"] = nil
		}
	}
	if params.UpdatedAt != nil {
		updates["updated_at"] = *params.UpdatedAt
	} else if params.UpdatedAt == nil {
		updates["updated_at"] = time.Now().UTC()
	}

	query, args, err := q.updater.SetMap(updates).ToSql()
	if err != nil {
		return fmt.Errorf("building updater query for table: %s: %w", cityGovTable, err)
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
		return fmt.Errorf("building deleter query for table: %s: %w", cityGovTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q CityGovQ) FilterUserID(userID uuid.UUID) CityGovQ {
	q.selector = q.selector.Where(sq.Eq{"user_id": userID})
	q.counter = q.counter.Where(sq.Eq{"user_id": userID})
	q.deleter = q.deleter.Where(sq.Eq{"user_id": userID})
	q.updater = q.updater.Where(sq.Eq{"user_id": userID})

	return q
}

func (q CityGovQ) FilterCityID(cityID uuid.UUID) CityGovQ {
	q.selector = q.selector.Where(sq.Eq{"city_id": cityID})
	q.counter = q.counter.Where(sq.Eq{"city_id": cityID})
	q.deleter = q.deleter.Where(sq.Eq{"city_id": cityID})
	q.updater = q.updater.Where(sq.Eq{"city_id": cityID})

	return q
}

func (q CityGovQ) FilterRole(role ...string) CityGovQ {
	q.selector = q.selector.Where(sq.Eq{"role": role})
	q.counter = q.counter.Where(sq.Eq{"role": role})
	q.deleter = q.deleter.Where(sq.Eq{"role": role})
	q.updater = q.updater.Where(sq.Eq{"role": role})

	return q
}

func (q CityGovQ) FilterLabelLike(label string) CityGovQ {
	cond := sq.Expr("label ILIKE ?", fmt.Sprintf("%%%s%%", label))
	q.selector = q.selector.Where(cond)
	q.counter = q.counter.Where(cond)
	q.updater = q.updater.Where(cond)
	q.deleter = q.deleter.Where(cond)

	return q
}

func (q CityGovQ) OrderByRole(asc bool) CityGovQ {
	dir := "ASC"
	if !asc {
		dir = "DESC"
	}
	q.selector = q.selector.OrderBy("role " + dir)
	q.counter = q.counter.OrderBy("role " + dir)
	return q
}

func (q CityGovQ) OrderByCreatedAt(asc bool) CityGovQ {
	dir := "ASC"
	if !asc {
		dir = "DESC"
	}
	q.selector = q.selector.OrderBy("created_at " + dir)
	q.counter = q.counter.OrderBy("created_at " + dir)
	return q
}

func (q CityGovQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building counter query for table: %s: %w", cityGovTable, err)
	}

	var count uint64
	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("scanning count result: %w", err)
	}

	return count, nil
}

func (q CityGovQ) Page(limit, offset uint64) CityGovQ {
	q.counter = q.counter.Limit(limit).Offset(offset)
	q.selector = q.selector.Limit(limit).Offset(offset)
	return q
}
