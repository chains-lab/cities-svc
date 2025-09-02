package responses

import (
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/resources"
	"github.com/chains-lab/pagi"
)

func City(m models.City) resources.City {
	resp := resources.City{
		Data: resources.CityData{
			Id:   m.ID.String(),
			Type: resources.CityType,
			Attributes: resources.CityAttributes{
				CountryId: m.CountryID.String(),
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

func CitiesCollection(ms []models.City, pag pagi.Response) resources.CitiesCollection {
	resp := resources.CitiesCollection{
		Data: make([]resources.CityData, 0, len(ms)),
		Links: resources.PaginationData{
			PageNumber: int64(pag.Page),
			PageSize:   int64(pag.Size),
			TotalItems: int64(pag.Total),
		},
	}

	for _, m := range ms {
		city := City(m).Data

		resp.Data = append(resp.Data, city)
	}

	return resp
}
