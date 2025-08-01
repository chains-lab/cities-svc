package enum

type CityAdminRole int

const (
	Owner CityAdminRole = iota
	Admin
	Moderator
)

var roleNames = [...]string{
	"owner",
	"admin",
	"moderator",
}

func (r CityAdminRole) String() string {
	if int(r) < len(roleNames) {
		return roleNames[r]
	}
	return "unknown"
}

func Compare(r1, r2 CityAdminRole) int {
	return int(r1) - int(r2)
}
