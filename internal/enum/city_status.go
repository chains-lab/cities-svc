package enum

type CityStatus string

const (
	CityStatusSupported   CityStatus = "supported"
	CityStatusSuspended   CityStatus = "suspended"
	CityStatusUnsupported CityStatus = "unsupported"
)

var cityStatuses = []CityStatus{
	CityStatusSupported,
	CityStatusSuspended,
	CityStatusUnsupported,
}

func ParseCityStatus(status string) (CityStatus, bool) {
	for _, s := range cityStatuses {
		if s == CityStatus(status) {
			return s, true
		}
	}

	return "", false
}

func GetAllCitiesStatuses() []CityStatus {
	return cityStatuses
}
