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
	Token      string
	UserID     *uuid.UUID
	AnsweredAt *time.Time
	ExpiresAt  time.Time
	CreatedAt  time.Time
}

func (i Invite) IsNil() bool {
	return i == Invite{}
}

type InviteTokenData struct {
	InviteID  uuid.UUID
	CityID    uuid.UUID
	Role      string
	ExpiresAt time.Time
}

func (i InviteTokenData) IsNil() bool {
	return i == InviteTokenData{}
}
