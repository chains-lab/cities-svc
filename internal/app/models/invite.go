package models

import (
	"time"

	"github.com/google/uuid"
)

type Invite struct {
	ID         uuid.UUID
	Status     string
	Role       string
	CityID     uuid.UUID
	UserID     *uuid.UUID
	AnsweredAt *time.Time
	Token      string
	ExpiresAt  time.Time
	CreatedAt  time.Time
}
