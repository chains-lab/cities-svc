package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type City struct {
	ID        uuid.UUID
	CountryID uuid.UUID
	Point     orb.Point // [lon, lat]
	Status    string
	Name      string
	Icon      *string
	Slug      *string
	Timezone  string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c City) IsNil() bool {
	return c.ID == uuid.Nil
}

type CitiesCollection struct {
	Data  []City `json:"data"`
	Page  uint64 `json:"page"`
	Size  uint64 `json:"size"`
	Total uint64 `json:"total"`
}
