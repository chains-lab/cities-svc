package models

import (
	"time"

	"github.com/google/uuid"
)

type CityAdmin struct {
	UserID    uuid.UUID
	CityID    uuid.UUID
	Role      string
	Position  *string
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

type CityAdminsWithUserData struct {
	Username string    `json:"username"`
	Avatar   *string   `json:"avatar"`
	Admin    CityAdmin `json:"admin"`
}

func (a CityAdminsWithUserData) IsNil() bool {
	return a.Admin.UserID == uuid.Nil
}

func (a CityAdmin) AddProfileData(profile Profile) CityAdminsWithUserData {
	return CityAdminsWithUserData{
		Username: profile.Username,
		Avatar:   profile.Avatar,
		Admin: CityAdmin{
			UserID:    a.UserID,
			CityID:    a.CityID,
			Role:      a.Role,
			Label:     a.Label,
			Position:  a.Position,
			CreatedAt: a.CreatedAt,
			UpdatedAt: a.UpdatedAt,
		},
	}
}

type CityAdminsWithUserDataCollection struct {
	Data  []CityAdminsWithUserData `json:"data"`
	Page  uint64                   `json:"page"`
	Size  uint64                   `json:"size"`
	Total uint64                   `json:"total"`
}

func (c CityAdminsCollection) AddProfileData(profiles map[uuid.UUID]Profile) CityAdminsWithUserDataCollection {
	employees := make([]CityAdminsWithUserData, 0, len(c.Data))
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
