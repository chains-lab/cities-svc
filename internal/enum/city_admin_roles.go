package enum

import "math"

type CityAdminRole string

const (
	CityOwner     CityAdminRole = "owner"
	CityAdmin     CityAdminRole = "admin"
	CityModerator CityAdminRole = "moderator"
)

var citiesAdminsRoles = []CityAdminRole{
	CityOwner,
	CityAdmin,
	CityModerator,
}

func ParseCityAdminRole(role string) (CityAdminRole, bool) {
	for _, r := range citiesAdminsRoles {
		if r == CityAdminRole(role) {
			return r, true
		}
	}
	return "", false
}

func CompareCityAdminRole(role1, role2 CityAdminRole) int {
	power := map[CityAdminRole]uint8{
		CityOwner:     math.MaxUint8,
		CityAdmin:     2,
		CityModerator: 1,
	}

	if power[role1] > power[role2] {
		return 1
	} else if power[role1] < power[role2] {
		return -1
	}

	return 0
}

func GetAllCitiesAdminsRoles() []CityAdminRole {
	return citiesAdminsRoles
}
