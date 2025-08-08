package renderer

import (
	svccitiesadmins "github.com/chains-lab/cities-dir-proto/gen/go/citiesadmins"
	"github.com/chains-lab/cities-dir-svc/internal/app/models"
)

func CityAdmin(cityAdmin models.CityAdmin) *svccitiesadmins.CityAdmin {
	return &svccitiesadmins.CityAdmin{
		CityId: cityAdmin.CityID.String(),
		UserId: cityAdmin.UserID.String(),
		Role:   string(cityAdmin.Role),
	}
}
