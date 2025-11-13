package models

import (
	"time"

	"github.com/google/uuid"
)

type Invite struct {
	ID          uuid.UUID `json:"id"`
	CityID      uuid.UUID `json:"city_id"`
	UserID      uuid.UUID `json:"user_id"`
	InitiatorID uuid.UUID `json:"initiator_id"`
	Status      string    `json:"status"`
	Role        string    `json:"role"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
}

func (i Invite) IsNil() bool {
	return i.ID == uuid.Nil
}
