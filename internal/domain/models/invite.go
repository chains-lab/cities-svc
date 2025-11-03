package models

import (
	"time"

	"github.com/google/uuid"
)

type Invite struct {
	ID        uuid.UUID
	CityID    uuid.UUID
	UserID    uuid.UUID
	Status    string
	Role      string
	ExpiresAt time.Time
	CreatedAt time.Time
}

func (i Invite) IsNil() bool {
	return i.ID == uuid.Nil
}
