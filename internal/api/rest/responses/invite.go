package responses

import (
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/cities-svc/resources"
)

func Invite(m models.Invite, token models.InviteToken) resources.Invite {
	resp := resources.Invite{
		Data: resources.InviteData{
			Id:   m.ID,
			Type: resources.InviteType,
			Attributes: resources.InviteDataAttributes{
				Status:    m.Status,
				Role:      m.Role,
				Token:     string(token),
				CityId:    m.CityID,
				ExpiresAt: m.ExpiresAt,
				CreatedAt: m.CreatedAt,
			},
		},
	}

	return resp
}
