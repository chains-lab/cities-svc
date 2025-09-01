package constant

import (
	"fmt"
)

const (
	CityGovRoleMayor     = "mayor"
	CityGovRoleAdvisor   = "advisor"
	CityGovRoleMember    = "member"
	CityGovRoleModerator = "moderator"
)

var citiesAdminsRoles = []string{
	CityGovRoleMayor,
	CityGovRoleAdvisor,
	CityGovRoleMember,
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

func GetAllCitiesAdminsRoles() []string {
	return citiesAdminsRoles
}
