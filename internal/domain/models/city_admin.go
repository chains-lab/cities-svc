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
	Position  *string
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
	Data     CityAdmin
	Username string  `json:"username"`
	Avatar   *string `json:"avatar"`
}

func (a CityAdminWithUserData) IsNil() bool {
	return a.Data.UserID == uuid.Nil
}

func (a CityAdmin) AddProfileData(profile Profile) CityAdminWithUserData {
	return CityAdminWithUserData{
		Data:     a,
		Username: profile.Username,
		Avatar:   profile.Avatar,
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
