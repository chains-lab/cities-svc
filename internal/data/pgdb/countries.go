package pgdb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const countriesTable = "countries"

type Country struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
type CountriesQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewCountriesQ(db *sql.DB) CountriesQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return CountriesQ{
		db:       db,
		selector: builder.Select("*").From(countriesTable),
		inserter: builder.Insert(countriesTable),
		updater:  builder.Update(countriesTable),
		deleter:  builder.Delete(countriesTable),
		counter:  builder.Select("COUNT(*) AS count").From(countriesTable),
	}
}

func (q CountriesQ) New() CountriesQ {
	return NewCountriesQ(q.db)
}

func (q CountriesQ) Insert(ctx context.Context, input Country) error {
	values := map[string]interface{}{
		"id":         input.ID,
		"name":       input.Name,
		"status":     input.Status,
		"created_at": input.CreatedAt,
		"updated_at": input.UpdatedAt,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return err
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q CountriesQ) Get(ctx context.Context) (Country, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return Country{}, fmt.Errorf("building selector query for table: %s: %w", countriesTable, err)
	}

	var model Country
	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(
		&model.ID,
		&model.Name,
		&model.Status,
		&model.UpdatedAt,
		&model.CreatedAt,
	)

	return model, err
}

func (q CountriesQ) Select(ctx context.Context) ([]Country, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building selector query for table: %s: %w", countriesTable, err)
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

	var models []Country
	for rows.Next() {
		var model Country
		if err := rows.Scan(
			&model.ID,
			&model.Name,
			&model.Status,
			&model.UpdatedAt,
			&model.CreatedAt,
		); err != nil {
			return nil, err
		}
		models = append(models, model)
	}

	return models, rows.Err()
}

func (q CountriesQ) Update(ctx context.Context, updatedAt time.Time) error {
	q.updater = q.updater.Set("updated_at", updatedAt)

	query, args, err := q.updater.ToSql()
	if err != nil {
		return fmt.Errorf("building update query for %s: %w", countriesTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q CountriesQ) UpdateName(name string) CountriesQ {
	q.updater = q.updater.Set("name", name)
	return q
}

func (q CountriesQ) UpdateStatus(status string) CountriesQ {
	q.updater = q.updater.Set("status", status)
	return q
}

func (q CountriesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building deleter query for table: %s: %w", countriesTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q CountriesQ) FilterID(ID uuid.UUID) CountriesQ {
	q.selector = q.selector.Where(sq.Eq{"id": ID})
	q.counter = q.counter.Where(sq.Eq{"id": ID})
	q.deleter = q.deleter.Where(sq.Eq{"id": ID})
	q.updater = q.updater.Where(sq.Eq{"id": ID})

	return q
}

func (q CountriesQ) FilterName(name string) CountriesQ {
	q.selector = q.selector.Where(sq.Eq{"name": name})
	q.counter = q.counter.Where(sq.Eq{"name": name})
	q.deleter = q.deleter.Where(sq.Eq{"name": name})
	q.updater = q.updater.Where(sq.Eq{"name": name})

	return q
}

func (q CountriesQ) FilterStatus(status ...string) CountriesQ {
	q.selector = q.selector.Where(sq.Eq{"status": status})
	q.counter = q.counter.Where(sq.Eq{"status": status})
	q.deleter = q.deleter.Where(sq.Eq{"status": status})
	q.updater = q.updater.Where(sq.Eq{"status": status})

	return q
}

func (q CountriesQ) FilterNameLike(name string) CountriesQ {
	likePattern := fmt.Sprintf("%%%s%%", name)
	q.selector = q.selector.Where(sq.ILike{"name": likePattern})
	q.counter = q.counter.Where(sq.ILike{"name": likePattern})
	q.deleter = q.deleter.Where(sq.ILike{"name": likePattern})
	q.updater = q.updater.Where(sq.ILike{"name": likePattern})

	return q
}

func (q CountriesQ) OrderByAlphabetical(asc bool) CountriesQ {
	if asc {
		q.selector = q.selector.OrderBy("name ASC")
	} else {
		q.selector = q.selector.OrderBy("name DESC")
	}
	return q
}

func (q CountriesQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for table: %s: %w", countriesTable, err)
	}

	var count uint64
	if tx, ok := TxFromCtx(ctx); ok {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&count)
	} else {
		err = q.db.QueryRowContext(ctx, query, args...).Scan(&count)
	}

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (q CountriesQ) Page(limit, offset uint64) CountriesQ {
	q.counter = q.counter.Limit(limit).Offset(offset)
	q.selector = q.selector.Limit(limit).Offset(offset)
	return q
}
