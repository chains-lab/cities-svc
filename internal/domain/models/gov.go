package models

import (
	"time"

	"github.com/google/uuid"
)

type CityModer struct {
	UserID    uuid.UUID
	CityID    uuid.UUID
	Role      string
	Label     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m CityModer) IsNil() bool {
	return m.UserID == uuid.Nil
}

type CityModersCollection struct {
	Data  []CityModer `json:"data"`
	Page  uint64      `json:"page"`
	Size  uint64      `json:"size"`
	Total uint64      `json:"total"`
}
