package responses

import (
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/resources"
)

func Invite(m models.Invite) resources.Invite {
	resp := resources.Invite{
		Data: resources.InviteData{
			Id:   m.ID.String(),
			Type: resources.InviteType,
			Attributes: resources.InviteDataAttributes{
				Status:      m.Status,
				Role:        m.Role,
				CityId:      m.CityID.String(),
				InitiatorId: m.InitiatorID.String(),
				ExpiresAt:   m.ExpiresAt,
				CreatedAt:   m.CreatedAt,
			},
		},
	}

	if m.UserID != nil {
		userIDStr := m.UserID.String()
		resp.Data.Attributes.UserId = &userIDStr
	}
	if m.AnsweredAt != nil {
		resp.Data.Attributes.AnsweredAt = m.AnsweredAt
	}

	return resp
}
