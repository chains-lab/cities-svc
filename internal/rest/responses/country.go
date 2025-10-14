package responses

import (
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/cities-svc/resources"
)

func Country(m models.Country) resources.Country {
	resp := resources.Country{
		Data: resources.CountryData{
			Id:   m.ID,
			Type: resources.CountryType,
			Attributes: resources.CountryAttributes{
				Name:      m.Name,
				Status:    m.Status,
				CreatedAt: m.CreatedAt,
				UpdatedAt: m.UpdatedAt,
			},
		},
	}

	return resp
}

func CountriesCollection(ms models.CountriesCollection) resources.CountriesCollection {
	resp := resources.CountriesCollection{
		Data: make([]resources.CountryData, 0, len(ms.Data)),
		Links: resources.PaginationData{
			PageNumber: int64(ms.Page),
			PageSize:   int64(ms.Size),
			TotalItems: int64(ms.Total),
		},
	}

	for _, m := range ms.Data {
		country := Country(m).Data

		resp.Data = append(resp.Data, country)
	}

	return resp
}
