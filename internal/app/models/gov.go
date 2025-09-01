package models

import (
	"time"

	"github.com/google/uuid"
)

type CityGov struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	CityID    uuid.UUID
	Active    bool
	Role      string
	Label     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
