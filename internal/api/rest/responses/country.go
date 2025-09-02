package responses

import (
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/resources"
	"github.com/chains-lab/pagi"
)

func Country(m models.Country) resources.Country {
	resp := resources.Country{
		Data: resources.CountryData{
			Id:   m.ID.String(),
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

func CountriesCollection(ms []models.Country, pag pagi.Response) resources.CountriesCollection {
	resp := resources.CountriesCollection{
		Data: make([]resources.CountryData, 0, len(ms)),
		Links: resources.PaginationData{
			PageNumber: int64(pag.Page),
			PageSize:   int64(pag.Size),
			TotalItems: int64(pag.Total),
		},
	}

	for _, m := range ms {
		country := Country(m).Data

		resp.Data = append(resp.Data, country)
	}

	return resp
}
