package response

import (
	cityAdminProto "github.com/chains-lab/cities-dir-proto/gen/go/citygov"
	"github.com/chains-lab/cities-dir-svc/internal/app/models"
)

func CityAdmin(cityAdmin models.CityGov) *cityAdminProto.CityGov {
	return &cityAdminProto.CityGov{
		CityId: cityAdmin.CityID.String(),
		UserId: cityAdmin.UserID.String(),
		Role:   cityAdmin.Role,
	}
}
