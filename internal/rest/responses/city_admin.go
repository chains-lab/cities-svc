package responses

import (
	"fmt"

	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/cities-svc/resources"
)

func CityAdmin(m models.CityAdmin) resources.CityAdmin {
	return resources.CityAdmin{
		Data: resources.CityAdminData{
			Id:   fmt.Sprintf("%s:%s", m.UserID, m.CityID),
			Type: resources.CityAdminType,
			Attributes: resources.CityAdminAttributes{
				Label:     m.Label,
				Position:  m.Position,
				Role:      m.Role,
				CreatedAt: m.CreatedAt,
				UpdatedAt: m.UpdatedAt,
			},
		},
	}
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
