package enum

import (
	"fmt"
)

const (
	CityAdminRoleAdmin     = "admin"
	CityAdminRoleModerator = "moderator"
)

var citiesAdminsRoles = []string{
	CityAdminRoleAdmin,
	CityAdminRoleModerator,
}

var ErrorInvalidCityAdminRole = fmt.Errorf("invalid city admin role mus be one of: %s", GetAllCitiesAdminsRoles())

func ParseCityAdminRole(role string) (string, error) {
	for _, r := range citiesAdminsRoles {
		if r == role {
			return r, nil
		}
	}

	return "", fmt.Errorf("'%s', %w", role, ErrorInvalidCityAdminRole)
}

func CompareCityAdminRole(role1, role2 string) int {
	power := map[string]uint8{
		CityAdminRoleAdmin:     2,
		CityAdminRoleModerator: 1,
	}

	if power[role1] > power[role2] {
		return 1
	} else if power[role1] < power[role2] {
		return -1
	}

	return 0
}

func GetAllCitiesAdminsRoles() []string {
	return citiesAdminsRoles
}
