package response

import (
	pagProto "github.com/chains-lab/cities-dir-proto/gen/go/common/pagination"
	countryProto "github.com/chains-lab/cities-dir-proto/gen/go/country"
	"github.com/chains-lab/cities-dir-svc/internal/app/models"
	"github.com/chains-lab/cities-dir-svc/internal/pagination"
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
