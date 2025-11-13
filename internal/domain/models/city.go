package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type City struct {
	ID        uuid.UUID `json:"id"`
	CountryID string    `json:"country_id"`
	Point     orb.Point `json:"point"` // [lon, lat]
	Status    string    `json:"status"`
	Name      string    `json:"name"`
	Icon      *string   `json:"icon,omitempty"`
	Slug      *string   `json:"slug,omitempty"`
	Timezone  string    `json:"timezone"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
