package enum

type CityAdminRole int

const (
	Owner CityAdminRole = iota
	Admin
	Moderator
)

var citiesAdminsRoles = []string{
	"owner",
	"admin",
	"moderator",
}

func (r CityAdminRole) String() string {
	if int(r) < len(citiesAdminsRoles) {
		return citiesAdminsRoles[r]
	}
	return ""
}

func CheckRole(role string) bool {
	for _, name := range citiesAdminsRoles {
		if name == role {
			return true
		}
	}
	return false
}

func GetAllCitiesAdminsRoles() []string {
	return citiesAdminsRoles
}
