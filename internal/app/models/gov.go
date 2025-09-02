package models

import (
	"time"

	"github.com/google/uuid"
)

type Gov struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	CityID        uuid.UUID
	Status        string
	Role          string
	Label         string
	DeactivatedAt *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
