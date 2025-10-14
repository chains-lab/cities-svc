package responses

import (
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/cities-svc/resources"
)

func City(m models.City) resources.City {
	resp := resources.City{
		Data: resources.CityData{
			Id:   m.ID,
			Type: resources.CityType,
			Attributes: resources.CityAttributes{
				CountryId: m.CountryID,
				Status:    m.Status,
				Name:      m.Name,
				Point: resources.Point{
					Latitude:  m.Point[0],
					Longitude: m.Point[1],
				},
				Timezone:  m.Timezone,
				CreatedAt: m.CreatedAt,
				UpdatedAt: m.UpdatedAt,
			},
		},
	}

	if m.Icon != nil {
		resp.Data.Attributes.Icon = m.Icon
	}
	if m.Slug != nil {
		resp.Data.Attributes.Slug = m.Slug
	}

	return resp
}

func CitiesCollection(ms models.CitiesCollection) resources.CitiesCollection {
	resp := resources.CitiesCollection{
		Data: make([]resources.CityData, 0, len(ms.Data)),
		Links: resources.PaginationData{
			PageNumber: int64(ms.Page),
			PageSize:   int64(ms.Size),
			TotalItems: int64(ms.Total),
		},
	}

	for _, m := range ms.Data {
		city := City(m).Data

		resp.Data = append(resp.Data, city)
	}

	return resp
}
