package constant

import (
	"fmt"
)

const (
	CityGovRoleMayor      = "mayor"
	CityGovRoleGovernment = "government"
	CityGovRoleModerator  = "moderator"
)

var citiesAdminsRoles = []string{
	CityGovRoleMayor,
	CityGovRoleGovernment,
	CityGovRoleModerator,
}

var ErrorInvalidCityGovRole = fmt.Errorf("invalid city government role mus be one of: %s", GetAllCitiesAdminsRoles())

func ParseCityGovRole(role string) error {
	for _, r := range citiesAdminsRoles {
		if r == role {
			return nil
		}
	}

	return fmt.Errorf("'%s', %w", role, ErrorInvalidCityGovRole)
}

func CompareCityGovRole(role1, role2 string) int {
	power := map[string]uint8{
		CityGovRoleMayor:      3,
		CityGovRoleGovernment: 2,
		CityGovRoleModerator:  1,
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
