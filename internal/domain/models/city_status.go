package models

import "time"

type CityStatus struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Accessible   bool      `json:"accessible"`
	AllowedAdmin bool      `json:"allowed_admins"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (c CityStatus) IsNil() bool {
	return c.ID == ""
}

type CityStatusesCollection struct {
	Data  []CityStatus `json:"data"`
	Page  uint64       `json:"page"`
	Size  uint64       `json:"size"`
	Total uint64       `json:"total"`
}
