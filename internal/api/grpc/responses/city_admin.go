package responses

import (
	cityAdminProto "github.com/chains-lab/cities-dir-proto/gen/go/svc/gov"
	"github.com/chains-lab/cities-dir-svc/internal/app/models"
)

func CityAdmin(cityAdmin models.CityGov) *cityAdminProto.CityGov {
	return &cityAdminProto.CityGov{
		CityId: cityAdmin.CityID.String(),
		UserId: cityAdmin.UserID.String(),
		Role:   cityAdmin.Role,
	}
}
