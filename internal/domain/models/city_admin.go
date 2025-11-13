package models

import (
	"time"

	"github.com/google/uuid"
)

type CityAdmin struct {
	UserID    uuid.UUID `json:"user_id"`
	CityID    uuid.UUID `json:"city_id"`
	Role      string    `json:"role"`
	Label     *string   `json:"label,omitempty"`
	Position  *string   `json:"position,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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

func (c CityAdminsCollection) GetUserIDs() []uuid.UUID {
	ids := make([]uuid.UUID, 0, len(c.Data))
	for _, admin := range c.Data {
		ids = append(ids, admin.UserID)
	}
	return ids
}
