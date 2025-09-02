package models

import (
	"time"

	"github.com/google/uuid"
)

type Invite struct {
	ID          uuid.UUID
	Status      string
	Role        string
	CityID      uuid.UUID
	InitiatorID uuid.UUID
	UserID      *uuid.UUID
	AnsweredAt  *time.Time
	ExpiresAt   time.Time
	CreatedAt   time.Time
}
