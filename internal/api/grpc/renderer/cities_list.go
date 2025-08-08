package renderer

import (
	svc "github.com/chains-lab/cities-dir-proto/gen/go/cities"
	"github.com/chains-lab/cities-dir-proto/gen/go/common/pagination"
	"github.com/chains-lab/cities-dir-svc/internal/app/models"
)

func CitiesList(cities []models.City) *svc.CitiesList {
	cityList := make([]*svc.City, len(cities))
	for i, city := range cities {
		cityList[i] = City(city)
	}

	return &svc.CitiesList{
		Cities: cityList,
		Pagination: &pagination.Response{
			Page:  1, //TODO
			Limit: 1, //TODO
		},
	}
}
