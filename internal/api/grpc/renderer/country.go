package renderer

import (
	"github.com/chains-lab/cities-dir-proto/gen/go/countries"
	"github.com/chains-lab/cities-dir-svc/internal/app/models"
)

func Country(country models.Country) *countries.Country {
	return &countries.Country{
		Id:     country.ID.String(),
		Name:   country.Name,
		Status: string(country.Status),
	}
}
