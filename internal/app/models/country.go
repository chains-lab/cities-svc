package models

import (
	"time"

	"github.com/google/uuid"
)

type CountryModel struct {
	ID        uuid.UUID
	Name      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
