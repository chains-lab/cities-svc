package response

import (
	cityAdminProto "github.com/chains-lab/cities-dir-proto/gen/go/citygov"
	pagProto "github.com/chains-lab/cities-dir-proto/gen/go/common/pagination"
	"github.com/chains-lab/cities-dir-svc/internal/app/models"
	"github.com/chains-lab/cities-dir-svc/internal/pagination"
)

func CitiesAdminsList(cityAdmins []models.CityGov, response pagination.Response) *cityAdminProto.ListCityGovs {
	cityAdminsList := make([]*cityAdminProto.CityGov, len(cityAdmins))
	for i, cityAdmin := range cityAdmins {
		cityAdminsList[i] = CityAdmin(cityAdmin)
	}

	return &cityAdminProto.ListCityGovs{
		Government: cityAdminsList,
		Pagination: &pagProto.Response{
			Page:  response.Page,
			Size:  response.Size,
			Total: response.Total,
		},
	}
}
