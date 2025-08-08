package renderer

import (
	"github.com/chains-lab/cities-dir-proto/gen/go/common/pagination"
	"github.com/chains-lab/cities-dir-proto/gen/go/countries"
	"github.com/chains-lab/cities-dir-svc/internal/app/models"
)

func CountriesList(arr []models.Country) *countries.CountriesList {
	countryList := make([]*countries.Country, len(arr))
	for i, country := range arr {
		countryList[i] = Country(country)
	}

	return &countries.CountriesList{
		Countries: countryList,
		Pagination: &pagination.Response{
			Page:  1, //TODO
			Limit: 1, //TODO
		},
	}
}
