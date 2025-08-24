package enum

import (
	"fmt"
)

const (
	CityGovRoleAdmin     = "admin"
	CityGovRoleModerator = "moderator"
)

var citiesAdminsRoles = []string{
	CityGovRoleAdmin,
	CityGovRoleModerator,
}

var ErrorInvalidCityGovRole = fmt.Errorf("invalid city government role mus be one of: %s", GetAllCitiesAdminsRoles())

func ParseCityGovRole(role string) (string, error) {
	for _, r := range citiesAdminsRoles {
		if r == role {
			return r, nil
		}
	}

	return "", fmt.Errorf("'%s', %w", role, ErrorInvalidCityGovRole)
}

func CompareCityGovRole(role1, role2 string) int {
	power := map[string]uint8{
		CityGovRoleAdmin:     2,
		CityGovRoleModerator: 1,
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
