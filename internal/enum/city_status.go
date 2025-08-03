package enum

type CityStatus int

const (
	CityNotSupport CityStatus = iota
	CitySupport
)

var citiesStatus = []string{
	"not_support",
	"support",
}

func (s CityStatus) String() string {
	if int(s) < len(citiesStatus) {
		return citiesStatus[s]
	}
	return ""
}

func CheckCityStatus(s string) bool {
	for _, name := range citiesStatus {
		if name == s {
			return true
		}
	}
	return false
}

func GetAllCitiesStatuses() []string {
	return citiesStatus
}
