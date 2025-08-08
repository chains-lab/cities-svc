package models

import (
	"time"

	"github.com/google/uuid"
)

type CityAdmin struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	CityID    uuid.UUID
	Role      string
	UpdatedAt time.Time
	CreatedAt time.Time
}
