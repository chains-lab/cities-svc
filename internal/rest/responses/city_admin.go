package responses

import (
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/cities-svc/resources"
)

func CityAdmin(m models.CityAdmin) resources.CityAdmin {
	return resources.CityAdmin{
		Data: resources.CityAdminData{
			Id:   m.Data.UserID,
			Type: resources.GovType,
			Attributes: resources.CityAdminAttributes{
				CityId:    m.Data.CityID,
				Label:     m.Data.Label,
				Username:  m.Username,
				Avatar:    m.Avatar,
				Role:      m.Data.Role,
				CreatedAt: m.Data.CreatedAt,
				UpdatedAt: m.Data.UpdatedAt,
			},
		},
	}
}

func CityAdminsCollection(ms models.CityAdminCollection) resources.CityAdminsCollection {
	resp := resources.CityAdminsCollection{
		Data: make([]resources.CityAdminData, 0, len(ms.Data)),
		Links: resources.PaginationData{
			PageNumber: int64(ms.Page),
			PageSize:   int64(ms.Size),
			TotalItems: int64(ms.Total),
		},
	}

	for _, m := range ms.Data {
		gov := CityAdmin(m).Data

		resp.Data = append(resp.Data, gov)
	}

	return resp
}
