package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type City struct {
	ID        uuid.UUID
	CountryID uuid.UUID
	Status    string
	Zone      orb.MultiPolygon
	Name      string
	Icon      string
	Slug      string
	Timezone  string

	CreatedAt time.Time
	UpdatedAt time.Time
}
