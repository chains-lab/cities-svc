package pgdb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const CityAdminsTable = "city_admins"

type CityAdmin struct {
	UserID    uuid.UUID `db:"user_id"`
	CityID    uuid.UUID `db:"city_id"`
	Role      string    `db:"role"`
	Position  *string   `db:"position"`
	Label     *string   `db:"label"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type CityAdminsQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewCityAdminsQ(db *sql.DB) CityAdminsQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	cols := []string{
		"user_id",
		"city_id",
		"role",
		"position",
		"label",
		"created_at",
		"updated_at",
	}

	return CityAdminsQ{
		db:       db,
		selector: b.Select(cols...).From(CityAdminsTable),
		inserter: b.Insert(CityAdminsTable),
		updater:  b.Update(CityAdminsTable),
		deleter:  b.Delete(CityAdminsTable),
		counter:  b.Select("COUNT(*) AS count").From(CityAdminsTable),
	}
}

func (q CityAdminsQ) New() CityAdminsQ { return NewCityAdminsQ(q.db) }

func (q CityAdminsQ) Insert(ctx context.Context, in CityAdmin) error {
	values := map[string]interface{}{
		"user_id": in.UserID,
		"city_id": in.CityID,
		"role":    in.Role,
	}

	if in.Position != nil {
		values["position"] = in.Position
	}
	if in.Label != nil {
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
		return fmt.Errorf("building insert query for %s: %w", CityAdminsTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q CityAdminsQ) Get(ctx context.Context) (CityAdmin, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return CityAdmin{}, fmt.Errorf("building select query for %s: %w", CityAdminsTable, err)
	}

	var m CityAdmin
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
		&m.Position,
		&m.Label,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	return m, err
}

func (q CityAdminsQ) Select(ctx context.Context) ([]CityAdmin, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", CityAdminsTable, err)
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

	var out []CityAdmin
	for rows.Next() {
		var m CityAdmin
		if err := rows.Scan(
			&m.UserID,
			&m.CityID,
			&m.Role,
			&m.Position,
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

func (q CityAdminsQ) Update(ctx context.Context, updatedAt time.Time) error {
	q.updater = q.updater.Set("updated_at", updatedAt)

	query, args, err := q.updater.ToSql()
	if err != nil {
		return fmt.Errorf("building update query for %s: %w", CityAdminsTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q CityAdminsQ) UpdateCityID(cityID uuid.UUID) CityAdminsQ {
	q.updater = q.updater.Set("city_id", cityID)
	return q
}

func (q CityAdminsQ) UpdateStatus(status string) CityAdminsQ {
	q.updater = q.updater.Set("status", status)
	return q
}

func (q CityAdminsQ) UpdateRole(role string) CityAdminsQ {
	q.updater = q.updater.Set("role", role)
	return q
}

func (q CityAdminsQ) UpdatePosition(position sql.NullString) CityAdminsQ {
	q.updater = q.updater.Set("position", position)
	return q
}

func (q CityAdminsQ) UpdateLabel(label sql.NullString) CityAdminsQ {
	q.updater = q.updater.Set("label", label)
	return q
}

func (q CityAdminsQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", CityAdminsTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q CityAdminsQ) FilterUserID(userID uuid.UUID) CityAdminsQ {
	q.selector = q.selector.Where(sq.Eq{"user_id": userID})
	q.deleter = q.deleter.Where(sq.Eq{"user_id": userID})
	q.updater = q.updater.Where(sq.Eq{"user_id": userID})
	q.counter = q.counter.Where(sq.Eq{"user_id": userID})
	return q
}

func (q CityAdminsQ) FilterCityID(cityID uuid.UUID) CityAdminsQ {
	q.selector = q.selector.Where(sq.Eq{"city_id": cityID})
	q.deleter = q.deleter.Where(sq.Eq{"city_id": cityID})
	q.updater = q.updater.Where(sq.Eq{"city_id": cityID})
	q.counter = q.counter.Where(sq.Eq{"city_id": cityID})
	return q
}

func (q CityAdminsQ) FilterRole(role ...string) CityAdminsQ {
	q.selector = q.selector.Where(sq.Eq{"role": role})
	q.deleter = q.deleter.Where(sq.Eq{"role": role})
	q.updater = q.updater.Where(sq.Eq{"role": role})
	q.counter = q.counter.Where(sq.Eq{"role": role})
	return q
}

// FilterCountryID Unsupported: this method is commented out and should not be used.
func (q CityAdminsQ) FilterCountryID(countryID string) CityAdminsQ {
	join := fmt.Sprintf("LEFT JOIN %s c ON c.id = cg.city_id", citiesTable)
	q.selector = q.selector.LeftJoin(join).Where(sq.Eq{"c.country_id": countryID})
	q.counter = q.counter.LeftJoin(join).Where(sq.Eq{"c.country_id": countryID})

	sub := sq.
		Select("1").
		From(citiesTable + " c").
		Where("c.id = " + CityAdminsTable + ".city_id").
		Where(sq.Eq{"c.country_id": countryID})

	subSQL, subArgs, _ := sub.ToSql()

	q.updater = q.updater.Where(sq.Expr("EXISTS ("+subSQL+")", subArgs...))
	q.deleter = q.deleter.Where(sq.Expr("EXISTS ("+subSQL+")", subArgs...))

	return q
}

func (q CityAdminsQ) FilterLabelLike(label string) CityAdminsQ {
	q.selector = q.selector.Where("label ILIKE ?", "%"+label+"%")
	q.deleter = q.deleter.Where("label ILIKE ?", "%"+label+"%")
	q.updater = q.updater.Where("label ILIKE ?", "%"+label+"%")
	q.counter = q.counter.Where("label ILIKE ?", "%"+label+"%")
	return q
}

func (q CityAdminsQ) FilterPositionLike(position string) CityAdminsQ {
	q.selector = q.selector.Where("position ILIKE ?", "%"+position+"%")
	q.deleter = q.deleter.Where("position ILIKE ?", "%"+position+"%")
	q.updater = q.updater.Where("position ILIKE ?", "%"+position+"%")
	q.counter = q.counter.Where("position ILIKE ?", "%"+position+"%")
	return q
}

func (q CityAdminsQ) OrderByRole(asc bool) CityAdminsQ {
	dir := "ASC"
	if !asc {
		dir = "DESC"
	}
	q.selector = q.selector.OrderBy("role " + dir)
	return q
}

func (q CityAdminsQ) OrderByCreatedAt(asc bool) CityAdminsQ {
	dir := "ASC"
	if !asc {
		dir = "DESC"
	}
	q.selector = q.selector.OrderBy("created_at " + dir)
	return q
}

func (q CityAdminsQ) OrderByUpdatedAt(asc bool) CityAdminsQ {
	dir := "ASC"
	if !asc {
		dir = "DESC"
	}
	q.selector = q.selector.OrderBy("updated_at " + dir)
	return q
}

func (q CityAdminsQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", CityAdminsTable, err)
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

func (q CityAdminsQ) Page(limit, offset uint64) CityAdminsQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	return q
}
