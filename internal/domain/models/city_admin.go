package models

import (
	"time"

	"github.com/google/uuid"
)

type CityAdmin struct {
	UserID    uuid.UUID
	CityID    uuid.UUID
	Role      string
	Label     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (a CityAdmin) IsNil() bool {
	return a.UserID == uuid.Nil
}

type CityAdminsCollection struct {
	Data  []CityAdmin `json:"data"`
	Page  uint64      `json:"page"`
	Size  uint64      `json:"size"`
	Total uint64      `json:"total"`
}
