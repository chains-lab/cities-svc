package responses

import (
	cityAdminProto "github.com/chains-lab/cities-dir-proto/gen/go/citygov"
	"github.com/chains-lab/cities-dir-svc/internal/app/models"
)

func CityAdmin(cityAdmin models.CityAdmin) *cityAdminProto.CityAdmin {
	return &cityAdminProto.CityAdmin{
		CityId: cityAdmin.CityID.String(),
		UserId: cityAdmin.UserID.String(),
		Role:   string(cityAdmin.Role),
	}
}
