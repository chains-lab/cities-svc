package response

import (
	pagProto "github.com/chains-lab/cities-dir-proto/gen/go/common/pagination"
	cityProto "github.com/chains-lab/cities-dir-proto/gen/go/svc/city"
	"github.com/chains-lab/cities-dir-svc/internal/app/models"
	"github.com/chains-lab/cities-dir-svc/internal/pagination"
)

func CitiesList(cities []models.City, pag pagination.Response) *cityProto.CitiesList {
	cityList := make([]*cityProto.City, len(cities))
	for i, city := range cities {
		cityList[i] = City(city)
	}

	return &cityProto.CitiesList{
		Cities: cityList,
		Pagination: &pagProto.Response{
			Page:  pag.Page,
			Size:  pag.Size,
			Total: pag.Total,
		},
	}
}
