package responses

import (
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/cities-svc/resources"
)

func Gov(m models.CityModer) resources.Gov {
	resp := resources.Gov{
		Data: resources.GovData{
			Id:   m.UserID,
			Type: resources.GovType,
			Attributes: resources.GovAttributes{
				CityId:    m.CityID,
				Role:      m.Role,
				CreatedAt: m.CreatedAt,
				UpdatedAt: m.UpdatedAt,
			},
		},
	}
	if m.Label != nil {
		resp.Data.Attributes.Label = m.Label
	}

	return resp
}

func GovsCollection(ms models.CityModersCollection) resources.GovsCollection {
	resp := resources.GovsCollection{
		Data: make([]resources.GovData, 0, len(ms.Data)),
		Links: resources.PaginationData{
			PageNumber: int64(ms.Page),
			PageSize:   int64(ms.Size),
			TotalItems: int64(ms.Total),
		},
	}

	for _, m := range ms.Data {
		gov := Gov(m).Data

		resp.Data = append(resp.Data, gov)
	}

	return resp
}
