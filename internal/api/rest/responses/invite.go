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
				Status:    m.Status,
				Role:      m.Role,
				CityId:    m.CityID.String(),
				Token:     m.Token,
				ExpiresAt: m.ExpiresAt,
				CreatedAt: m.CreatedAt,
			},
		},
	}

	return resp
}
