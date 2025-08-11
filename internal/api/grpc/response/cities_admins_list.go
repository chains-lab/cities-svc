package response

import (
	cityAdminProto "github.com/chains-lab/cities-dir-proto/gen/go/citygov"
	pagProto "github.com/chains-lab/cities-dir-proto/gen/go/common/pagination"
	"github.com/chains-lab/cities-dir-svc/internal/app/models"
	"github.com/chains-lab/cities-dir-svc/internal/pagination"
)

func CitiesAdminsList(cityAdmins []models.CityAdmin, response pagination.Response) *cityAdminProto.ListCitiesAdmins {
	cityAdminsList := make([]*cityAdminProto.CityAdmin, len(cityAdmins))
	for i, cityAdmin := range cityAdmins {
		cityAdminsList[i] = CityAdmin(cityAdmin)
	}

	return &cityAdminProto.ListCitiesAdmins{
		Admins: cityAdminsList,
		Pagination: &pagProto.Response{
			Page:  response.Page,
			Size:  response.Size,
			Total: response.Total,
		},
	}
}
