package responses

import (
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/cities-svc/resources"
)

func CityAdmin(m models.CityAdmin) resources.CityAdmin {
	resp := resources.CityAdmin{
		Data: resources.CityAdminData{
			Id:   m.UserID,
			Type: resources.GovType,
			Attributes: resources.CityAdminAttributes{
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

func CityAdminsCollection(ms models.CityAdminsCollection) resources.CityAdminsCollection {
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
