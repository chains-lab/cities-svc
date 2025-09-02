package models

import (
	"time"

	"github.com/google/uuid"
)

type Gov struct {
	UserID    uuid.UUID
	CityID    uuid.UUID
	Role      string
	Label     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
