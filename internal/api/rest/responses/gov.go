package responses

import (
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/resources"
	"github.com/chains-lab/pagi"
)

func Gov(m models.Gov) resources.Gov {
	resp := resources.Gov{
		Data: resources.GovData{
			Id:   m.UserID.String(),
			Type: resources.GovType,
			Attributes: resources.GovAttributes{
				CityId:    m.CityID.String(),
				Role:      m.Role,
				CreatedAt: m.CreatedAt,
				UpdatedAt: m.UpdatedAt,
			},
		},
	}
	if m.Label != nil {
		resp.Data.Attributes.Label = *m.Label
	}

	return resp
}

func GovsCollection(ms []models.Gov, pag pagi.Response) resources.GovsCollection {
	resp := resources.GovsCollection{
		Data: make([]resources.GovData, 0, len(ms)),
		Links: resources.PaginationData{
			PageNumber: int64(pag.Page),
			PageSize:   int64(pag.Size),
			TotalItems: int64(pag.Total),
		},
	}

	for _, m := range ms {
		gov := Gov(m).Data

		resp.Data = append(resp.Data, gov)
	}

	return resp
}
