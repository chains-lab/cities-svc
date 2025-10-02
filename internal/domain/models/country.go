package models

import (
	"time"

	"github.com/google/uuid"
)

type Country struct {
	ID        uuid.UUID
	Name      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c Country) IsNil() bool {
	return c.ID == uuid.Nil
}

type CountriesCollection struct {
	Data  []Country `json:"data"`
	Page  uint64    `json:"page"`
	Size  uint64    `json:"size"`
	Total uint64    `json:"total"`
}
