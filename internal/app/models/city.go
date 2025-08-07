package models

import (
	"time"

	"github.com/chains-lab/cities-dir-svc/internal/enum"
	"github.com/google/uuid"
)

type City struct {
	ID        uuid.UUID
	CountryID uuid.UUID
	Name      string
	Status    enum.CityStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}
