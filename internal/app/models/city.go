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
	Center    orb.Point        // [lon, lat]
	Boundary  orb.MultiPolygon // многоугольники границы
	Icon      string
	Slug      string
	Timezone  string

	Details CityDetail

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CityDetail struct {
	Name        string
	Description *string
	Language    string
}
