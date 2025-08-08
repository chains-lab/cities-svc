package renderer

import (
	"github.com/chains-lab/cities-dir-proto/gen/go/cities"
	"github.com/chains-lab/cities-dir-svc/internal/app/models"
)

func City(city models.City) *cities.City {
	return &cities.City{
		Id:        city.ID.String(),
		Name:      city.Name,
		Status:    string(city.Status),
		CountryId: city.CountryID.String(),
	}
}
