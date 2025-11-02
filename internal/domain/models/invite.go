package models

import (
	"time"

	"github.com/google/uuid"
)

type Invite struct {
	ID        uuid.UUID
	CityID    uuid.UUID
	UserID    uuid.UUID
	Target    string
	Role      string
	Status    string
	ExpiresAt time.Time
	CreatedAt time.Time
}

func (i Invite) IsNil() bool {
	return i.ID == uuid.Nil
}

type InviteWithUserData struct {
	Username string  `json:"username"`
	Avatar   *string `json:"avatar"`
	Invite   Invite  `json:"invite"`
}

func (i InviteWithUserData) IsNil() bool {
	return i.Invite.ID == uuid.Nil
}

func (i Invite) AddProfileData(profile Profile) InviteWithUserData {
	return InviteWithUserData{
		Username: profile.Username,
		Avatar:   profile.Avatar,
		Invite: Invite{
			ID:        i.ID,
			CityID:    i.CityID,
			UserID:    i.UserID,
			Target:    i.Target,
			Role:      i.Role,
			Status:    i.Status,
			ExpiresAt: i.ExpiresAt,
			CreatedAt: i.CreatedAt,
		},
	}
}

type InvitesWithUserDataCollection struct {
	Data  []InviteWithUserData `json:"data"`
	Page  uint64               `json:"page"`
	Size  uint64               `json:"size"`
	Total uint64               `json:"total"`
}

func (c InvitesWithUserDataCollection) AddProfileData(profiles map[uuid.UUID]Profile) InvitesWithUserDataCollection {
	invites := make([]InviteWithUserData, 0, len(c.Data))
	for _, invite := range c.Data {
		profile, exists := profiles[invite.Invite.UserID]
		if !exists {
			continue
		}

		invites = append(invites, invite.Invite.AddProfileData(profile))
	}

	return InvitesWithUserDataCollection{
		Data:  invites,
		Page:  c.Page,
		Size:  c.Size,
		Total: c.Total,
	}
}
