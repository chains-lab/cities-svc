package responses

import (
	cityProto "github.com/chains-lab/cities-dir-proto/gen/go/city"
	"github.com/chains-lab/cities-dir-proto/gen/go/common/pagination"
	"github.com/chains-lab/cities-dir-svc/internal/app/models"
)

func CitiesList(cities []models.City) *cityProto.CitiesList {
	cityList := make([]*cityProto.City, len(cities))
	for i, city := range cities {
		cityList[i] = City(city)
	}

	return &cityProto.CitiesList{
		Cities: cityList,
		Pagination: &pagination.Response{
			Page:  1, //TODO
			Limit: 1, //TODO
		},
	}
}
