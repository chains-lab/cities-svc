package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const invitesTable = "invites"

// Invite — модель строки из invites.
type Invite struct {
	ID          uuid.UUID     `db:"id"`
	Status      string        `db:"status"` // 'sent' | 'accepted' | 'rejected'
	Role        string        `db:"role"`   // enum employee_roles
	CityID      uuid.UUID     `db:"city_id"`
	InitiatorID uuid.UUID     `db:"initiator_id"` // пользователь-инициатор (из внешнего сервиса)
	UserID      uuid.NullUUID `db:"user_id"`      // может быть NULL до акцепта
	AnsweredAt  sql.NullTime  `db:"answered_at"`  // NULL для sent
	ExpiresAt   time.Time     `db:"expires_at"`
	CreatedAt   time.Time     `db:"created_at"`
}

// InviteQ — билдер запросов к invites (в стиле твоих Q-структур).
type InviteQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewInviteQ(db *sql.DB) InviteQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	cols := []string{
		"id",
		"status",
		"role",
		"city_id",
		"initiator_id",
		"user_id",
		"answered_at",
		"expires_at",
		"created_at",
	}
	return InviteQ{
		db:       db,
		selector: b.Select(cols...).From(invitesTable),
		inserter: b.Insert(invitesTable),
		updater:  b.Update(invitesTable),
		deleter:  b.Delete(invitesTable),
		counter:  b.Select("COUNT(*) AS count").From(invitesTable),
	}
}

func (q InviteQ) New() InviteQ { return NewInviteQ(q.db) }

// Insert — вставка новой записи. Если ID == uuid.Nil, возьмётся из переданного значения (требуется задать снаружи).
// created_at/updated_at можно не заполнять — если в схеме стоят DEFAULT, но ты их явно задаёшь в других местах.
func (q InviteQ) Insert(ctx context.Context, in Invite) error {
	values := map[string]interface{}{
		"id":           in.ID,
		"status":       in.Status,
		"role":         in.Role,
		"city_id":      in.CityID,
		"initiator_id": in.InitiatorID,
		"expires_at":   in.ExpiresAt,
	}

	if in.UserID.Valid {
		values["user_id"] = in.UserID
	}
	if in.AnsweredAt.Valid {
		values["answered_at"] = in.AnsweredAt
	}
	if !in.CreatedAt.IsZero() {
		values["created_at"] = in.CreatedAt
	}

	sqlStr, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("build insert %s: %w", invitesTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, sqlStr, args...)
	} else {
		_, err = q.db.ExecContext(ctx, sqlStr, args...)
	}
	return err
}

// Get — вернуть одну запись по текущим фильтрам.
func (q InviteQ) Get(ctx context.Context) (Invite, error) {
	sqlStr, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return Invite{}, fmt.Errorf("build select %s: %w", invitesTable, err)
	}

	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, sqlStr, args...)
	} else {
		row = q.db.QueryRowContext(ctx, sqlStr, args...)
	}

	var m Invite
	var userID uuid.NullUUID
	var answeredAt sql.NullTime

	if err := row.Scan(
		&m.ID,
		&m.Status,
		&m.Role,
		&m.CityID,
		&m.InitiatorID,
		&userID,
		&answeredAt,
		&m.ExpiresAt,
		&m.CreatedAt,
	); err != nil {
		return Invite{}, err
	}

	if userID.Valid {
		m.UserID = userID
	}
	if answeredAt.Valid {
		m.AnsweredAt = answeredAt
	}

	return m, nil
}

