package responses

import (
	cityProto "github.com/chains-lab/cities-dir-proto/gen/go/city"
	"github.com/chains-lab/cities-dir-svc/internal/app/models"
)

func City(city models.City) *cityProto.City {
	return &cityProto.City{
		Id:        city.ID.String(),
		Name:      city.Name,
		Status:    string(city.Status),
		CountryId: city.CountryID.String(),
	}
}
