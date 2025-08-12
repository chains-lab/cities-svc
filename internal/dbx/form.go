package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const formsToCreateCityTable = "forms_to_create_city"

type FormToCreateCityModel struct {
	ID           uuid.UUID `db:"id"`
	Status       string    `db:"status"`
	CityName     string    `db:"city_name"`
	CountryID    uuid.UUID `db:"country_id"`
	InitiatorID  uuid.UUID `db:"initiator_id"`
	ContactEmail string    `db:"contact_email"`
	ContactPhone string    `db:"contact_phone"`
	Text         string    `db:"text"`
	UserRevID    uuid.UUID `db:"user_reviewed_id,omitempty"`
	CreatedAt    time.Time `db:"create_at"`
}

type FormToCreateCityQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewFormToCreateCityQ(db *sql.DB) FormToCreateCityQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return FormToCreateCityQ{
		db:       db,
		selector: builder.Select("*").From(formsToCreateCityTable),
		inserter: builder.Insert(formsToCreateCityTable),
		updater:  builder.Update(formsToCreateCityTable),
		deleter:  builder.Delete(formsToCreateCityTable),
		counter:  builder.Select("COUNT(*) AS count").From(formsToCreateCityTable),
	}
}

func (q FormToCreateCityQ) New() FormToCreateCityQ {
	return NewFormToCreateCityQ(q.db)
}

func (q FormToCreateCityQ) Insert(ctx context.Context, input FormToCreateCityModel) error {
	values := map[string]interface{}{
		"id":               input.ID,
		"status":           input.Status,
		"city_name":        input.CityName,
		"country_id":       input.CountryID,
		"initiator_id":     input.InitiatorID,
		"contact_email":    input.ContactEmail,
		"contact_phone":    input.ContactPhone,
		"text":             input.Text,
		"user_reviewed_id": input.UserRevID,
		"create_at":        input.CreatedAt,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building inserter query for table: %s: %w", formsToCreateCityTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q FormToCreateCityQ) Get(ctx context.Context) (FormToCreateCityModel, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return FormToCreateCityModel{}, fmt.Errorf("building selector query for table: %s: %w", formsToCreateCityTable, err)
	}

	var form FormToCreateCityModel
	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	err = row.Scan(
		&form.ID,
		&form.Status,
		&form.CityName,
		&form.CountryID,
		&form.InitiatorID,
		&form.ContactEmail,
		&form.ContactPhone,
		&form.Text,
		&form.UserRevID,
		&form.CreatedAt,
	)

	return form, err
}

func (q FormToCreateCityQ) Select(ctx context.Context) ([]FormToCreateCityModel, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building selector query for table: %s: %w", formsToCreateCityTable, err)
	}

	var forms []FormToCreateCityModel
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

	for rows.Next() {
		var form FormToCreateCityModel
		if err := rows.Scan(
			&form.ID,
			&form.Status,
			&form.CityName,
			&form.CountryID,
			&form.InitiatorID,
			&form.ContactEmail,
			&form.ContactPhone,
			&form.Text,
			&form.UserRevID,
			&form.CreatedAt,
		); err != nil {
			return nil, err
		}
		forms = append(forms, form)
	}

	return forms, nil
}

type UpdateFormToCreateCityInput struct {
	Status       *string
	ContactEmail *string
	ContactPhone *string
	Text         *string
	UserRevID    *uuid.UUID
}

func (q FormToCreateCityQ) Update(ctx context.Context, input UpdateFormToCreateCityInput) error {
	updates := map[string]interface{}{}
	if input.Status != nil {
		updates["status"] = *input.Status
	}
	if input.ContactEmail != nil {
		updates["contact_email"] = *input.ContactEmail
	}
	if input.ContactPhone != nil {
		updates["contact_phone"] = *input.ContactPhone
	}
	if input.Text != nil {
		updates["text"] = *input.Text
	}
	if input.UserRevID != nil && *input.UserRevID != uuid.Nil {
		updates["user_reviewed_id"] = *input.UserRevID
	}

	query, args, err := q.updater.SetMap(updates).ToSql()
	if err != nil {
		return fmt.Errorf("building updater query for table: %s: %w", formsToCreateCityTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q FormToCreateCityQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building deleter query for table: %s: %w", formsToCreateCityTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q FormToCreateCityQ) FilterID(ID uuid.UUID) FormToCreateCityQ {
	q.selector = q.selector.Where(sq.Eq{"id": ID})
	q.counter = q.counter.Where(sq.Eq{"id": ID})
	q.deleter = q.deleter.Where(sq.Eq{"id": ID})
	q.updater = q.updater.Where(sq.Eq{"id": ID})
	return q
}

func (q FormToCreateCityQ) FilterStatus(status string) FormToCreateCityQ {
	q.selector = q.selector.Where(sq.Eq{"status": status})
	q.counter = q.counter.Where(sq.Eq{"status": status})
	q.deleter = q.deleter.Where(sq.Eq{"status": status})
	q.updater = q.updater.Where(sq.Eq{"status": status})
	return q
}

func (q FormToCreateCityQ) FilterInitiatorID(initiatorID uuid.UUID) FormToCreateCityQ {
	q.selector = q.selector.Where(sq.Eq{"initiator_id": initiatorID})
	q.counter = q.counter.Where(sq.Eq{"initiator_id": initiatorID})
	q.deleter = q.deleter.Where(sq.Eq{"initiator_id": initiatorID})
	q.updater = q.updater.Where(sq.Eq{"initiator_id": initiatorID})
	return q
}

func (q FormToCreateCityQ) FilterCountryID(countryID uuid.UUID) FormToCreateCityQ {
	q.selector = q.selector.Where(sq.Eq{"country_id": countryID})
	q.counter = q.counter.Where(sq.Eq{"country_id": countryID})
	q.deleter = q.deleter.Where(sq.Eq{"country_id": countryID})
	q.updater = q.updater.Where(sq.Eq{"country_id": countryID})
	return q
}

func (q FormToCreateCityQ) CityNameLike(name string) FormToCreateCityQ {
	pattern := fmt.Sprintf("%%%s%%", name)
	q.selector = q.selector.Where("city_name ILIKE ?", pattern)
	q.counter = q.counter.Where("city_name ILIKE ?", pattern)
	return q
}

// Pagination
func (q FormToCreateCityQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for table: %s: %w", formsToCreateCityTable, err)
	}

	var count uint64
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&count)
	} else {
		err = q.db.QueryRowContext(ctx, query, args...).Scan(&count)
	}

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (q FormToCreateCityQ) Page(limit, offset uint64) FormToCreateCityQ {
	q.counter = q.counter.Limit(limit).Offset(offset)
	q.selector = q.selector.Limit(limit).Offset(offset)
	return q
}