// Select — выбрать множество по текущим фильтрам.
func (q InviteQ) Select(ctx context.Context) ([]Invite, error) {
	sqlStr, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select %s: %w", invitesTable, err)
	}

	var rows *sql.Rows
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		rows, err = tx.QueryContext(ctx, sqlStr, args...)
	} else {
		rows, err = q.db.QueryContext(ctx, sqlStr, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Invite
	for rows.Next() {
		var m Invite
		var userID uuid.NullUUID
		var answeredAt sql.NullTime
		if err := rows.Scan(
			&m.ID,
			&m.Status,
			&m.Role,
			&m.CityID,
			&m.InitiatorID,
			&userID,
			&answeredAt,
			&m.ExpiresAt,
			&m.CreatedAt,
		); err != nil {
			return nil, err
		}
		if userID.Valid {
			m.UserID = userID
		}
		if answeredAt.Valid {
			m.AnsweredAt = answeredAt
		}
		out = append(out, m)
	}
	return out, nil
}

// UpdateInviteParams — частичное обновление (DAL без бизнес-логики).
type UpdateInviteParams struct {
	Status      *string
	Role        *string
	CityID      *uuid.UUID
	InitiatorID *uuid.UUID
	UserID      *uuid.NullUUID // двойной указатель, чтобы различать "не менять", "установить", "обнулить(nil)"
	AnsweredAt  *sql.NullTime  // аналогично: nil — не менять; *nil — записать NULL; *&t — записать t
	ExpiresAt   *time.Time
}

// Update — выполнит UPDATE по текущим фильтрам.
// ВАЖНО: как и в твоих Q, без вызова Filter* получишь UPDATE всех строк — ответственность на слое выше.
func (q InviteQ) Update(ctx context.Context, p UpdateInviteParams) error {
	updates := map[string]interface{}{}
	if p.Role != nil {
		updates["role"] = *p.Role
	}
	if p.Status != nil {
		updates["status"] = *p.Status
	}
	if p.CityID != nil {
		updates["city_id"] = *p.CityID
	}
	if p.InitiatorID != nil {
		updates["initiator_id"] = *p.InitiatorID
	}
	if p.UserID != nil { // задано желание изменить user_id
		if p.UserID.Valid {
			updates["user_id"] = p.UserID.UUID
		} else {
			updates["user_id"] = nil
		}
	}
	if p.AnsweredAt != nil {
		if p.AnsweredAt.Valid {
			updates["answered_at"] = p.AnsweredAt.Time
		} else {
			updates["answered_at"] = nil
		}
	}
	if p.ExpiresAt != nil {
		updates["expires_at"] = *p.ExpiresAt
	}

	if len(updates) == 0 {
		return nil
	}

	sqlStr, args, err := q.updater.SetMap(updates).ToSql()
	if err != nil {
		return fmt.Errorf("build update %s: %w", invitesTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, sqlStr, args...)
	} else {
		_, err = q.db.ExecContext(ctx, sqlStr, args...)
	}
	return err
}

func (q InviteQ) Delete(ctx context.Context) error {
	sqlStr, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("build delete %s: %w", invitesTable, err)
	}
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, sqlStr, args...)
	} else {
		_, err = q.db.ExecContext(ctx, sqlStr, args...)
	}
	return err
}

// --------- Фильтры ---------

func (q InviteQ) FilterID(id uuid.UUID) InviteQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})
	return q
}

func (q InviteQ) FilterCityID(cityID uuid.UUID) InviteQ {
	q.selector = q.selector.Where(sq.Eq{"city_id": cityID})
	q.updater = q.updater.Where(sq.Eq{"city_id": cityID})
	q.deleter = q.deleter.Where(sq.Eq{"city_id": cityID})
	q.counter = q.counter.Where(sq.Eq{"city_id": cityID})
	return q
}

func (q InviteQ) FilterInitiatorID(initiatorID uuid.UUID) InviteQ {
	q.selector = q.selector.Where(sq.Eq{"initiator_id": initiatorID})
	q.updater = q.updater.Where(sq.Eq{"initiator_id": initiatorID})
	q.deleter = q.deleter.Where(sq.Eq{"initiator_id": initiatorID})
	q.counter = q.counter.Where(sq.Eq{"initiator_id": initiatorID})
	return q
}

func (q InviteQ) FilterUserID(userID uuid.UUID) InviteQ {
	q.selector = q.selector.Where(sq.Eq{"user_id": userID})
	q.updater = q.updater.Where(sq.Eq{"user_id": userID})
	q.deleter = q.deleter.Where(sq.Eq{"user_id": userID})
	q.counter = q.counter.Where(sq.Eq{"user_id": userID})
	return q
}

