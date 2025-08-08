package renderer

import (
	svccitiesadmins "github.com/chains-lab/cities-dir-proto/gen/go/citiesadmins"
	"github.com/chains-lab/cities-dir-svc/internal/app/models"
)

func CitiesAdminsList(cityAdmins []models.CityAdmin) *svccitiesadmins.ListCitiesAdmins {
	cityAdminsList := make([]*svccitiesadmins.CityAdmin, len(cityAdmins))
	for i, cityAdmin := range cityAdmins {
		cityAdminsList[i] = CityAdmin(cityAdmin)
	}

	return &svccitiesadmins.ListCitiesAdmins{
		Admins: cityAdminsList,
		// TODO: add pagination fields if needed
	}
}
