package enum

import (
	"fmt"
)

const (
	//This role is about city gov administration
	CityAdminRoleChief     = "chief"
	CityAdminRoleViceChief = "vice-chief"
	CityAdminRoleMember    = "member"

	//This role is about city tech administration
	CityAdminRoleTechLead  = "tech-lead"
	CityAdminRoleModerator = "moderator"
)

var citiesAdminRoles = []string{
	CityAdminRoleChief,
	CityAdminRoleViceChief,
	CityAdminRoleMember,

	CityAdminRoleTechLead,
	CityAdminRoleModerator,
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

func RightCityAdminsTechPolitics(first, second string) bool {
	switch first {
	case CityAdminRoleModerator:
		if second != CityAdminRoleTechLead && second != CityAdminRoleChief {
			return true
		}
		return false
	case CityAdminRoleTechLead:
		return true
	default:
		return false
	}
}
