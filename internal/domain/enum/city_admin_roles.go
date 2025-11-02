package enum

import (
	"fmt"
)

const (
	CityAdminChief         = "chief"
	CityAdminRoleExecutive = "government"
	CityAdminModerator     = "moderator"
)

var citiesAdminRoles = []string{
	CityAdminChief,
	CityAdminRoleExecutive,
	CityAdminModerator,
}

var ErrorInvalidCityAdminRole = fmt.Errorf("invalid city admin role must be one of: %s", GetAllCityAdminRoles())

func CheckCityAdminRole(role string) error {
	for _, r := range citiesAdminRoles {
		if r == role {
			return nil
		}
	}

	return fmt.Errorf("'%s', %w", role, ErrorInvalidCityAdminRole)
}

func GetAllCityAdminRoles() []string {
	return citiesAdminRoles
}
