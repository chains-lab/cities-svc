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

var ErrorInvalidCityGovRole = fmt.Errorf("invalid city government role mus be one of: %s", GetAllCityGovsRoles())

func CheckCityGovRole(role string) error {
	for _, r := range citiesAdminsRoles {
		if r == role {
			return nil
		}
	}

	return fmt.Errorf("'%s', %w", role, ErrorInvalidCityGovRole)
}

func GetAllCityGovsRoles() []string {
	return citiesAdminsRoles
}

var citiesGovsPriority = map[string]int{
	CityGovRoleMayor:     4,
	CityGovRoleAdvisor:   3,
	CityGovRoleModerator: 2,
	CityGovRoleMember:    1,
}

// CompareCityGovRoles
// res : 1, if first role is higher priority
// res : -1, if second role is higher priority
// res : 0, if roles are equal
func CompareCityGovRoles(role1, role2 string) (int, error) {
	err := CheckCityGovRole(role1)
	if err != nil {
		return -1, err
	}

	err = CheckCityGovRole(role2)
	if err != nil {
		return -1, err
	}

	p1, ok1 := citiesGovsPriority[role1]
	p2, ok2 := citiesGovsPriority[role2]

	if !ok1 || !ok2 {
		return -1, nil
	}

	if p1 > p2 {
		return 1, nil
	} else if p1 < p2 {
		return -1, nil
	}

	return 0, nil
}
