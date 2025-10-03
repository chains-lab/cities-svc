package enum

import (
	"fmt"
)

const (
	CityAdminRoleHead      = "head"
	CityAdminMember        = "member"
	CityAdminRoleModerator = "moderator"
)

var citiesAdminsRoles = []string{
	CityAdminRoleHead,
	CityAdminMember,
	CityAdminRoleModerator,
}

var ErrorInvalidCityAdminRole = fmt.Errorf("invalid city admins role mus be one of: %s", GetAllCityAdminRoles())

func CheckCityAdminRole(role string) error {
	for _, r := range citiesAdminsRoles {
		if r == role {
			return nil
		}
	}

	return fmt.Errorf("'%s', %w", role, ErrorInvalidCityAdminRole)
}

func GetAllCityAdminRoles() []string {
	return citiesAdminsRoles
}
