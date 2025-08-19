package responses

import (
	countryProto "github.com/chains-lab/cities-proto/gen/go/svc/country"
	"github.com/chains-lab/cities-svc/internal/app/models"
)

func Country(country models.Country) *countryProto.Country {
	return &countryProto.Country{
		Id:     country.ID.String(),
		Name:   country.Name,
		Status: country.Status,
	}
}
