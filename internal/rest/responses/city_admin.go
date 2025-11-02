package responses

import (
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/cities-svc/resources"
)

func CityAdmin(m models.CityAdminsWithUserData) resources.CityAdmin {
	return resources.CityAdmin{
		Data: resources.CityAdminData{
			Id:   m.UserID,
			Type: resources.adminType,
			Attributes: resources.CityAdminAttributes{
				CityId:    m.CityID,
				Label:     m.Label,
				Username:  m.Username,
				Avatar:    m.Avatar,
				Role:      m.Role,
				CreatedAt: m.CreatedAt,
				UpdatedAt: m.UpdatedAt,
			},
		},
	}
}

func CityAdminsCollection(ms models.CityAdminsWithUserDataCollection) resources.CityAdminsCollection {
	resp := resources.CityAdminsCollection{
		Data: make([]resources.CityAdminData, 0, len(ms.Data)),
		Links: resources.PaginationData{
			PageNumber: int64(ms.Page),
			PageSize:   int64(ms.Size),
			TotalItems: int64(ms.Total),
		},
	}

	for _, m := range ms.Data {
		admin := CityAdmin(m).Data

		resp.Data = append(resp.Data, admin)
	}

	return resp
}
