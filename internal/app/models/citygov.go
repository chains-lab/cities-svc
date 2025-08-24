package models

import (
	"time"

	"github.com/google/uuid"
)

type CityGov struct {
	UserID    uuid.UUID
	CityID    uuid.UUID
	Role      string
	UpdatedAt time.Time
	CreatedAt time.Time
}
