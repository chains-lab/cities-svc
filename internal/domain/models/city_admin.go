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

type CityAdminWithUserData struct {
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	Avatar    *string   `json:"avatar"`
	CityID    uuid.UUID `json:"city_id"`
	Role      string    `json:"role"`
	Label     *string   `json:"label"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (a CityAdminWithUserData) IsNil() bool {
	return a.UserID == uuid.Nil
}

func (a CityAdmin) AddProfileData(profile Profile) CityAdminWithUserData {
	return CityAdminWithUserData{
		UserID:    a.UserID,
		Username:  profile.Username,
		Avatar:    profile.Avatar,
		CityID:    a.CityID,
		Role:      a.Role,
		Label:     a.Label,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

type CityAdminsWithUserDataCollection struct {
	Data  []CityAdminWithUserData `json:"data"`
	Page  uint64                  `json:"page"`
	Size  uint64                  `json:"size"`
	Total uint64                  `json:"total"`
}

func (c CityAdminsCollection) AddProfileData(profiles map[uuid.UUID]Profile) CityAdminsWithUserDataCollection {
	employees := make([]CityAdminWithUserData, 0, len(c.Data))
	for _, emp := range c.Data {
		empWithProfile := emp.AddProfileData(profiles[emp.UserID])
		employees = append(employees, empWithProfile)
	}
	return CityAdminsWithUserDataCollection{
		Data:  employees,
		Page:  c.Page,
		Size:  c.Size,
		Total: c.Total,
	}
}
