package responses

import (
	cityAdminProto "github.com/chains-lab/cities-dir-proto/gen/go/citygov"
	"github.com/chains-lab/cities-dir-svc/internal/app/models"
)

func CitiesAdminsList(cityAdmins []models.CityAdmin) *cityAdminProto.ListCitiesAdmins {
	cityAdminsList := make([]*cityAdminProto.CityAdmin, len(cityAdmins))
	for i, cityAdmin := range cityAdmins {
		cityAdminsList[i] = CityAdmin(cityAdmin)
	}

	return &cityAdminProto.ListCitiesAdmins{
		Admins: cityAdminsList,
		// TODO: add pagination fields if needed
	}
}
