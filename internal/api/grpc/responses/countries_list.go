package responses

import (
	pagProto "github.com/chains-lab/cities-proto/gen/go/common/pagination"
	countryProto "github.com/chains-lab/cities-proto/gen/go/svc/country"
	"github.com/chains-lab/cities-svc/internal/app/models"
)

func CountriesList(arr []models.Country, response pagination.Response) *countryProto.CountriesList {
	countryList := make([]*countryProto.Country, len(arr))
	for i, country := range arr {
		countryList[i] = Country(country)
	}

	return &countryProto.CountriesList{
		Countries: countryList,
		Pagination: &pagProto.Response{
			Page:  response.Page,
			Size:  response.Size,
			Total: response.Total,
		},
	}
}