func (q InviteQ) FilterStatus(status ...string) InviteQ {
	q.selector = q.selector.Where(sq.Eq{"status": status})
	q.updater = q.updater.Where(sq.Eq{"status": status})
	q.deleter = q.deleter.Where(sq.Eq{"status": status})
	q.counter = q.counter.Where(sq.Eq{"status": status})
	return q
}

func (q InviteQ) FilterRole(role ...string) InviteQ {
	q.selector = q.selector.Where(sq.Eq{"role": role})
	q.updater = q.updater.Where(sq.Eq{"role": role})
	q.deleter = q.deleter.Where(sq.Eq{"role": role})
	q.counter = q.counter.Where(sq.Eq{"role": role})
	return q
}

func (q InviteQ) FilterExpiresBefore(t time.Time) InviteQ {
	q.selector = q.selector.Where(sq.LtOrEq{"expires_at": t})
	q.updater = q.updater.Where(sq.LtOrEq{"expires_at": t})
	q.deleter = q.deleter.Where(sq.LtOrEq{"expires_at": t})
	q.counter = q.counter.Where(sq.LtOrEq{"expires_at": t})
	return q
}

func (q InviteQ) FilterExpiresAfter(t time.Time) InviteQ {
	q.selector = q.selector.Where(sq.Gt{"expires_at": t})
	q.updater = q.updater.Where(sq.Gt{"expires_at": t})
	q.deleter = q.deleter.Where(sq.Gt{"expires_at": t})
	q.counter = q.counter.Where(sq.Gt{"expires_at": t})
	return q
}

func (q InviteQ) FilterAnswered(answered bool) InviteQ {
	if answered {
		q.selector = q.selector.Where("answered_at IS NOT NULL")
		q.updater = q.updater.Where("answered_at IS NOT NULL")
		q.deleter = q.deleter.Where("answered_at IS NOT NULL")
		q.counter = q.counter.Where("answered_at IS NOT NULL")
	} else {
		q.selector = q.selector.Where("answered_at IS NULL")
		q.updater = q.updater.Where("answered_at IS NULL")
		q.deleter = q.deleter.Where("answered_at IS NULL")
		q.counter = q.counter.Where("answered_at IS NULL")
	}
	return q
}

func (q InviteQ) FilterCreatedBetween(from, to time.Time) InviteQ {
	q.selector = q.selector.Where(sq.And{
		sq.GtOrEq{"created_at": from},
		sq.LtOrEq{"created_at": to},
	})
	q.updater = q.updater.Where(sq.And{
		sq.GtOrEq{"created_at": from},
		sq.LtOrEq{"created_at": to},
	})
	q.deleter = q.deleter.Where(sq.And{
		sq.GtOrEq{"created_at": from},
		sq.LtOrEq{"created_at": to},
	})
	q.counter = q.counter.Where(sq.And{
		sq.GtOrEq{"created_at": from},
		sq.LtOrEq{"created_at": to},
	})
	return q
}

func (q InviteQ) OrderByCreatedAt(asc bool) InviteQ {
	dir := "ASC"
	if !asc {
		dir = "DESC"
	}
	q.selector = q.selector.OrderBy("created_at " + dir)
	return q
}

func (q InviteQ) OrderByUpdatedAt(asc bool) InviteQ {
	dir := "ASC"
	if !asc {
		dir = "DESC"
	}
	q.selector = q.selector.OrderBy("updated_at " + dir)
	return q
}

func (q InviteQ) OrderByExpiresAt(asc bool) InviteQ {
	dir := "ASC"
	if !asc {
		dir = "DESC"
	}
	q.selector = q.selector.OrderBy("expires_at " + dir)
	return q
}

func (q InviteQ) Count(ctx context.Context) (uint64, error) {
	sqlStr, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build count %s: %w", invitesTable, err)
	}

	var n uint64
	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, sqlStr, args...)
	} else {
		row = q.db.QueryRowContext(ctx, sqlStr, args...)
	}
	if err := row.Scan(&n); err != nil {
		return 0, fmt.Errorf("scan count %s: %w", invitesTable, err)
	}
	return n, nil
}

func (q InviteQ) Page(limit, offset uint64) InviteQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	return q
}
