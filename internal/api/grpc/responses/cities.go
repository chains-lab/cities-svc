package responses

import (
	cityProto "github.com/chains-lab/cities-proto/gen/go/svc/city"
	"github.com/chains-lab/cities-svc/internal/app/models"
)

func City(city models.City) *cityProto.City {
	return &cityProto.City{
		Id:        city.ID.String(),
		Name:      city.Name,
		Status:    city.Status,
		CountryId: city.CountryID.String(),
	}
}
