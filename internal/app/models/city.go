package models

import (
	"time"

	"github.com/google/uuid"
)

type City struct {
	ID        uuid.UUID
	CountryID uuid.UUID
	Name      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
