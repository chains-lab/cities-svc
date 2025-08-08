package responses

import (
	"github.com/chains-lab/cities-dir-proto/gen/go/common/pagination"
	countryProto "github.com/chains-lab/cities-dir-proto/gen/go/country"
	"github.com/chains-lab/cities-dir-svc/internal/app/models"
)

func CountriesList(arr []models.Country) *countryProto.CountriesList {
	countryList := make([]*countryProto.Country, len(arr))
	for i, country := range arr {
		countryList[i] = Country(country)
	}

	return &countryProto.CountriesList{
		Countries: countryList,
		Pagination: &pagination.Response{
			Page:  1, //TODO
			Limit: 1, //TODO
		},
	}
}
